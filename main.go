package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {

	max := 10
	prefix := "test:"
	rand.Seed(time.Now().UnixNano())

	ctx := context.Background()
	url := os.Args[1]
	username := os.Args[2]
	password := os.Args[3]
	now := time.Now().Format(time.RFC3339)
	fmt.Printf("Connect to %s as %s at %s\n", url, username, now)
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        []string{url + ":6379"},
		Username:     username,
		Password:     password,
		MaxRedirects: 10,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})
	defer rdb.Close()

	fmt.Println("---")

	var keys []string
	for i := 0; i < max; i++ {
		s := fmt.Sprintf("%024d",rand.Uint64())
		se := base64.StdEncoding.EncodeToString([]byte(s))
		key := prefix + se
		fmt.Printf("set %s, %s\n", key, now)
		err := rdb.Set(ctx, key, now, 0).Err()
		if err != nil {
			fmt.Println(err.Error())
		}
		keys = append(keys,key)
	}

	fmt.Println("---")
	for i:=0; i < max; i++ {
		r := rdb.Get(ctx, keys[i])
		if r.Err()!= nil {
			fmt.Println(r.Err())
		}
		fmt.Printf("get %s, %s\n", keys[i], r.Val())
	}

	fmt.Println("---")
	v := rdb.MGet(ctx, keys...)
	if v.Err() != nil {
		fmt.Println(v.Err())
	}
	fmt.Printf("%+v\n",v.Val()...)

	fmt.Println("---")
	fmt.Printf("SET a: %s\n", now)
	err := rdb.Set(ctx, "a", now, 0).Err()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("---")
	val, err := rdb.Get(ctx, "a").Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("GET a: %s\n", val)
}
