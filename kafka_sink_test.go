package zapx_test

import (
	"net/url"
	"testing"

	"github.com/jeffzhangme/zapx"
	"go.uber.org/zap"
)

func TestKafkaSink(t *testing.T) {
	stderr := zapx.SinkURL{url.URL{Opaque: "stderr"}}
	sinkUrl := zapx.SinkURL{url.URL{Scheme: zapx.SchemeKafka, Host: "127.0.0.1:9092", RawQuery: "topic=test_log_topic"}}
	logger, _ := zapx.NewCachedLoggerConfig().AddSinks(stderr, sinkUrl).Build()
	defer logger.Flush()
	logger.Info("key", zap.String("k", "v"))
}
