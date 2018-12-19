package faststorage

import (
	"context"
	"encoding/json"
	"fmt"
)

func (rd *RedisDB) Put(ctx context.Context, asset RedisAsset) (interface{}, error) {
	payload, err := json.Marshal(asset)
	if err != nil {
		return nil, err
	}
	if asset.GetField() == "" {
		return nil, fmt.Errorf("cannot use empty string as redis field - asset key: %s", asset.GetKey())
	}
	c := rd.GetConn()
	defer c.Close()
	return c.Do("HSET", asset.GetKey(), asset.GetField(), payload)
}

func (rd *RedisDB) Get(ctx context.Context, asset RedisAsset) (error){
	c := rd.GetConn()

	var reply interface{}
	var err error
	if asset.GetField() != "" {
		reply, err = c.Do("HGET", asset.GetKey(), asset.GetField())
	} else {
		reply, err = c.Do("HGETALL", asset.GetKey())
	}
	c.Close()
	if err != nil {
		return err
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
