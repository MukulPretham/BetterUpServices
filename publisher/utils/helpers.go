package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"

	"github.com/redis/go-redis/v9"
)

func CreateRedisClient(Addr string,DB int, Password string, Protocol int) *redis.Client{
	return redis.NewClient(&redis.Options{
		Addr: Addr,
		DB: DB,
		Password: Password,
		Protocol: Protocol,
	})
}

func CreateRedisGroup(client *redis.Client, streamName string, gorupName string)error{
	ctx := context.Background()
	_, err := client.XGroupCreateMkStream(ctx, streamName, gorupName, "$").Result()
	if err != nil {
		if err.Error() != "BUSYGROUP Consumer Group name already exists" {
			log.Fatalf("failed to create consumer group: %v", err)
			return errors.New("redis error")
		}
		fmt.Print("group already exist")	
	}
	return nil
}

func ReadXGroup(client *redis.Client ,stream []string,gorupName string,)([]redis.XStream,error){
	ctx := context.Background()
	res, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Streams:  stream,
		Group:    gorupName,
		Consumer: fmt.Sprintf("%s - %d", gorupName,rand.Intn(1000)),
		Block:    0,
		Count:    1,
	}).Result()
	if err != nil{
		return nil,fmt.Errorf("redis XGroupRead error: %w",err)
	}
	return res,nil
}