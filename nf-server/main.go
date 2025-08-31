package main

import (
	"fmt"
	"log"
	"mukulpretham/betterUpPublisher/utils"
)

func main(){
	client := utils.CreateRedisClient("localhost:6379",0,"",2)
	res,err := utils.CreateRedisGroup(client,"notifications","notificationGroup")
	if err != nil{
		log.Fatal(err)
	}
	fmt.Print(res)
}