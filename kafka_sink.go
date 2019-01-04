package zapx

import (
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

var (
	kafkaSinkInsts = map[string]kafkaSink{}
)

type kafkaSink struct {
	kafkaProducer sarama.SyncProducer
	isAsync       bool
	topic         string
}

func getKafkaSink(brokers []string, topic string, config *sarama.Config) kafkaSink {
	producerInst, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		panic(err)
	}
	kafkaSinkInst := kafkaSink{
		kafkaProducer: producerInst,
		topic:         topic,
	}
	return kafkaSinkInst
}

// InitKafkaSink  create kafka sink instance
func InitKafkaSink(u *url.URL) (zap.Sink, error) {
	topic := kafkaDefaultTopic
	if t := u.Query().Get("topic"); len(t) > 0 {
		topic = t
	}
	brokers := []string{u.Host}
	instKey := strings.Join(brokers, ",")
	if v, ok := kafkaSinkInsts[instKey]; ok {
		return v, nil
	}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	if ack := u.Query().Get("acks"); len(ack) > 0 {
		if iack, err := strconv.Atoi(ack); err == nil {
			config.Producer.RequiredAcks = sarama.RequiredAcks(iack)
		} else {
			log.Printf("kafka producer acks value '%s' invalid  use default value %d\n", ack, config.Producer.RequiredAcks)
		}
	}
	if retries := u.Query().Get("retries"); len(retries) > 0 {
		if iretries, err := strconv.Atoi(retries); err == nil {
			config.Producer.Retry.Max = iretries
		} else {
			log.Printf("kafka producer retries value '%s' invalid  use default value %d\n", retries, config.Producer.Retry.Max)
		}
	}
	kafkaSinkInsts[instKey] = getKafkaSink(brokers, topic, config)
	return kafkaSinkInsts[instKey], nil
}

// Close implement zap.Sink func Close
func (p kafkaSink) Close() error {
	return nil
}

// Write implement zap.Sink func Write
func (p kafkaSink) Write(b []byte) (n int, err error) {
	var multiErr MultiError
	for _, topic := range strings.Split(p.topic, ",") {
		_, _, err = p.kafkaProducer.SendMessage(&sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(time.Now().String()),
			Value: sarama.ByteEncoder(b),
		})
		if err != nil {
			multiErr = append(multiErr, err)
		}
	}
	return len(b), multiErr
}

// Sync implement zap.Sink func Sync
func (p kafkaSink) Sync() error {
	return nil
}
