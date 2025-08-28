package main

type Website struct{
	Id string
	Name string
	Url string
}

func (Website) TableName() string {
	return "Website"
}