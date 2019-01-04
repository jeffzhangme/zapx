package zapx

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

var (
	redisSinkInsts = map[string]redisSink{}
)

func getRedisSink(host, pwd string, db int, typee, key string) redisSink {
	redisSinkInst := redisSink{
		redisClient: *redis.NewClient(&redis.Options{
			Addr:     host,
			Password: pwd,
			DB:       db,
		}),
		key:   key,
		typee: typee,
	}
	if err := redisSinkInst.redisClient.Ping().Err(); err != nil {
		panic(err)
	}
	return redisSinkInst
}

// InitRedisSink init redis sink
func InitRedisSink(u *url.URL) (zap.Sink, error) {
	var pwd, key, typee = redisDefaultPwd, redisDefaultKey, redisDefaultType
	var db int
	db, _ = strconv.Atoi(u.Query().Get("db"))
	if k := u.Query().Get("key"); len(k) > 0 {
		key = k
	}
	if t := u.Query().Get("type"); len(t) > 0 {
		typee = t
	}
	instKey := u.Host + typee + strconv.Itoa(db)
	if v, ok := redisSinkInsts[instKey]; ok {
		return v, nil
	}
	if u.User != nil {
		pwd, _ = u.User.Password()
	}
	redisSinkInsts[instKey] = getRedisSink(u.Host, pwd, db, typee, key)
	return redisSinkInsts[instKey], nil
}

type redisSink struct {
	redisClient redis.Client
	key         string
	typee       string
	isCluster   bool
}

// Close implement zap.Sink func Close
func (p redisSink) Close() error {
	return nil
}

// Write implement zap.Sink func Write
func (p redisSink) Write(b []byte) (n int, err error) {
	var multiErr MultiError
	switch p.typee {
	case redisPubSubType:
		for _, key := range strings.Split(p.key, ",") {
			if err := p.redisClient.Publish(key, string(b)).Err(); err != nil {
				multiErr = append(multiErr, err)
			}
		}
	case redisDefaultType:
		for _, key := range strings.Split(p.key, ",") {
			if err := p.redisClient.RPush(key, string(b)).Err(); err != nil {
				multiErr = append(multiErr, err)
			}
		}
	default:
		for _, key := range strings.Split(p.key, ",") {
			if err := p.redisClient.RPush(key, string(b)).Err(); err != nil {
				multiErr = append(multiErr, err)
			}
		}
	}
	return len(b), multiErr
}

// Sync implement zap.Sink func Sync
func (p redisSink) Sync() error {
	return nil
}
