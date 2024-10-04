package main

import (
	"10_1_simple_pipeline/forecast"
	"10_1_simple_pipeline/location"
	"10_1_simple_pipeline/predict_models"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	done := make(chan struct{})

	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	requestsChan := forecast.RequestRandomGenerator(done)

	wg.Add(2)

	cityWeather := weatherCalculation(done, requestsChan, &wg)

	fullInfo := cityCoordinates(done, cityWeather, &wg)

	print(fullInfo)

	wg.Wait()
}

func weatherCalculation(done <-chan struct{}, requestsChan <-chan forecast.ForecastRequest, wg *sync.WaitGroup) <-chan forecast.ForecastPrediction {
	weatherChan := make(chan forecast.ForecastPrediction)
	go func() {
		defer wg.Done()        
		defer close(weatherChan) 

		for currentRequest := range requestsChan {
			select {
			case <-done:
				return
			case weatherChan <- predict_models.NewModel1().Predict(currentRequest):
			}
		}
	}()
	return weatherChan
}

func cityCoordinates(done <-chan struct{}, requestsChan <-chan forecast.ForecastPrediction, wg *sync.WaitGroup) chan string {
	coordinatesChan := make(chan string)
	go func() {
		defer wg.Done()           
		defer close(coordinatesChan) 

		for currentRequest := range requestsChan {
			select {
			case <-done:
				return
			default:
				loc := location.FindLocation(currentRequest.Location)
				formattedOutput := fmt.Sprintf(
					"Location: %s (Lat: %.6f, Long: %.6f), Date: %s, Temp: %d°C, Humidity: %d%%, Wind Speed: %d km/h",
					loc.CityName,
					loc.Latitude,
					loc.Longitude,
					currentRequest.Time.Format("2006-01-02 15:04:05"),
					currentRequest.TemperatureCelsius,
					currentRequest.HumidityPercent,
					currentRequest.ProbabilityPercent)
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
