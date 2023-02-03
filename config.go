package zapx

import (
	"net/url"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	sink    zapcore.WriteSyncer
	encoder zapcore.Encoder
	level   zapcore.Level
)

func init() {
	err := zap.RegisterSink(SchemeRedis, InitRedisSink)
	if err != nil {
		panic(err)
	}
	err = zap.RegisterSink(SchemeKafka, InitKafkaSink)
	if err != nil {
		panic(err)
	}
	eConfig := zap.NewProductionEncoderConfig()
	eConfig.LineEnding = cachedLogLineEnding
	eConfig.MessageKey = cachedLogMessageKey
	sink, _, _ = zap.Open(cachedLogSinkURL)
	encoder = zapcore.NewJSONEncoder(eConfig)
}

// ENV env
type ENV int

const (
	// DEV dev
	DEV ENV = iota
	// TEST test
	TEST
	// PROD prod
	PROD
)

const (
	SchemeRedis = "redis"
	SchemeKafka = "kafka"
)

const (
	redisPubSubType      = "channel"
	redisDefaultType     = "list"
	redisClusterNodesKey = "nodes"
	redisDefaultPwd      = ""
	redisDefaultKey      = "just_a_test_key"
	kafkaDefaultTopic    = "just_a_test_topic"
	kafkaAsyncKey        = "isAsync"

	cachedLogLineEnding = ","
	cachedLogMessageKey = ""
	cachedLogSinkURL    = "stderr"
	globalKeyPrefix     = "__"
)

var stdErrSink, _, _ = zap.Open("stderr")

// SinkURL sink url
type SinkURL struct {
	url.URL
}

// CachedLogConfig cached log config
type CachedLogConfig struct {
	sinkURLs         []SinkURL
	level            zapcore.Level
	env              ENV
	econfig          zapcore.EncoderConfig
	isEConfigChanged bool
}

// Level set level
func (cl CachedLogConfig) Level(level zapcore.Level) CachedLogConfig {
	cl.level = level
	return cl
}

// AddSinks add sinks
func (cl CachedLogConfig) AddSinks(url ...SinkURL) CachedLogConfig {
	cl.sinkURLs = append(cl.sinkURLs, url...)
	return cl
}

// Env set env
func (cl CachedLogConfig) Env(env ENV) CachedLogConfig {
	cl.env = env
	return cl
}

// EncoderConfig set encoder config
func (cl CachedLogConfig) EncoderConfig(econfig zapcore.EncoderConfig) CachedLogConfig {
	cl.econfig = econfig
	cl.isEConfigChanged = true
	return cl
}

// Build create logger
func (cl CachedLogConfig) Build() (*CachedLogger, error) {
	ws := stdErrSink
	core := getCachedCore()
	core.LevelEnabler = level
	urls := make([]string, 0, len(cl.sinkURLs)+1)
	urls = append(urls, cachedLogSinkURL)
	if len(cl.sinkURLs) > 0 {
		urls = urls[:0]
		for _, url := range cl.sinkURLs {
			urls = append(urls, url.String())
		}
		var err error
		ws, _, err = zap.Open(urls...)
		if err != nil {
			return nil, err
		}
	}
	core.out = ws
	if cl.isEConfigChanged {
		core.enc = zapcore.NewJSONEncoder(cl.econfig)
	}
	switch cl.env {
	case DEV:
		logger := zap.New(core, zap.Development(), zap.AddCaller())
		return &CachedLogger{Logger: logger}, nil
	case PROD:
		logger := zap.New(core)
		return &CachedLogger{Logger: logger}, nil
	default:
		logger := zap.New(core)
		return &CachedLogger{Logger: logger}, nil
	}
}

// NewCachedLoggerConfig create logger config
func NewCachedLoggerConfig() CachedLogConfig {
	return CachedLogConfig{env: PROD}
}
