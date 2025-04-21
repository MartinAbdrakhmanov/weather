package api

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"weather/internal/client"
	"weather/internal/service"

	"github.com/gorilla/mux"
)


func GetWeather(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	city := vars["city"]
	if city == "" {
		city = "Moscow"
	}
	log.Printf("Got city for getWeather: %v\n", city)
	forecast := client.GetForecast(city, 8) //
	if forecast.Status == "No matching location found." {
		forecast.Location.Name = "Location isn't found"
	} else {
		service.FilterForecastByTime(&forecast)
	}

	tmpl := template.Must(template.New("layout.html").Funcs(template.FuncMap{
		"json": jsonMarshal,
	}).ParseFiles("templates/layout.html", "templates/weather.html"))

	err := tmpl.ExecuteTemplate(w, "layout.html", forecast)
	if err != nil {
		http.Error(w, "Ошибка отображения шаблона", http.StatusInternalServerError)
	}

}

func jsonMarshal(v interface{}) string {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(jsonData)
}
