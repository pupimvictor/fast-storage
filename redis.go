package faststorage

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"reflect"
	"time"
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

	buildSetArgs(asset)
	val, err := c.Do("HGET", asset.GetKey())
	fmt.Printf("Val: %+v\n", val)

	replyAsset := &TestAsset{}
	buildAsset(val, replyAsset)
	fmt.Printf("replyasset: %+v\n", replyAsset)

	value, err := redis.Values(val, err)

	err = redis.ScanStruct(value, asset)
	if err != nil {
		return err
	}
	return nil
}

//todo check for value != nil
func buildSetArgs(asset RedisAsset) ([]string, error) {
	fmt.Printf("asset::: %+v\n", asset)
	var args []string
	t := reflect.TypeOf(asset).Elem()
	v := reflect.ValueOf(asset).Elem()
	for i:= 0 ; i< t.NumField(); i++{
		n := t.Field(i)
		v := v.Field(i)
		fieldName := n.Tag.Get("redis")
		fieldVal := v.Interface()
		args = append(args, fieldName, fmt.Sprint(fieldVal))
	}
	fmt.Println(args)
	return args, nil
}

func buildGetArgs(asset RedisAsset) ([]string, error) {
	t := reflect.TypeOf(asset).Elem()
	fields := make([]string, t.NumField())
	for i:= 0 ; i< t.NumField(); i++ {
		n := t.Field(i)
		fieldName := n.Tag.Get("redis")
		fields[i] = fieldName
	}
	return fields, nil
}

func buildAsset(reply interface{}, asset RedisAsset) (error) {
	//need to check how the redigo reply looks like
	//a := reply.([]string)

	return nil
}












//






//
type TestAsset struct{
	Id  string `redis:"TestAsset.Id"`
	Val string `redis:"TestAsset.Val"`
	A   []string `redis:"TestAsset.A"`
	S   Sub `redis:"TestAsset.S"`
}

type Sub struct{
	B []string `redis:Sub.B`
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




