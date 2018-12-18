package faststorage

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"reflect"
	"strings"
)

func (rd *RedisDB) Put(ctx context.Context, asset RedisAsset) (interface{}, error) {
	c := rd.GetConn()
	defer c.Close()

	return c.Do("HMSET", redis.Args{asset.GetKey()}.AddFlat(asset)...)}

func (rd *RedisDB) Get(ctx context.Context, asset RedisAsset) (error){
	fields, err := scanAssetSctruct(asset)
	if err != nil {
		return err
	}
	args := append([]interface{}{asset.GetKey()}, strings.Join(fields, " "))

	fmt.Printf("get args %+v\n", args)

	c := rd.GetConn()
	defer c.Close()

	ci := c.(*redigomock.Conn)

	val, err := ci.Do("HMGET", args)

	fmt.Printf("reply: %+v\n", val)

	value, err := redis.Values(val, err)
	err = redis.ScanStruct(value, asset)
	if err != nil {
		return err
	}
	return nil
}

func scanAssetSctruct(asset RedisAsset) ([]string, error){
	var fields []string
	t := reflect.TypeOf(asset).Elem()
	for i := 0; i < t.NumField() ; i++ {
		tag := t.FieldByIndex([]int{i}).Tag.Get("redis")
		if tag != "-" {
			fields = append(fields, tag)
		}
	}
	return fields, nil
}