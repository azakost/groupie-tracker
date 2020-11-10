package pack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// Standard parsing API function
func ParseAPI(url string, d interface{}) string {
	got, eget := http.Get(url)
	body, ered := ioutil.ReadAll(got.Body)
	emar := json.Unmarshal(body, d)
	if eget != nil || ered != nil || emar != nil {
		return "Error with parsing " + url
	}
	return "Success"
}

// Type for parsing "https://groupietrackers.herokuapp.com/api/artists"
type tmpart []struct {
	Id        int      `json:"id"`
	Image     string   `json:"image"`
	Name      string   `json:"name"`
	Members   []string `json:"members"`
	Establish int      `json:"creationDate"`
	Album     string   `json:"firstAlbum"`
}

// Type for parsing "https://groupietrackers.herokuapp.com/api/relation"
type tmprel struct {
	Index []struct {
		Dl map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

type Bands []struct {
	Id        int
	Name      string
	Image     string
	Members   []string
	Establish int
	Album     int64
	Concerts  []Event
}

type Event struct {
	Place  string
	Coords []float64
	Dates  []int64
}

// setting global variable for use to process data
var Artists Bands

func GrabData(url string) string {

	// Read base.json file if exist - if not parse, organize and create it
	out, ered := ioutil.ReadFile("pack/base.json")
	if ered == nil {
		// Unmarshal existings file and assign to global variable
		json.Unmarshal([]byte(out), &Artists)
		return "Success"
	} else {

		fmt.Println("Parsing API...")
		// Parse raw data form given API
		var rawArt tmpart
		var rawRel tmprel
		eart := ParseAPI(url+"artists", &rawArt)
		erel := ParseAPI(url+"relation", &rawRel)

		// Return error message if parsing failed
		if eart != "Success" {
			return eart
		} else if erel != "Success" {
			return erel
		} else {
			fmt.Println("Putting parsed data to organized structure and getting coordinates from Yandex API...")
			// Put raw data to organized struct 'Artists'
			b := make(Bands, len(rawArt))
			for i, val := range rawArt {
				b[i].Id = val.Id
				b[i].Name = val.Name
				b[i].Image = val.Image
				b[i].Members = val.Members
				b[i].Establish = val.Establish
				b[i].Album = convertDate([]string{val.Album})[0]
				csrt := make([]Event, len(rawRel.Index[i].Dl))
				n := 0
				for loc, dat := range rawRel.Index[i].Dl {
					loc = cleanLocaton(loc)
					csrt[n].Place = loc
					cord, e := GetCoordinates(loc)
					if e != "Success" {
						return e
					}
					csrt[n].Coords = cord
					csrt[n].Dates = convertDate(dat)
					n++
				}
				b[i].Concerts = csrt
			}

			// Assign organized data to a global variable
			Artists = b

			// Create json file for the next time
			createJson(b, "pack/base.json")
		}
		fmt.Println("Creation base.json is done!")
		return "Success"
	}
}

// Function to create base.json file
func createJson(base interface{}, dir string) {
	b, _ := json.Marshal(base)
	f, _ := os.Create(dir)
	ioutil.WriteFile(dir, b, 0644)
	defer f.Close()
}

// Edit date to make it possible for further parsing to time.Time type
const layoutISO = "2006-01-02"

func convertDate(s []string) []int64 {
	var res []int64
	for _, r := range s {
		arr := strings.Split(r, "-")
		tmp := make([]string, 3)
		tmp[0] = arr[2]
		tmp[1] = arr[1]
		tmp[2] = arr[0]
		msec, _ := DateToNano(strings.Join(tmp, "-"))
		res = append(res, msec)
	}
	return res
}

func DateToNano(date string) (int64, error) {
	d, e := time.Parse(layoutISO, date)
	return d.UnixNano(), e
}

func NanoToDate(nano int64) string {
	d := time.Unix(0, nano).Format(layoutISO)
	return d
}

// Make Location titles look fancy
func cleanLocaton(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	arr := strings.Split(s, "-")
	if len(arr) == 1 {
		return strings.Title(arr[0])
	}
	city := strings.Title(arr[0])
	country := ""
	if arr[1] == "usa" || arr[1] == "uk" {
		country = strings.ToUpper(arr[1])
	} else {
		country = strings.Title(arr[1])
	}
	return city + ", " + country
}
