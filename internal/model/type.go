package model

type ForecastResponse struct {
	Location struct {
		Name      string `json:"name"`
		Country   string `json:"country"`
		LocalTime string `json:"localtime"`
	} `json:"location"`

	Current struct {
		TempC       float64 `json:"temp_c"`
		Wind_kph    float64 `json:"wind_kph"`
		Feelslike_c float64 `json:"feelslike_c"`
		Humidity    float64 `json:"humidity"`
		Condition   struct {
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
	Status     string
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

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type CityResponse struct {
	City   string `json:"city"`
	Status string `json:"status"`
}

type FormData struct {
	Surname     string
	Name        string
	Patronymic  string
	Approval    string
	Suggestions string
	Email       string
}
