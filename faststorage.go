package faststorage

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

type (
	RedisDB struct {
		redisPool *redis.Pool
		GetConn func() (redis.Conn)
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
		defer c.Close()

		fmt.Println("new conn!")
		if err != nil {
			return nil, err
		}
		if redisPassword != "" {
			if _, err := c.Do("AUTH", redisPassword); err != nil {
				return nil, err
			}
		}
		if _, err := c.Do("SELECT", db); db != "" && err != nil {
			return nil, err
		}
		return c, nil
	}

	dl := &DataLayer{
		Redis: RedisDB{
			redisPool: &redis.Pool{
				Dial:            dial,
				MaxActive:       maxActive,
				MaxConnLifetime: maxConnLifetime,
				IdleTimeout:     idleTimeout,
				MaxIdle:         maxIdle,
				Wait:            wait,
				TestOnBorrow: func(c redis.Conn, t time.Time) error {
					_, err := c.Do("PING")
					return err
				},
			}},
		DS: DatastoreDB{
			client: client,
		},
	}
	dl.Redis.GetConn = dl.Redis.redisPool.Get
	return dl, nil
}

/*
Put asset in Datastore and Redis. Uses the Asset interface to get info about asset (namespace, kind, TTL, etc)
 */
func (dl *DataLayer) Put(ctx context.Context, asset Asset, dsParent *datastore.Key) (*datastore.Key, interface{}, error) {
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

/*
Lookup in Redis db for the asset Id, in case of cahce miss, try to get it from Datastore.
If asset is present in Datastore it will put it back in Redis and return the asset found in Datastore
 */
func (dl *DataLayer) Get(ctx context.Context, asset Asset) (error) {
	err := dl.Redis.Get(ctx, asset, []interface{}{})
	if err != nil {
		if err.Error() == "not found in cache" /*todo: create error struct for this*/ {
			err = dl.DS.Get(ctx, asset)
			if err != nil {
				return err
			}
			_, err = dl.Redis.Put(ctx, asset)
			if err != nil {
				return err
			}
		}
		return err
	}
	return nil
}

/*
Asset impl for testing
 */
type TestAsset struct{
	Id   string   `redis:"-"`
	Val1 string   `redis:"TestAsset.Val1"`
	Val2 []string `redis:"TestAsset.Val2"`
}

func (ta TestAsset) GetDSKind() string {
	return "test-kind"
}

func (ta TestAsset) GetNameKey() (string, bool) {
	return ta.Id, true
}

func (ta TestAsset) GetIDKey() (int64, bool) {
	return 0, false
}

func (ta TestAsset) GetDSNamespace() string {
	return "test-namespace"
}

func (ta TestAsset) GetKey() interface{} {
	return ta.Id
}

func (ta TestAsset) GetTTL() time.Duration {
	return 1 * time.Hour
}

func (ta TestAsset) GetStructType() string {
	return "HASH"
}


