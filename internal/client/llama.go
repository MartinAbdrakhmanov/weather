package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"weather/internal/model"
)

func GetClothingSuggestion(forecast model.ForecastResponse) string {
	prompt := generatePrompt(forecast)

	request := model.OllamaRequest{
		Model:  "llama3",
		Prompt: prompt,
		Stream: true,
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		log.Fatal("Marshall error:", err)
	}

	req, err := http.NewRequest("POST", os.Getenv("LLAMA_API_URL"), bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatal("Request error:", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Ollama request error:", err)
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
			log.Fatal("Read response error:", err)
		}
		if len(line) == 0 {
			continue
		}

		var chunk model.OllamaResponse
		if err := json.Unmarshal(line, &chunk); err != nil {
			log.Printf("Unmarshall error: %v\n%s", err, line)
			continue
		}

		result += chunk.Response
	}

	return result
}

func generatePrompt(forecast model.ForecastResponse) string {
	current_temperature := forecast.Current.TempC
	current_feels := forecast.Current.Feelslike_c
	current_humid := forecast.Current.Humidity
	current_windspd := forecast.Current.Wind_kph
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
	prompt := fmt.Sprintf(
		`I'm trying to decide what to wear based on the current and upcoming weather.

	Right now in %s, it's %.1f°C but feels like %.1f°C. The sky is %s.
	Humidity is at %v%% and the wind speed is %.1f kph.
	Local time is %s.

	Throughout the day, weather conditions may include: %v.
	The temperature will range from a low of %.1f°C to a high of %.1f°C.

	Given these conditions, what kind of clothing would be appropriate to stay comfortable and prepared?
	Keep response in range of 150 word and please try not to use extra \n between paragraphs if its not needed (MAX 2 PARAGRAPHS ALLOWED)`,
		location,
		current_temperature,
		current_feels,
		current_condition,
		current_humid,
		current_windspd,
		current_time,
		conditions,
		mintemp,
		maxtemp,
	)
	return prompt
}
