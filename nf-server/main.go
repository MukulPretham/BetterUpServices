package main

import (
	// "encoding/json"
	// "context"
	"encoding/json"
	"fmt"
	"log"
	"mukulpretham/betterUpPublisher/utils"
)

func main() {
	client := utils.CreateRedisClient("localhost:6379", 0, "", 2)
	_, err := utils.CreateRedisGroup(client, "notifications", "notificationGroup")
	if err != nil {
		log.Fatal(err)
	}
	for {
		readRes, readErr := utils.ReadXGroup(client, []string{"notifications", ">"}, "notificationGroup")
		if readErr != nil {
			log.Fatal(readErr)
		}
		currMessage := readRes[0].Values["site"].(string)
		m :=  make(map[string]string)
		if err := json.Unmarshal([]byte(currMessage),&m); err==nil{
			fmt.Println(m["siteId"])
		}
	}
}
