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

	/*
	DSAsset is the interface that a struct must implement to use the Datastore operations.
	GetDSKind should return the Datastore Kind for the object to be save/get from.
	GetDSNameKey and GetDSIDKey: should return a identifier to be used in the Datastore Key along with a bool flag indicating which one (name key or id key) should be used.
	- If both true, Id key will be used, if both false will use Datasore Incomplete Key.
	GetDSNamespace should return the Datastore Namespace to be used.
	 */
	DSAsset interface {
		GetDSKind() string
		GetDSNameKey() (string, bool)
		GetDSIDKey() (int64, bool)
		GetDSNamespace() string
	}

	/*
	RedisAsset is the interface that a struct must implement to use the Redis operations.
	GetRedisKey should return the Key for the Redis Hash
	GetRedisField should return the Field for the Redis Hash
	GetTTL should return the duration of ttl for a key. Min value: 1 second

	The Redis operations rely on encoding/json Marshal and Unmarshal methods of the RedisAsset implementer struct. The bytes slice on these methods are use to load and save the structs in Redis
	 */
	RedisAsset interface {
		GetRedisKey() interface{}
		GetRedisField() string
		GetTTL() time.Duration
	}

	/*
	Asset is a interface combining both Redis and Datastore operations. A struct must implement this interface to take advantage of the combined usage of both Redis and Datastore.
	 */
	Asset interface {
		DSAsset
		RedisAsset
	}

	/*
	Config struct to use with the NewWithCfg wrapper.
 	*/
	Cfg struct {
		RedisAddrs      string
		RedisPassword   string
		Db              string
		MaxIdle         int
		MaxActive       int
		IdleTimeout     time.Duration
		MaxConnLifetime time.Duration
		Wait            bool
		TestOnBorrow    func(c redis.Conn, t time.Time) (error)
		Client          *datastore.Client
	}

	AssetMissError struct {
		DSmissErr DSmissErr
		CachemissErr CachemissErr
	}
)

func NewWithCfg(cfg Cfg) (*DataLayer, error) {
	dial := func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", cfg.RedisAddrs)

		if err != nil {
			return nil, err
		}
		if _, err := c.Do("AUTH", cfg.RedisPassword); cfg.RedisPassword != "" && err != nil {
			return nil, err
		}
		if _, err := c.Do("SELECT", cfg.Db); cfg.Db != "" && err != nil {
			return nil, err
		}
		return c, nil
	}

	dl := &DataLayer{
		Redis: RedisDB{
			redisPool: &redis.Pool{
				Dial:            dial,
				MaxActive:       cfg.MaxActive,
				MaxConnLifetime: cfg.MaxConnLifetime,
				IdleTimeout:     cfg.IdleTimeout,
				MaxIdle:         cfg.MaxIdle,
				Wait:            cfg.Wait,
				TestOnBorrow:    cfg.TestOnBorrow,
			}},
		DS: DatastoreDB{
			client: cfg.Client,
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
	redisReply, err := dl.Redis.Put(ctx, asset)
	if err != nil {
		return dsKey, nil, err
	}
	return dsKey, redisReply, nil
}

/*
Lookup in Redis db for the asset Id, in case of cahce miss, try to get it from Datastore.
If asset is present in Datastore it will put it back in Redis and return the asset found in Datastore
 */
func (dl *DataLayer) Get(ctx context.Context, asset Asset) error {
	err := dl.Redis.Get(ctx, asset, []interface{}{})
	if err != nil {
		if  cacheMiss, ok := err.(CachemissErr) ; ok{

			err = dl.DS.Get(ctx, asset)

			if dsMiss, ok := err.(DSmissErr) ; err != nil && ok {
				return AssetMissError{CachemissErr:cacheMiss, DSmissErr:dsMiss}
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

func (amr AssetMissError) Error() string {
	return fmt.Sprintf("asset not found: %s - %s", amr.CachemissErr.Error(), amr.DSmissErr.Error())
}

/*
Asset impl for testing
*/
type TestAsset struct {
	Id   string
	Val1 string
	Val2 []string
	Val3 int
	Val4 *Sub
}

type Sub struct {
	SubVal string
}

func (ta TestAsset) GetDSKind() string {
	return "test-kind"
}

func (ta TestAsset) GetDSNameKey() (string, bool) {
	return ta.Id, true
}

func (ta TestAsset) GetDSIDKey() (int64, bool) {
	return 0, false
}

func (ta TestAsset) GetDSNamespace() string {
	return "test-namespace"
}

func (ta TestAsset) GetRedisKey() interface{} {
	return ta.Id
}

func (ta TestAsset) GetRedisField() string {
	return "test2"
}

func (ta TestAsset) GetTTL() time.Duration {
	return 1 * time.Hour
}
