package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"github.com/russross/blackfriday/v2"
)

func getReport(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join("data", "Report.md")

	report, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Не удалось прочитать файл", http.StatusInternalServerError)
	}

	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	tmpl := template.Must(template.New("layout.html").Funcs(funcMap).ParseFiles("templates/layout.html", "templates/report.html"))
	reportb := string(blackfriday.Run(report))
	data := struct {
		Title      string
		ReportHTML string
	}{
		Title:      "Report",
		ReportHTML: reportb,
	}
	err = tmpl.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		log.Fatal(err)
	}

}

func getForm(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.New("layout.html").ParseFiles("templates/layout.html", "templates/form.html"))

	err := tmpl.Execute(w, "layout.html")
	if err != nil {
		log.Fatal(err)
	}

}

type FormData struct {
	Surname     string
	Name        string
	Patronymic  string
	Approval    string
	Suggestions string
	Email       string
}

func submitForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	data := FormData{
		Surname:     r.FormValue("surname"),
		Name:        r.FormValue("name"),
		Patronymic:  r.FormValue("patronymic"),
		Approval:    r.FormValue("approval"),
		Suggestions: r.FormValue("suggestions"),
		Email:       r.FormValue("email"),
	}

	fmt.Fprintf(w, "Полученные данные:\n%+v", data)

}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

type ForecastResponse struct {
	Location struct {
		Name      string `json:"name"`
		Country   string `json:"country"`
		LocalTime string `json:"localtime"`
	} `json:"location"`

	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
		} `json:"condition"`
	} `json:"current"`

	Forecast struct {
		Forecastday []struct {
			Date    string `json:"date"`
			WeekDay string
			Day     struct {
				MaxtempC  float64 `json:"maxtemp_c"`
				MintempC  float64 `json:"mintemp_c"`
				Condition struct {
					Text string `json:"text"`
					Icon string `json:"icon"`
				} `json:"condition"`
			} `json:"day"`

			Hour HourD `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
	Suggestion string
	HourWindow HourD
}

type HourD []struct {
	Time      string  `json:"time"`
	TempC     float64 `json:"temp_c"`
	IsDay     int     `json:"is_day"`
	Condition struct {
		Text string `json:"text"`
		Icon string `json:"icon"`
	} `json:"condition"`
}

func getForecast(city string, days int) ForecastResponse {
	apitoken := "86fb988a827648b4a25140324251604"
	city = strings.Replace(city, " ", "+", -1)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%v&q=%v&days=%v&aqi=no&alerts=yes&hour_fields=time,temp_c,is_day,condition", apitoken, city, days)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	forecast := ForecastResponse{}
	json.Unmarshal(body, &forecast)
	return forecast

}

func filterForecastByTime(forecast *ForecastResponse) {
	locTime, err := time.Parse("2006-01-02 15:04", forecast.Location.LocalTime)
	if err != nil {
		log.Println("Ошибка парсинга времени:", err)
		return
	}

	cutoff := locTime.Add(18 * time.Hour)

	for i := range []int{0, 1} { //forecast.Forecast.Forecastday

		var filteredHours HourD

		for _, hour := range forecast.Forecast.Forecastday[i].Hour {
			hourTime, err := time.Parse("2006-01-02 15:04", hour.Time)
			if err != nil {
				log.Println("Ошибка парсинга времени для часа:", err)
				continue
			}

			if hourTime.After(locTime) && hourTime.Before(cutoff) {
				hour.Time = hour.Time[12:]
				filteredHours = append(filteredHours, hour)
			}
		}

		forecast.HourWindow = filteredHours
	}

	for i, day := range forecast.Forecast.Forecastday { //

		date, err := time.Parse("2006-01-02", day.Date)
		if err != nil {
			log.Println("Ошибка парсинга времени для дня:", err)
			continue
		}
		forecast.Forecast.Forecastday[i].WeekDay = strconv.Itoa(date.Day()) + " " + date.Month().String() + "\n" + date.Weekday().String()
	}

}

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // No password set
	DB:       0,  // Use default DB
	Protocol: 2,  // Connection protocol
})

func getWeather(w http.ResponseWriter, r *http.Request) {
	forecast := getForecast("Saint-Petersburg", 10) //
	filterForecastByTime(&forecast)

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

func getSuggestion(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		log.Fatal("city is required", http.StatusBadRequest)
		return
	}
	suggestion, err := rdb.Get(ctx, city).Result()
	if err == redis.Nil {
		forecast := getForecast(city, 2)
		filterForecastByTime(&forecast)
		suggestion = getClothingSuggestion(forecast)
		err = rdb.Set(ctx, city, suggestion, 1*time.Hour).Err() // change me
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"suggestion": suggestion,
	})
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"` // true — чтобы понимать, что будем стримить
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func getClothingSuggestion(forecast ForecastResponse) string {
	prompt := generatePrompt(forecast)

	request := OllamaRequest{
		Model:  "llama3",
		Prompt: prompt,
		Stream: true,
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		log.Fatal("Ошибка при маршалинге запроса:", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:11434/api/generate", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatal("Ошибка при создании запроса:", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Ошибка при выполнении запроса к Ollama:", err)
	}
	defer resp.Body.Close()

	var result string
	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Ошибка при чтении ответа:", err)
		}
		if len(line) == 0 {
			continue
		}

		var chunk OllamaResponse
		if err := json.Unmarshal(line, &chunk); err != nil {
			log.Printf("Ошибка при разборе чанка: %v\n%s", err, line)
			continue
		}

		result += chunk.Response
	}

	return result
}

func generatePrompt(forecast ForecastResponse) string {
	// log.Printf("Пришел прогноз %v\n", forecast)
	current_temperature := forecast.Current.TempC
	current_condition := forecast.Current.Condition.Text
	location := forecast.Location.Name + " " + forecast.Location.Country
	current_time := forecast.Location.LocalTime
	conditions := []string{}
	for i := range forecast.Forecast.Forecastday {
		for _, hour := range forecast.Forecast.Forecastday[i].Hour {
			if !slices.Contains(conditions, hour.Condition.Text) {
				conditions = append(conditions, hour.Condition.Text)
			}
		}
	}
	maxtemp := forecast.Forecast.Forecastday[0].Day.MaxtempC
	mintemp := forecast.Forecast.Forecastday[0].Day.MintempC
	prompt := fmt.Sprintf("Скажи, какую одежду стоит надеть, если температура %.1f°C и состояние погоды: %s. Мы находимся в %s, Текущее время: %v. В течение дня будет погода: %v, максимальная температура будет достигать %v, и падать до %v. ОТВЕЧАЙ СТРОГО НА РУССКОМ ЯЗЫКЕ",
		current_temperature, current_condition, location, current_time, conditions, maxtemp, mintemp)
	// log.Printf("Промпт: %v \n", prompt)
	return prompt
}

func main() {
	router := mux.NewRouter()

	router.Use(logRequest)
	router.HandleFunc("/", getWeather)
	router.HandleFunc("/report", getReport)
	router.HandleFunc("/form", getForm)
	router.HandleFunc("/submit", submitForm).Methods(http.MethodPost)
	router.HandleFunc("/suggestion", getSuggestion)

	log.Fatal(http.ListenAndServe(":8080", router))
}
