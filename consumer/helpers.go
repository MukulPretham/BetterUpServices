package main

import (
	// "fmt"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect to database and get all websites
func connectDB() gorm.DB {
	const dsn = "host=localhost user=postgres password=9059015626 dbname=postgres port=5432"
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatal("unable to connect ot the database")
	}
	return *db
}

func getRegions(db *gorm.DB) []string {
	var regions []Region
	db.Find(&regions)
	var regionList []string
	for _, region := range regions {
		regionList = append(regionList, region.Id)
	}
	return regionList
}

func getRegionId(db *gorm.DB,regionName string)string{
	var currRegion Region
	db.Where("name = ?",regionName).First(&currRegion)
	return currRegion.Id
}

func getSiteId(db *gorm.DB, url string) string {
	var website Website
	db.Where("url = ?", url).First(&website)
	return website.Id
}

func setStatus(db *gorm.DB, siteId string, regionId string, status bool) bool {
	err := db.Model(&Status{}).Where(`"siteId" = ?`, siteId).Update("status", status)
	if err.Error != nil {
		return false
	}
	return true
}

func setLatency(db *gorm.DB,siteId string, regionId string,latency float64){
	latencyRepot := Latency{
		Id: uuid.NewString(),
		SiteId: siteId,
		RegionId: regionId,
		Latency: latency,
		Time: time.Now(),
	}
	db.Create(&latencyRepot)
}

func fetch(url string)int{
	res, err := http.Get(fmt.Sprintf("https://%s", url))
	fmt.Print(err)
	if res.StatusCode == 200{
		return 200
	}else{
		return 0
	}
}

func WriteToDB(url string) {
	start := time.Now()
	var currLatency float64
	
	res := fetch(url)
	currLatency = float64(time.Since(start).Milliseconds())
	if res == 0 {
		db := connectDB()

		for _, regionId := range getRegions(&db) {
			setLatency(&db, getSiteId(&db, url), regionId, 404)
		}
		
		for _, regionId := range getRegions(&db) {
			setStatus(&db, getSiteId(&db, url), regionId, false)
		}
	} else {
		
		if res == 200 {
			db := connectDB()
			for _, regionId := range getRegions(&db) {
				setStatus(&db, getSiteId(&db, url), regionId, true)
			}
			for _, regionId := range getRegions(&db) {
				setLatency(&db, getSiteId(&db, url), regionId,currLatency)
			}
		} else {
			db := connectDB()
			for _, regionId := range getRegions(&db) {
				setStatus(&db, getSiteId(&db, url), regionId, false)
			}
			for _, regionId := range getRegions(&db) {
				setLatency(&db, getSiteId(&db, url), regionId,404)
			}
		}
	}
	fmt.Println("updated",url)
}
