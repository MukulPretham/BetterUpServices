package main

import (
	"encoding/json"
	"fmt"
	"log"
	
	"time"

	"mukulpretham/betterUpPublisher/utils"

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
	client := utils.CreateRedisClient("localhost:6379",0,"",2)

	// Create redis consumerGroup and a stream
	res,err := utils.CreateRedisGroup(client,"websites","consumerGroup")
	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf(res)
	for {
		func(db *gorm.DB,client *redis.Client){
			var cueeWebsites []utils.Website
			db.Find(&cueeWebsites)

			ctx := context.Background()
			
			for _,rec := range cueeWebsites{
				data, err := json.Marshal(rec)
				if err != nil {
					log.Println("Failed to marshal:", err)
					continue
				}
				response,err := client.XAdd(ctx,&redis.XAddArgs{
					Stream: "websites",
					Values: map[string]any{
						"site": string(data),
					},
				}).Result()
				fmt.Print(response)
			}

		}(db,client)
		fmt.Println("iteration completed")
		time.Sleep(3 * time.Second)
	}
}
