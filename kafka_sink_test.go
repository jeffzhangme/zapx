package zapx_test

import (
	"net/url"
	"sync"
	"testing"
	. "zapx"

	"go.uber.org/zap"
)

func TestKafkaSink(t *testing.T) {
	var wg sync.WaitGroup
	stderr := SinkURL{url.URL{Opaque: "stderr"}}
	sinkUrl := SinkURL{url.URL{Scheme: "kafka", Host: "127.0.0.1:9092", RawQuery: "topic=test_log_topic"}}
	logger, _ := NewCachedLoggerConfig().AddSinks(stderr, sinkUrl).Build()
	defer logger.Flush(&wg)
	logger.Info("key", zap.String("k", "v"))
}
