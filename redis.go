package faststorage

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func (rd *RedisDB) Put(ctx context.Context, asset RedisAsset) (interface{}, error) {
	c := rd.GetConn()
	defer c.Close()

	return c.Do("HSET", redis.Args{asset.GetKey()}.AddFlat(asset)...)
}

func (rd *RedisDB) Get(ctx context.Context, asset RedisAsset, args ...interface{}) (error){
	c := rd.GetConn()
	defer c.Close()

	args = append([]interface{}{asset.GetKey()}, args)

	val, err := c.Do("HGET", asset.GetKey())
	fmt.Printf("Val: %+v\n", val)
	value, err := redis.Values(val, err)

	err = redis.ScanStruct(value, asset)
	if err != nil {
		return err
	}
	return nil
}
