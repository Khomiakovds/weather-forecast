package location

import "strings"

var (
	//"Saint Petersburg", "Kazan", "Nizhniy Novgorod", "Novosibirsk", "Samara"
	locations = map[string]Location{
		"moscow": {
			CityName:  "Moscow",
			Latitude:  55.751244,
			Longitude: 37.618423,
		},
		"saint petersburg": {
			CityName:  "Saint Petersburg",
			Latitude:  59.937500,
			Longitude: 30.308611,
		},
		"kazan": {
			CityName:  "Kazan",
			Latitude:  55.796391,
			Longitude: 49.108891,
		},
		"nizhniy novgorod": {
			CityName:  "Nizhniy Novgorod",
			Latitude:  56.296505,
			Longitude: 43.936058,
		},
		"novosibirsk": {
			CityName:  "Novosibirsk",
			Latitude:  55.018803,
			Longitude: 82.933952,
		},
		"samara": {
			CityName:  "Samara",
			Latitude:  53.241505,
			Longitude: 50.221245,
		},
	}
)

type Location struct {
	CityName  string
	Latitude  float64
	Longitude float64
}

func FindLocation(cityName string) *Location {
	if l, ok := locations[strings.ToLower(cityName)]; ok {
		return &l
	}
	return &Location{}
}
