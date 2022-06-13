package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "TECH21know!!"
	dbname   = "project1"
)

func GetMySQLDB() (db *sql.DB, err error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	return
}

func (orgModel OrgModel) FindPartner(id int) (Org, error) {
	rows, err := orgModel.Db.Query(`SELECT * FROM org_partners WHERE org_partners."Id" = ` + fmt.Sprint(id))
	if err != nil {
		return Org{}, err
	} else {
		var org Org
		for rows.Next() {
			var id int64
			var partner int64
			err2 := rows.Scan(&id, &partner)
			if err2 != nil {
				return Org{}, err2
			} else {
				org = Org{id, partner}
			}
		}
		return org, nil
	}
}

func (orgModel OrgModel) FindLocation(zip string) (Location, error) {
	rows, err := orgModel.Db.Query(`SELECT * FROM country_zip_city WHERE zip = '` + zip + `'`)
	if err != nil {
		return Location{}, err
	} else {
		var zipCountry Location
		for rows.Next() {
			var country string
			var zip string
			var city string
			err2 := rows.Scan(&country, &zip, &city)
			if err2 != nil {
				return Location{}, err2
			} else {
				zipCountry = Location{country, zip, city}
			}
		}
		return zipCountry, nil
	}
}
