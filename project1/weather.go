package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/fegoa89/zipcodes"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var partner int64
var requestData WeatherRequestStruct
var zipcodesDataset zipcodes.Zipcodes

func getRequest(context *gin.Context) *http.Request {
	var requestDataError Response

	city := context.Param("city")
	country := context.Param("country")
	zipCode := context.Param("zipcode")

	location, err2 := orgModel.FindLocation(zipCode)
	if err2 != nil {
		requestDataError.Status = http.StatusNotFound
		requestDataError.Error = "Not found, Zip Code not in database"
		requestDataError.Message = err2.Error()
		context.IndentedJSON(requestDataError.Status, requestDataError)
		return nil
	}

	if country != location.country || city != location.city {
		requestDataError.Status = http.StatusBadRequest
		requestDataError.Error = "Bad request, location does not match zip code"
		requestDataError.Message = fmt.Sprintf("[Country: %s, Zip Code Matching Country: %s, City: %s, Zip Code Matching City: %s, ZipCode: %s]", country, location.country, city, location.city, zipCode)
		context.IndentedJSON(requestDataError.Status, requestDataError)
		return nil
	}

	if len(country) != 2 {
		requestDataError.Status = http.StatusBadRequest
		requestDataError.Error = "Bad Request, invalid country code"
		requestDataError.Message = fmt.Sprintf("[City: %s, Country: %s, ZipCode: %s]", city, country, zipCode)
		context.IndentedJSON(requestDataError.Status, requestDataError)
		return nil
	}

	body := `{"city": "` + city + `", "country": "` + country + `", "zipCode": "` + zipCode + `"}`

	if err := json.Unmarshal([]byte(body), &requestData); err != nil {
		requestDataError.Status = http.StatusBadRequest
		requestDataError.Error = "Bad Request, format incorrect: " + body
		requestDataError.Message = err.Error()
		context.IndentedJSON(requestDataError.Status, requestDataError)
		return nil
	}

	orgStr := context.GetHeader("Org")

	orgInt, err := strconv.Atoi(orgStr)
	if err != nil {
		requestDataError.Status = http.StatusNotAcceptable
		requestDataError.Error = "Not Acceptable, org header failed converting to integer"
		requestDataError.Message = err.Error()
		context.IndentedJSON(requestDataError.Status, requestDataError)
		return nil
	}

	org, err2 := orgModel.FindPartner(orgInt)
	if err2 != nil {
		requestDataError.Status = http.StatusUnauthorized
		requestDataError.Error = "Unauthorized, org not found in database"
		requestDataError.Message = err.Error()
		context.IndentedJSON(requestDataError.Status, requestDataError)
		return nil
	}

	partner = org.Partner

	var text string
	var url string
	var auth string
	var method string
	if partner == 1 {
		text = `{"location": ` + body + `}`
		url = "https://us-central1-business-services-platform-dev.cloudfunctions.net/cnm-welcome-app-partner1"
		auth = "123"
		method = http.MethodPost
	} else if partner == 2 {
		text = body
		url = "https://us-central1-business-services-platform-dev.cloudfunctions.net/cnm-welcome-app-partner2"
		auth = "234"
		method = http.MethodPost
	} else if partner == 3 {
		var location string
		if country == "US" {
			location = zipCode
		} else {
			location = city + "," + country
		}
		text = ``
		currentTime := time.Now()
		url = `https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/` + location + "/" + currentTime.Format("2006-01-02") + `?unitGroup=metric&key=ZF5NDDZV63RC2RTRNY2SY6VQU&contentType=json`
		auth = ``
		method = http.MethodGet
	} else {
		requestDataError.Status = http.StatusNotFound
		requestDataError.Error = "Not Found, partner unknown"
		requestDataError.Message = "org.Partner"
		context.IndentedJSON(requestDataError.Status, requestDataError)
		return nil
	}

	textBytes := []byte(text)
	request, err := http.NewRequest(method, url, bytes.NewBuffer(textBytes))
	request.Header.Set("Authorization", auth)

	return request
}

func getResponseBody(response *http.Response) WeatherResponseStruct {

	body, _ := ioutil.ReadAll(response.Body)

	var responseBody WeatherResponseStruct

	if partner == 1 {
		var w WeatherResponseStruct1

		json.Unmarshal(body, &w)
		responseBody = WeatherResponseStruct{
			City:        w.Location.City,
			Country:     w.Location.Country,
			ZipCode:     w.Location.ZipCode,
			MaxTemp:     w.Temperature.Max,
			MinTemp:     w.Temperature.Min,
			CurrentTemp: w.Temperature.Current,
			Humidity:    w.Humidity,
			Pressure:    w.Pressure,
		}

	} else if partner == 2 {
		json.Unmarshal(body, &responseBody)
	} else if partner == 3 {
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		w := result["days"].([]interface{})[0].(map[string]interface{})
		maxtemp := w["tempmax"].(float64)
		mintemp := w["tempmin"].(float64)
		temp := w["temp"].(float64)
		humidity := int(w["humidity"].(float64))
		pressure := int(w["pressure"].(float64))

		responseBody = WeatherResponseStruct{
			City:        requestData.City,
			Country:     requestData.Country,
			ZipCode:     requestData.ZipCode,
			MaxTemp:     maxtemp,
			MinTemp:     mintemp,
			CurrentTemp: temp,
			Humidity:    humidity,
			Pressure:    pressure,
		}
	}

	return responseBody
}

func weather(context *gin.Context) {
	request := getRequest(context)

	if request == nil {
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	responseBody := getResponseBody(response)
	context.IndentedJSON(http.StatusOK, responseBody)
}
