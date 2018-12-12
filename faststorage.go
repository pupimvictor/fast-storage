package faststorage

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/gomodule/redigo/redis"
	"time"
)

type (
	RedisDB struct {
		redisPool *redis.Pool
	}

	DatastoreDB struct {
		client *datastore.Client
	}

	DataLayer struct {
		Redis RedisDB
		DS    DatastoreDB
	}

	DSAsset interface {
		GetDSKind() string
		GetNameKey() (string, bool)
		GetIDKey() (int64, bool)
		GetDSNamespace() string
	}

	RedisAsset interface {
		GetKey() interface{}
		GetTTL() time.Duration
		GetStructType() string
	}

	Asset interface {
		DSAsset
		RedisAsset
	}
)

func New(redisAddrs, redisPassword, db string, maxIdle, maxActive int, idleTimeout, maxConnLifetime time.Duration, wait bool, client *datastore.Client) (*DataLayer, error) {
	dial := func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", redisAddrs)
		if err != nil {
			return nil, err
		}
		if redisPassword != "" {
			if _, err := c.Do("AUTH", redisPassword); err != nil {
				c.Close()
				return nil, err
			}
		}
		if _, err := c.Do("SELECT", db); err != nil {
			c.Close()
			return nil, err
		}
		return c, nil
	}

	return &DataLayer{
		Redis: RedisDB{
			redisPool: &redis.Pool{
				Dial:            dial,
				MaxActive:       maxActive,
				MaxConnLifetime: maxConnLifetime,
				IdleTimeout:     idleTimeout,
				MaxIdle:         maxIdle,
				Wait:            wait,
			}},
		DS: DatastoreDB{
			client: client,
		},
	}, nil
}

func (dl *DataLayer) Put(ctx context.Context, asset Asset, dsParent *datastore.Key) (*datastore.Key, *string, error) {
	dsKey, err := dl.DS.Put(ctx, asset, dsParent)
	if err != nil {
		return nil, nil, err
	}
	redisKey, err := dl.Redis.Put(ctx, asset)
	if err != nil {
		return dsKey, nil, err
	}
	return dsKey, redisKey, nil
}

