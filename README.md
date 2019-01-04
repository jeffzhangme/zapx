# zapx

Based on [uber-go/zap](https://github.com/uber-go/zap), log caching and redis/kafka sink implemented in Go

## Doc

[zapx](https://godoc.org/github.com/jeffzhangme/zapx) ( Please read the [uber-go/zap](https://github.com/uber-go/zap/blob/master/README.md) documentation first. )

## Quick Start

Using redis sink

```go
stderr := zapx.SinkURL{url.URL{Opaque: "stderr"}}
sinkUrl := zapx.SinkURL{url.URL{Scheme: "redis", Host: "127.0.0.1:6379", RawQuery: "db=0&type=list&key=log:for:test"}}
logger, _ := zapx.NewCachedLoggerConfig().AddSinks(stderr, sinkUrl).Build()
defer logger.Flush(nil)
logger.Info("key", zap.String("k", "v"))
```
Log example

```json
{
    "key": {
        "level": "info",
        "ts": 1546572659.1465247,
        "k": "v"
    },
    "level": "info",
    "ts": 1546572659.1465852
}
```