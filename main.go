package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
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
