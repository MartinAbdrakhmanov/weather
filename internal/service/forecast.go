package service

import (
	"log"
	"strconv"
	"time"
	"weather/internal/model"
)

func FilterForecastByTime(forecast *model.ForecastResponse) {
	locTime, err := time.Parse("2006-01-02 15:04", forecast.Location.LocalTime)
	if err != nil {
		log.Println("Ошибка парсинга времени:", err)
		return
	}

	cutoff := locTime.Add(18 * time.Hour)
	var filteredHours model.HourD
	for i := range []int{0, 1} { //forecast.Forecast.Forecastday

		for _, hour := range forecast.Forecast.Forecastday[i].Hour {
			hourTime, err := time.Parse("2006-01-02 15:04", hour.Time)
			if err != nil {
				log.Println("Ошибка парсинга времени для часа:", err)
				continue
			}

			if hourTime.After(locTime) && hourTime.Before(cutoff) {
				hour.Time = hour.Time[11:]
				filteredHours = append(filteredHours, hour)
			}
		}

	}
	forecast.HourWindow = filteredHours

	for i, day := range forecast.Forecast.Forecastday { //

		date, err := time.Parse("2006-01-02", day.Date)
		if err != nil {
			log.Println("Ошибка парсинга времени для дня:", err)
			continue
		}
		forecast.Forecast.Forecastday[i].WeekDay = strconv.Itoa(date.Day()) + " " + date.Month().String() + "\n" + date.Weekday().String()
	}

}
