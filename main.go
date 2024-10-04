package main

import (
	"10_1_simple_pipeline/forecast"
	"10_1_simple_pipeline/location"
	"10_1_simple_pipeline/predict_models"
	"fmt"
	"time"
)

func main() {

	done := make(chan struct{})

	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	requestsChan := forecast.RequestRandomGenerator(done)

	model := predict_models.NewModel1()

	cityWeather := weatherCalculation(done, requestsChan, model)

	fullInfo := cityCoordinates(done, cityWeather)

	print(fullInfo)
}

func weatherCalculation(done <-chan struct{}, requestsChan <-chan forecast.ForecastRequest, model *predict_models.Model) <-chan forecast.ForecastPrediction {
	weatherChan := make(chan forecast.ForecastPrediction)

	rateLimit := time.NewTicker(time.Minute / 60)

	go func() {
		defer close(weatherChan) 
		defer rateLimit.Stop()   

		for {
			select {
			case <-done: 
				return
			case currentRequest, ok := <-requestsChan:
				if !ok {
					
					return
				}
				select {
				case <-rateLimit.C: 
					
					weatherChan <- model.Predict(currentRequest)
				case <-done:
					return
				}
			}
		}
	}()

	return weatherChan
}

func cityCoordinates(done <-chan struct{}, weatherChan <-chan forecast.ForecastPrediction) chan string {
	coordinatesChan := make(chan string)

	go func() {
		defer close(coordinatesChan) 

		for {
			select {
			case <-done: 
				return
			case currentPrediction, ok := <-weatherChan:
				if !ok {
					
					return
				}
				loc := location.FindLocation(currentPrediction.Location)
				formattedOutput := fmt.Sprintf(
					"Location: %s (Lat: %.6f, Long: %.6f), Date: %s, Temp: %d°C, Humidity: %d%%, Wind Speed: %d km/h",
					loc.CityName,
					loc.Latitude,
					loc.Longitude,
					currentPrediction.Time.Format("2006-01-02 15:04:05"),
					currentPrediction.TemperatureCelsius,
					currentPrediction.HumidityPercent,
					currentPrediction.ProbabilityPercent)
				coordinatesChan <- formattedOutput
			}
		}
	}()

	return coordinatesChan
}

func print(resultsChan <-chan string) {
	for result := range resultsChan {
		fmt.Println(result)
	}
}
