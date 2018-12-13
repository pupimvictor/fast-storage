package faststorage

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"testing"
)

func TestRedisDB_Get(t *testing.T){
	dl := &DataLayer{
		Redis: RedisDB{
			redisPool: nil,
		},
	}

	mockedConnFn := getMockedConn("asset1", "HGET", "asset:1", "a")
	dl.Redis.GetConn = mockedConnFn

	ctx := context.Background()
	r, err := dl.Redis.Get(ctx, "asset:1", "a")

	if err != nil {
		fmt.Printf("err: %v\n", err)
		t.Fail()
	}

	fmt.Printf("reply: %+v\n", r)

	if r != "asset1" {
		fmt.Printf("expect 'asset1' got %v\n", r)
		t.Fail()
	}
}


func getMockedConn(expectResp interface{}, redisCmd string, args ...interface{}) (func() (redis.Conn)){

	mockedConn := redigomock.NewConn()
	mockedConn.Command(redisCmd, args...).Expect(expectResp)

	var connInterface interface{}
	connInterface = mockedConn
	conn := connInterface.(redis.Conn)

	getConnFn := func() (redis.Conn) {
		return conn
	}
	return getConnFn
}
