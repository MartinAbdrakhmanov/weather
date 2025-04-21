package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"weather/internal/model"
)

func GetForecast(city string, days int) model.ForecastResponse {
	apitoken := os.Getenv("WEATHER_API_TOKEN")
	city = strings.Replace(city, " ", "+", -1)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%v&q=%v&days=%v&aqi=no&alerts=yes&hour_fields=time,temp_c,is_day,condition", apitoken, city, days)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		return model.ForecastResponse{Status: "No matching location found."} // CHANGE OR  MATCHING ESSEX, UK
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	forecast := model.ForecastResponse{}
	json.Unmarshal(body, &forecast)
	return forecast

}
