package main

import (
	"context"
	"encoding/json"

	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"math/rand"
	"net/http"
)

// import (
// 	"fmt"
// )

func main() {

	// Redis client created
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	ctx := context.Background()
	for {
		// Read messages form redis stremas via a consumer group
		res, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Streams:  []string{"websites", ">"},
			Group:    "consumerGroup",
			Consumer: fmt.Sprintf("consumerGroup%d", rand.Intn(1000)),
			Block:    0,
			Count:    1,
		}).Result()
		if err != nil {
			panic(err)
		}
		// res structure
		// [{
		// 	  Stream
		// 	  Messages	: [{
		// 	     ID
		// 	     Values : {
		//		   // data stores in redis
		//       }
		//    }]
		// }]

		for _, msg := range res[0].Messages {
			if currMesssage, ok := msg.Values["site"].(string); ok {
				var m map[string]string
				if err := json.Unmarshal([]byte(currMesssage), &m); err != nil {
					panic("error parsing string")
				}

				fmt.Println(m["Url"])
				if res, err := http.Get(fmt.Sprintf("https://%s", m["Url"])); err != nil {
					log.Fatal("error while fetch")
					panic("")
				} else {
					if res.StatusCode == 200 {

					} else {

					}
				}
			}
			client.XAck(ctx, "websites", "consumerGroup", msg.ID)
		}

		//Check website ping
		//Update database
		//Acknowledge the message

	}
}

// func main(){
// 	db := connectDB()
// 	fmt.Println(getRegions(&db))
// 	for _,regionId := range getRegions(&db){
// 		setStatus(&db,getSiteId(&db,"www.google.com"),regionId,false)
// 	}
// }

type StreamMsg struct {
	Id   string
	Name string
	Url  string
}
