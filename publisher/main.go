package main

import (
	"encoding/json"
	"fmt"
	"log"
	
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"context"

	"github.com/redis/go-redis/v9"
)

func main() {

	//Connect to database and get all websites
	dsn := "host=localhost user=postgres password=9059015626 dbname=postgres port=5432"
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatal("unable to connect ot the database")
	}
	//Creating redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	// Create redis consumerGroup and a stream
	ctx := context.Background()
	response, err := client.XGroupCreateMkStream(ctx, "websites", "consumerGroup", "$").Result()
	if err != nil {
		if err.Error() != "BUSYGROUP Consumer Group name already exists" {
			log.Fatalf("failed to create consumer group: %v", err)
		}
		fmt.Print("group already exist")
	}
	fmt.Print(response)
	
	for {
		func(db *gorm.DB,client *redis.Client){
			var cueeWebsites []Website
			db.Find(&cueeWebsites)

			ctx := context.Background()
			
			for _,rec := range cueeWebsites{
				data, err := json.Marshal(rec)
				if err != nil {
					log.Println("Failed to marshal:", err)
					continue
				}
				response,err = client.XAdd(ctx,&redis.XAddArgs{
					Stream: "websites",
					Values: map[string]any{
						"site": string(data),
					},
				}).Result()
			}

		}(db,client)
		fmt.Println("iteration completed")
		time.Sleep(3 * time.Second)
	}
}
