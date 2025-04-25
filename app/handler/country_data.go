package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

var countries map[string]string

type Country struct {
	Name       string `json:"name"`
	Alpha2Code string `json:"alpha2Code"`
}

func init() {
	ReadCountryData()
}

func ReadCountryData() {
	countries = make(map[string]string)

	res, err := http.Get("https://restcountries.com/v2/all")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var countriesArray []Country
	err = json.NewDecoder(res.Body).Decode(&countriesArray)
	if err != nil {
		log.Fatal(err)
	}

	for _, country := range countriesArray {
		countries[country.Alpha2Code] = strings.ToUpper(country.Name)
	}
}

func GetCountries() map[string]string {
	return countries
}

func CheckCountry(ISO2 string, country string) int {
	val, ok := countries[ISO2]
	if !ok {
		return 1
	}

	if val != country {
		return 2
	}

	return 0
}
