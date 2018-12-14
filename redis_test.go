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

	mockedConnFn, _ := getMockedConn(&TestAsset{Id: "asset:1", Val: "asset1"}, "HGET", "asset:1")
	dl.Redis.GetConn = mockedConnFn

	ctx := context.Background()
	asset := TestAsset{Id: "asset:1"}
	err := dl.Redis.Get(ctx, &asset, []interface{}{})

	if err != nil {
		fmt.Printf("err: %v\n", err)
		t.Fail()
	}

	fmt.Printf("reply: %+v\n", asset)

	if asset.Val != "asset1" {
		fmt.Printf("expect 'Val: asset1' got %+v\n", asset)
		t.Fail()
	}
}


func getMockedConn(expectResp RedisAsset, redisCmd string, args ...interface{}) (func() (redis.Conn), *redigomock.Cmd){
	mockedConn := redigomock.NewConn()
	cmd := mockedConn.Command(redisCmd, args...).Expect(redis.Args{}.AddFlat(expectResp))

	var connInterface interface{}
	connInterface = mockedConn
	conn := connInterface.(redis.Conn)

	getConnFn := func() (redis.Conn) {
		return conn
	}
	return getConnFn, cmd
}
