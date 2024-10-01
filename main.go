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
	go func(){
		time.Sleep(5*time.Second)
		done <- struct{}{}
	}()
	stopCh := make(chan struct{})
	requestsChan := forecast.RequestRandomGenerator(stopCh)
	
	
cityWether := weathercalculation(done,requestsChan)
fullInfo := cityСoordinates(done,cityWether )
print(fullInfo)
	

}


func weathercalculation(done <- chan struct{}, requestsChan<-chan forecast.ForecastRequest)<-chan forecast.ForecastPrediction{
	wetherTemplate := make(chan forecast.ForecastPrediction)
	go func(){
		for currentRequst := range requestsChan{
			defer close(wetherTemplate)
			select{
			case<-done:
				return
			case wetherTemplate <- predict_models.NewModel1().Predict(currentRequst):
				
	
			}
			
		}
	}()
	
	return wetherTemplate
}
func cityСoordinates(done<- chan struct{},requestsChan<-chan forecast.ForecastPrediction)chan string{
	coordinateWetherChan := make(chan string)
	go func(){
		defer close(coordinateWetherChan)
		select{
		case <-done:
			return
		default:
			for currentRequst := range requestsChan{
				locatName := location.FindLocation(currentRequst .Location)
				formattedOutput := fmt.Sprintf(
					"Location: %s (Lat: %.6f, Long: %.6f), Date: %s, Temp: %d°C, Humidity: %d%%, Wind Speed: %d km/h",
					locatName.CityName,                               
					locatName.Latitude,                           
					locatName.Longitude,                          
					currentRequst .Time.Format("2006-01-02 15:04:05"),     
					currentRequst .TemperatureCelsius,                            
					currentRequst .HumidityPercent ,                               
					currentRequst .ProbabilityPercent)   
				coordinateWetherChan <- formattedOutput
			
			}

		}
		
	}()
	return coordinateWetherChan

}
func print(requestsChan<-chan string ){
	for i := range requestsChan{
		fmt.Println(i)
	}
}






