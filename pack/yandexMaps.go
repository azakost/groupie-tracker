package pack

import (
	"strconv"
	"strings"
)

// Type for parsing Yandex Coordinations
type yandex struct {
	Response struct {
		GeoObjectCollection struct {
			FeatureMember []struct {
				GeoObject struct {
					Point struct {
						Pos string `json:"pos"`
					} `json:"Point"`
				} `json:"GeoObject"`
			} `json:"featureMember"`
		} `json:"GeoObjectCollection"`
	} `json:"response"`
}

// Getting location coordinates fo further operations
func yandexMaps(city string) ([]float64, string) {
	var pos yandex
	city = strings.ReplaceAll(city, " ", "%20")
	key := "ebc177a6-77e8-48e7-be65-78e871de82d5"
	epar := ParseAPI("https://geocode-maps.yandex.ru/1.x/?apikey="+key+"&format=json&geocode="+city, &pos)
	if epar != "Success" {
		return []float64{}, "Error with parsing Yandex API"
	}
	coor := pos.Response.GeoObjectCollection.FeatureMember[0].GeoObject.Point.Pos
	c := strings.Split(coor, " ")
	lon, elon := strconv.ParseFloat(c[0], 64)
	lat, elat := strconv.ParseFloat(c[1], 64)
	if elon != nil || elat != nil {
		return []float64{}, "Error with parsing coordinates to Float64 (" + coor + ")"
	}
	return []float64{lat, lon}, "Success"
}

var Latlon = make(map[string][]float64)

// Give location coordinates
func GetCoordinates(loc string) ([]float64, string) {

	// If we already found  coordinated for asked location just get it from stored variable 'coords'
	if val, ok := Latlon[loc]; ok {
		return val, "Success"
	} else {

		// if not - find it on yandex and store in variable 'coords'
		c, e := yandexMaps(loc)
		if e != "Success" {
			return c, e
		}
		Latlon[loc] = c
		return Latlon[loc], e
	}
}
