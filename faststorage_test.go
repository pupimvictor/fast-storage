package faststorage

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"testing"
)

func TestGet(t *testing.T){
	dl, err := newTestDB()
	if err != nil {
		fmt.Printf("err newTestDB: %v\n", err)
		t.Fail()
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

func newTestDB() (*DataLayer, error) {
	ctx := context.Background()
	cli, err := datastore.NewClient(ctx, "test")
	if err != nil {
		return nil, err
	}

	return &DataLayer{
		Redis: RedisDB{
			redisPool: nil,
		},
		DS: DatastoreDB{
			client: cli,
		},
	}, nil

}
