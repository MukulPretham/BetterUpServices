package main

import (
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
			fmt.Println(cueeWebsites)

			ctx := context.Background()
			for _,rec := range cueeWebsites{
				response,err = client.XAdd(ctx,&redis.XAddArgs{
					Stream: "websites",
					Values: map[string]interface{}{
						"site": rec,
					},
					ID: "0-2",
				}).Result()
			}

		}(db,client)
		
		time.Sleep(3 * time.Second)
	}
}
