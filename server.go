package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"text/template"

	. "./pack"
)

func main() {

	// 1. Grab all needed data before starting listening http port
	// Result of this Grabbing is "Artists" global variable that can be used for further operations
	err := GrabData("https://groupietrackers.herokuapp.com/api/")
	if err != "Success" {
		fmt.Println(err)
	} else {

		// Setting a file server to hold js & css files there
		assets := http.FileServer(http.Dir("html/assets"))
		http.Handle("/assets/", http.StripPrefix("/assets/", assets))

		// Handle request
		fmt.Println("Listening http://localhost:8080")
		http.HandleFunc("/", request)
		err := http.ListenAndServe(":8080", nil)

		fmt.Println(err)
	}

}

func request(w http.ResponseWriter, r *http.Request) {

	switch r.URL.Path {

	case "/api": // This case is only for javascript render (/html/assets/list.js)

		// 1. Store requests into spacial struct
		requests := WriteRequests(r)

		// 2. Make a filtered by requests list of bands + get search suggestions list
		bands, suggestions := BandList(requests)

		// 3. Sort bands if sorting request is set
		sorted := SortBands(bands, requests)

		// 4. Paginate all data + get numbers of bands, pages and requested page
		paginated, pages, page, found := Paginate(sorted, requests)

		// 5. Construct a final struct with list of bands including previously found search suggestions, pages and stuff
		final := Construct(requests, paginated, suggestions, pages, page, found)

		// 6. Return a json for further javascript use
		writeJson(final, w)

	case "/": // This case is for first application load

		// 1. Get a prepared template with needed links to js and css
		tmpl, e := template.ParseFiles("html/blank.html")
		if e != nil {
			// if something wrong with template - throw error
			err(w, 500)
		} else {
			// 2. Just execute template without sending there anything
			tmpl.Execute(w, "")
		}

	case "/map": // This case is for map.js requests
		r.ParseForm()

		// 1. Get requested id of band
		id, e := strconv.Atoi(r.Form["id"][0])
		if e != nil {
			w.WriteHeader(400)
			w.Write([]byte("Oops! Bad request! Method requre an 'id' request as well. =)"))
		}

		// 2. Throw json wiith concerts of band by requested id
		art := Artist()[id].Concerts
		writeJson(art, w)

	default: // Case for other URLs that only can get bands id from 1 to 52

		// 1. Get number value of requested URL
		n, e := strconv.Atoi(r.URL.Path[1:])

		// If requested value is not number or greater than 52 or less than 1 - throw error
		if e != nil || n <= 0 || n > len(Artists) {
			out, _ := ioutil.ReadFile("html/404.html")
			w.WriteHeader(404)
			w.Write(out)
		} else {

			// 2. Get previously prepared template
			tmpl, er := template.ParseFiles("html/artist.html")
			if er != nil {
				err(w, 500)
			} else {
				art := Artist()

				// 3. Execute template
				tmpl.Execute(w, art[n-1])
			}
		}
	}
}

// Function for writing json
func writeJson(d interface{}, w http.ResponseWriter) {
	js, _ := json.Marshal(d)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Function for throwing error
func err(w http.ResponseWriter, n int) {
	e := map[int]string{
		400: "400: Bad Request",
		404: "404: Not Found",
		500: "500: Internal Server Error",
	}
	w.WriteHeader(n)
	w.Write([]byte(e[n]))
}
