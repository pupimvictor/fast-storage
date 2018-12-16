package main

import (
	"context"
	"fmt"
	"github.com/pupimvictor/fast-storage"
	"time"
)

func main() {
	fmt.Println("hello fast storage")

	ts := faststorage.TestAsset{
		Id: "test:1",
		Val1: "val1",
		Val2: []string{"val2.a", "val2.b"},
	}

	dl, err := faststorage.New("redis-14804.c10.us-east-1-2.ec2.cloud.redislabs.com:14804", "lT93etfPybYMoSnBfX71vUcDjadxbol3", "0", 10, 10,  1 * time.Second,  1 * time.Second, false, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("new dl")
	ctx := context.Background()
	key, err := dl.Redis.Put(ctx, &ts)
	if err != nil {
		fmt.Printf("put redis err: %v\n", err)
	}
	fmt.Printf("redis key: %+v\n", key)

	tsReply := faststorage.TestAsset{Id: "test:1"}

	err = dl.Redis.Get(ctx, &tsReply)
	if err != nil {
		fmt.Printf("get redis err: %v", err)
	}

	fmt.Printf("redis reply: %+v\n", tsReply)

	time.Sleep(2 * time.Minute)
	fmt.Println("goobye")
}
