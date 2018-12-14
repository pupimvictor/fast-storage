package faststorage

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"testing"
	"time"
)

func TestGet(t *testing.T){
	t.Skip()
	dl, err := newTestDB()
	if err != nil {
		fmt.Printf("err newTestDB: %v\n", err)
		t.Fail()
	}

	mockedConnFn, _ := getMockedConn(TestAsset{Id: "asset:1", Val: "asset1"}, "HGET", "asset:1")
	dl.Redis.GetConn = mockedConnFn

	ctx := context.Background()
	asset := TestAsset{Id: "asset:1"}

	err = dl.Get(ctx, &asset)

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

type TestAsset struct{
	Id  string
	Val string
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




