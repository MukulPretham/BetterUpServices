package main

// TODOs
// handel consumer dead situation
// Notification feature.

import (
	"encoding/json"
	"log"

	"mukulpretham/betterUpPublisher/utils"
)

func main() {
	
	// Redis client created
	client := utils.CreateRedisClient("localhost:6379",0,"",2)

	if err := utils.CreateRedisGroup(client,"notifications","ntGroup"); err != nil{
		log.Fatal("redis error")
	}

	for {
		// Read messages form redis stremas via a consumer group
		res,err := utils.ReadXGroup(client,[]string{"websites",">"},"consumerGroup")
		if err != nil {
			log.Fatal(err)
		}

		for _, msg := range res[0].Messages {
			if currMesssage, ok := msg.Values["site"].(string); ok {
				var m map[string]string
				// Parsing to JSON.
				if err := json.Unmarshal([]byte(currMesssage), &m); err != nil {
					panic("error parsing string")
				}
				go WriteToDB(m["Url"],client,msg.ID)
			}
		}
	}
}

