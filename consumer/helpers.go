package main

import (
	// "fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect to database and get all websites
func connectDB()gorm.DB {
	const dsn = "host=localhost user=postgres password=9059015626 dbname=postgres port=5432"
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatal("unable to connect ot the database")
	}
	return *db
}

func getRegions(db *gorm.DB)[]string{
	var regions []Region
	db.Find(&regions)
	var regionList []string
	for _,region := range regions{
		regionList = append(regionList, region.Id)
	}
	return regionList
}

func getSiteId(db *gorm.DB,url string)string{
	var website Website
	db.Where("url = ?",url).First(&website)
	return website.Id
}

func setStatus(db *gorm.DB,siteId string, regionId string, status bool)bool{
	err := db.Create(&Status{
		SiteId: siteId,
		RegionId: regionId,
		Status: status,
		Time: time.Now(),
	})
	if err.Error != nil{
		return false
	}
	return true
}
