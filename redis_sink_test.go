package zapx_test

import (
	"net/url"
	"testing"

	"github.com/jeffzhangme/zapx"
	"go.uber.org/zap"
)

func TestRedisSink(t *testing.T) {
	stderr := zapx.SinkURL{url.URL{Opaque: "stderr"}}
	sinkUrl := zapx.SinkURL{url.URL{Scheme: zapx.SchemeRedis, Host: "127.0.0.1:6379", RawQuery: "db=0&type=list&key=log:for:test"}}
	logger, _ := zapx.NewCachedLoggerConfig().AddSinks(stderr, sinkUrl).Build()
	defer logger.Flush()
	logger.Info("key", zap.String("k", "v"))
}
