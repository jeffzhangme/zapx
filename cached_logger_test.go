package zapx_test

import (
	"net/url"
	"testing"

	"github.com/jeffzhangme/zapx"

	"go.uber.org/zap"
)

func TestZapLogger(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("msg string", zap.String("k", "v"))
	sugar := logger.Sugar()
	defer sugar.Sync()
	sugar.Infof("key %d", 1)
}

func TestCachedLogger(t *testing.T) {
	logger, _ := zapx.NewCachedLoggerConfig().Build()
	defer logger.Flush()
	logger.Info("key", zap.Int("key string", 1))
}

func TestCachedSugarLogger(t *testing.T) {
	logger, _ := zapx.NewCachedLoggerConfig().Build()
	sugar := logger.CachedSugar()
	sugar.Infow("key.subkey", "k", "v")
	sugar.Flush()
}

func BenchmarkZapLogger(b *testing.B) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		logger.Info("key", zap.String("k", "v"))
	}
}

func BenchmarkCachedLoggerStd(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger, _ := zapx.NewCachedLoggerConfig().Build()
			logger.Info("key", zap.String("k", "v"))
			logger.Flush()
		}
	})
}

func BenchmarkCachedLoggerSink(b *testing.B) {
	// sinkUrl := SinkURL{url.URL{Scheme: "redis", Host: "127.0.0.1:6379", RawQuery: "db=0&type=list&key=log:for:test"}}
	sinkUrl := zapx.SinkURL{url.URL{Scheme: "kafka", Host: "127.0.0.1:9092", RawQuery: "topic=test_log_topic"}}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log, _ := zapx.NewCachedLoggerConfig().AddSinks(sinkUrl).Build()
			log.Info("key", zap.String("k", "v"))
			log.Flush()
		}
	})
}
