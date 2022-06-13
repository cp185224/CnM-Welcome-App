package partner2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type PartnerResponse struct {
	Main struct {
		Temp     float64 `json:"temp"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
}

type WeatherRequest struct {
	City    string `json:"city"`
	Country string `json:"country"`
	ZipCode string `json:"zipCode"`
}

type WeatherResponse struct {
	City        string  `json:"city"`
	Country     string  `json:"country"`
	ZipCode     string  `json:"zipCode"`
	CurrentTemp float64 `json:"currentTemp"`
	MinTemp     float64 `json:"minTemp"`
	MaxTemp     float64 `json:"maxTemp"`
	Humidity    int     `json:"humidity"`
	Pressure    int     `json:"pressure"`
}

type ErrorResponse struct {
	Msg    string `json:"msg"`
	Status int    `json:"status"`
}

func Weather(w http.ResponseWriter, r *http.Request) {
	log.Println("running")
	var wr WeatherRequest

	if r.Header.Get("Authorization") != "234" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{
			Msg:    "Invalid Authorization.",
			Status: http.StatusUnauthorized,
		})
		return
	}

	err := json.NewDecoder(r.Body).Decode(&wr)

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?zip=%s,"+
		"%s&appid=be8d8edba8be14f65dcb483dbcdbd653&units=imperial", wr.ZipCode, wr.Country)
	resp, err := http.Get(url)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(string(body))
		var pResp PartnerResponse
		json.Unmarshal(body, &pResp)

		wResp := &WeatherResponse{
			City:        wr.City,
			Country:     wr.Country,
			ZipCode:     wr.ZipCode,
			MaxTemp:     pResp.Main.TempMax,
			MinTemp:     pResp.Main.TempMin,
			CurrentTemp: pResp.Main.Temp,
			Humidity:    pResp.Main.Humidity,
			Pressure:    pResp.Main.Pressure,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(wResp)
	} else if resp.StatusCode == http.StatusNotFound {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Msg:    "Could not find location.",
			Status: http.StatusNotFound,
		})
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Msg:    "Server Error.",
			Status: http.StatusInternalServerError,
		})
	}
}
