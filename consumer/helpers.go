package main

import (
	// "fmt"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
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

func getRegionId(db *gorm.DB, regionName string) (string,error) {
	var currRegion Region
	err := db.Where("name = ?", regionName).First(&currRegion)
	if err.Error == nil {
		return currRegion.Id,nil
	}else{
		return "",errors.New("region wich was passed via env variable is not a valid region")
	}
}

func getSiteId(db *gorm.DB, url string) string {
	var website Website
	db.Where("url = ?", url).First(&website)
	return website.Id
}

func setStatus(db *gorm.DB, siteId string, regionId string, status bool) bool {
	err := db.Model(&Status{}).Where(`"siteId" = ? AND "regionId" = ?`, siteId,regionId).Update("status", status)
	if err.Error != nil {
		return false
	}
	return true
}

func setLatency(db *gorm.DB, siteId string, regionId string, latency float64) {
	latencyRepot := Latency{
		Id:       uuid.NewString(),
		SiteId:   siteId,
		RegionId: regionId,
		Latency:  latency,
		Time:     time.Now(),
	}
	db.Create(&latencyRepot)
}

func fetch(url string) int {
	res, err := http.Get(fmt.Sprintf("https://%s", url))
	fmt.Print(err)
	if res.StatusCode == 200 {
		return 200
	} else {
		return 0
	}
}

func WriteToDB(url string) {
	var currLatency float64

	start := time.Now()
	
	db := connectDB()
	res := fetch(url)

	currLatency = float64(time.Since(start).Milliseconds())

	env := os.Getenv("REGION")
	currRegionId,err := getRegionId(&db, env)
	if err != nil{
		log.Fatal(err)
	}

	if res == 200 {
		setStatus(&db, getSiteId(&db, url), currRegionId, true)
		setLatency(&db, getSiteId(&db, url), currRegionId, currLatency)
	} else {
		setLatency(&db, getSiteId(&db, url), currRegionId, 404)
		setStatus(&db, getSiteId(&db, url), currRegionId, false)
	}
	fmt.Println("updated", url)
}
