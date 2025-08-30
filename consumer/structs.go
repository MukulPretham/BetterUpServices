package main

import "time"

type Status struct {
	Id       string `gorm:"column:id"`
	SiteId   string `gorm:"column:siteId"`
	RegionId string `gorm:"column:regionId"`
	Status   bool   `gorm:"column:status"`
}

type Latency struct {
	Id       string   `gorm:"column:id"`
	SiteId   string `gorm:"column:siteId"`
	RegionId string `gorm:"column:regionId"`
	Latency  float64 `gorm:"column:latency`
	Time     time.Time  `gorm:"column:time`
}

type Region struct {
	Id   string
	Name string
}

type Website struct {
	Id   string
	Name string
	Url  string
}

type StreamMsg struct {
	Id   string
	Name string
	Url  string
}

func (Region) TableName() string {
	return "Region"
}

func (Status) TableName() string {
	return "Status"
}

func (Website) TableName() string {
	return "Website"
}

func (Latency) TableName() string {
	return "Latency"
}
