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

	mockedConnFn, cmd := getMockedConn(&TestAsset{Id: "asset:1", Val1: "asset1", Val2:[]string{"a", "b"}}, "HMGET", "asset:1", "TestAsset.Val1", "TestAsset.Val2")
	dl.Redis.GetConn = mockedConnFn

	ctx := context.Background()
	asset := TestAsset{Id: "asset:1"}
	err := dl.Redis.Get(ctx, &asset)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		t.Fail()
	}

	ccc := mockedConnFn().(*redigomock.Conn)
	fmt.Printf("stats %d\n", ccc.Stats(cmd))

	fmt.Printf("test reply: %+v\n", asset)

	if asset.Val1 != "asset1" {
		fmt.Printf("expect 'Val1: asset1' got %+v\n", asset)
		t.Fail()
	}
}

/**
Mocked Conn apply the expected reply to the redis conn and returns it to be used by the client
 */
func getMockedConn(expectResp RedisAsset, redisCmd string, args ...interface{}) (func() (redis.Conn), *redigomock.Cmd){
	mockedConn := redigomock.NewConn()

	fmt.Printf("cmd args: rediscmd %s, args: %s, expected: %+v\n", redisCmd, args, redis.Args{}.AddFlat(expectResp))
	cmd := mockedConn.Command(redisCmd, args).Expect(redis.Args{}.AddFlat(expectResp))

	getConnFn := func() (redis.Conn) {
		return mockedConn
	}
	return getConnFn, cmd
}
