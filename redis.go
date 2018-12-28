package faststorage

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
)

func (rd *RedisDB) Put(ctx context.Context, asset RedisAsset) (interface{}, error) {
	payload, err := json.Marshal(asset)
	if err != nil {
		return nil, err
	}
	if asset.GetRedisField() == "" {
		return nil, fmt.Errorf("cannot use empty string as redis field - asset key: %s", asset.GetRedisKey())
	}
	c := rd.GetConn()
	defer c.Close()

	reply, err := c.Do("HSET", asset.GetRedisKey(), asset.GetRedisField(), payload)
	if err != nil {
		return reply, err
	}

	if ttl := math.Round(asset.GetTTL().Seconds()) ; ttl > 0 {
		return c.Do("EXPIRE", asset.GetRedisKey(), ttl)
	}

	return reply, err
}

func (rd *RedisDB) Get(ctx context.Context, asset RedisAsset, args ...interface{}) error {
	c := rd.GetConn()

	var reply interface{}
	var err error
	if asset.GetRedisField() != "" {
		reply, err = c.Do("HGET", asset.GetRedisKey(), asset.GetRedisField())
	} else {
		reply, err = c.Do("HGETALL", asset.GetRedisKey())
	}
	c.Close()
	if err != nil {
		return err
	}
	if reply == nil {
		return CachemissErr{AssetKey: fmt.Sprintf("%s", asset.GetRedisKey()), AssetField: asset.GetRedisField()}
	}

	var b []byte
	if v, ok := reply.([]interface{}); ok {
		for _, val := range v {
			r := val.([]byte)
			b = append(b, r...)
		}
	} else {
		b = reply.([]byte)
	}
	json.Unmarshal(b, asset)

	return nil
}

type CachemissErr struct {
	AssetKey string
	AssetField string
}

func (e CachemissErr) Error() string {
	return fmt.Sprintf("asset not found in Redis - Key: %s - Field: %s", e.AssetKey, e.AssetField)
}
