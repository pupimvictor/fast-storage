package faststorage

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func (rd *RedisDB) Put(ctx context.Context, asset Asset) (interface{}, error) {
	c := rd.redisPool.Get()
	defer c.Close()

	structType := asset.GetStructType()
	if structType == "HASH" {
		reply, err := c.Do("HSET", redis.Args{asset.GetKey()}.AddFlat(asset))
		if err != nil {
			return nil, err
		}
		return reply, nil
	}

	return nil, fmt.Errorf("not supported structType: %s", structType)
}

