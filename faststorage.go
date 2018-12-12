package faststorage

import (
	"cloud.google.com/go/datastore"
	"github.com/gomodule/redigo/redis"
	"time"
)

type (

	RedisDB struct{
		redisPool *redis.Pool
	}

	DatastoreDB struct{
		client datastore.Client
	}

	DataLayer struct {
		Redis RedisDB
		DS    DatastoreDB
	}

	DSAsset interface{
		GetDSKind() string
		GetDSStrKey() string
		GetDSIntKey() int64
		GetDSNamespace() string
	}

	RedisAsset interface{
		GetKey() interface{}
		GetTTL() time.Duration
		GetStructType() string
	}

	Asset interface{
		DSAsset
		RedisAsset
	}
)

func New(redisAddrs, redisPassword string, maxIdle, maxActive int, idleTimeout, maxConnLifetime time.Duration, wait bool) (DataLayer, error) {
	


}