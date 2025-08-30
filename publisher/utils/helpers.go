package utils

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func CreateRedisGroup(client *redis.Client, streamName string, gorupName string)error{
	ctx := context.Background()
	response, err := client.XGroupCreateMkStream(ctx, streamName, gorupName, "$").Result()
	if err != nil {
		if err.Error() != "BUSYGROUP Consumer Group name already exists" {
			log.Fatalf("failed to create consumer group: %v", err)
			return errors.New("redis error")
		}
		fmt.Print("group already exist")
		
	}
	fmt.Print(response)
	return nil
}