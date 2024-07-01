package main

import "reflect"

// Key unique key created using City and country for the map
type Key struct {
	City    string
	Country string
	Geo     string
}

// GetUniqueKey take some of the fields and returns the key struct
func GetUniqueKey(city LocationData) Key {
	return Key{
		City:    city.Name,
		Country: city.Country,
		Geo:     city.Geo,
	}
}

// LocationData struct for the json data
type LocationData struct {
	Latitude     string `json:"latitude"`      // Latitude coordinate
	Longitude    string `json:"longitude"`     // Longitude coordinate
	Geo          string `json:"geo"`           // Geographical coordinates as a string
	Name         string `json:"City"`          // City name
	ProvinceIcon string `json:"province_icon"` // URL to the province icon (can be null)
	Province     string `json:"province"`      // Province name (empty if not applicable)
	CountryIcon  string `json:"country_icon"`  // URL to the country flag icon
	Country      string `json:"country"`       // Country name
}

func hardCheck(verifyData, inp LocationData) bool {
	return reflect.DeepEqual(verifyData, inp)
}

func (inp LocationData) hardValidate(verifyData LocationData) bool {
	return reflect.DeepEqual(verifyData, inp)
}

func (inp LocationData) basicValidate(verifyData LocationData) bool {
	if inp.Name == verifyData.Name && inp.Geo == verifyData.Geo && inp.Country == verifyData.Country {
		return true
	}
	return false
}
