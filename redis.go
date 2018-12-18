package faststorage

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"reflect"
)

func (rd *RedisDB) Put(ctx context.Context, asset RedisAsset) (interface{}, error) {
	c := rd.GetConn()
	defer c.Close()

	return c.Do("HMSET", redis.Args{asset.GetKey()}.AddFlat(asset)...)}

func (rd *RedisDB) Get(ctx context.Context, asset RedisAsset, args ...interface{}) (error){
	fields, err := scanAssetSctruct(asset)
	if err != nil {
		return err
	}
	args = append([]interface{}{asset.GetKey()}, fields, args)

	fmt.Printf("get args %+v\n", args)

	c := rd.GetConn()
	defer c.Close()

	val, err := c.Do("HMGET", asset.GetKey(), args)

	value, err := redis.Values(val, err)
	err = redis.ScanStruct(value, asset)
	if err != nil {
		return err
	}
	return nil
}

func scanAssetSctruct(asset RedisAsset) ([]string, error){
	var fields []string
	t := reflect.TypeOf(asset)
	for i := 0; i < t.NumField() ; i++ {
		tag := t.FieldByIndex([]int{i}).Tag.Get("redis")
		if tag != "-" {
			fields = append(fields, tag)
		}
	}
	return fields, nil
}