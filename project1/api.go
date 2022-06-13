package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB
var orgModel OrgModel

func main() {
	db, err := GetMySQLDB()
	if err != nil {
		fmt.Println(err)
	} else {
		orgModel = OrgModel{
			Db: db,
		}
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	router := gin.Default()
	router.GET("/weather/city/:city/country/:country/zipcode/:zipcode", weather)
	router.Run("localhost:9090")
}
