package zapx_test

import (
	"net/url"
	"sync"
	"testing"
	. "zapx"

	"go.uber.org/zap"
)

func TestRedisSink(t *testing.T) {
	var wg sync.WaitGroup
	stderr := SinkURL{url.URL{Opaque: "stderr"}}
	sinkUrl := SinkURL{url.URL{Scheme: "redis", Host: "127.0.0.1:6379", RawQuery: "db=0&type=list&key=log:for:test"}}
	logger, _ := NewCachedLoggerConfig().AddSinks(stderr, sinkUrl).Build()
	defer logger.Flush(&wg)
	logger.Info("key", zap.String("k", "v"))
}
