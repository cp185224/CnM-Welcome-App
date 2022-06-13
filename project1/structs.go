package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type WeatherRequestStruct struct {
	City    string `json:"city"`
	Country string `json:"country"`
	ZipCode string `json:"zipCode"`
}

type WeatherResponseStruct struct {
	City        string  `json:"city"`
	Country     string  `json:"country"`
	ZipCode     string  `json:"zipCode"`
	CurrentTemp float64 `json:"currentTemp"`
	MinTemp     float64 `json:"minTemp"`
	MaxTemp     float64 `json:"maxTemp"`
	Humidity    int     `json:"humidity"`
	Pressure    int     `json:"pressure"`
}

type WeatherResponseStruct1 struct {
	Location struct {
		City    string `json:"city"`
		Country string `json:"country"`
		ZipCode string `json:"zipCode"`
	} `json:"location"`
	Temperature struct {
		Current float64 `json:"current"`
		Min     float64 `json:"min"`
		Max     float64 `json:"max"`
	} `json:"temperature"`
	Humidity int `json:"humidity"`
	Pressure int `json:"pressure"`
}

type Response struct {
	Status  int
	Error   string
	Message string
}

type Org struct {
	Id      int64
	Partner int64
}

type Location struct {
	country string
	zip     string
	city    string
}

type OrgModel struct {
	Db *sql.DB
}
