package main

// TODOs
// handel consumer dead situation
// Notification feature.

// import (
// 	"fmt"
// 	"mukulpretham/betterUpConsumer/helpers"
// 	"os"

// 	"github.com/joho/godotenv"
// )

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"mukulpretham/betterUpPublisher/utils"

	"mukulpretham/betterUpConsumer/helpers"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	fmt.Print(os.Getenv("SmptIP"))
	// Redis client created
	client := utils.CreateRedisClient("localhost:6379",0,"",2)

	res,err := utils.CreateRedisGroup(client,"notifications","ntGroup");
	if err != nil{
		log.Fatal(err)
	}
	fmt.Print(res)

	for {
		// Read messages form redis stremas via a consumer group
		res,err := utils.ReadXGroup(client,[]string{"websites",">"},"consumerGroup")
		if err != nil {
			log.Fatal(err)
		}

		msg := res[0]

		if currMesssage, ok := msg.Values["site"].(string); ok {
			var m map[string]string
			// Parsing to JSON.
			if err := json.Unmarshal([]byte(currMesssage), &m); err != nil {
				panic("error parsing string")
			}
			go helpers.WriteToDB(m["Url"],client,msg.ID)
		}
	}
}

// func main(){
// 	godotenv.Load(".env")
// 	fmt.Print(os.Getenv("REGION"))
// 	helpers.SendMain([]string{"schmunna@gmail.com"},"hi")
// }

