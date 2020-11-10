package pack

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type requests struct {
	Page                  int
	Txt, Sort, ByLocation string
	ByAlbum, ByEvent      []int64
	ByCreation, ByMembers []int
}

// 0. Make a global Request variable

func WriteRequests(r *http.Request) requests {
	var R requests
	r.ParseForm()
	R.Page = getInt(r, "page")
	R.Txt = strings.ToLower(getString(r, "txt"))
	R.Sort = strings.ToLower(getString(r, "sort"))
	R.ByLocation = strings.ToLower(getString(r, "bylocation"))
	R.ByCreation = getIntN(r, "bycreation")
	R.ByAlbum = getInt64(r, "byalbum")
	R.ByEvent = getInt64(r, "byevent")
	R.ByMembers = getIntN(r, "bymembers")
	return R
}

func getInt(r *http.Request, form string) int {
	if v, ok := r.Form[form]; ok {
		n, e := strconv.Atoi(v[0])
		if e != nil {
			return 1
		}
		return n
	}
	return 1
}

func getIntN(r *http.Request, form string) []int {
	if v, ok := r.Form[form]; ok {
		arr := strings.Split(v[0], "~")
		x, ex := strconv.Atoi(arr[0])
		y, ey := strconv.Atoi(arr[1])

		if ex != nil && ey != nil {
			return []int{0, 10000}
		}
		if ex != nil {
			return []int{0, y}
		}
		if ey != nil {
			return []int{x, 10000}
		}
		return []int{x, y}
	}
	return []int{0, 10000}
}

func getInt64(r *http.Request, form string) []int64 {
	x, _ := DateToNano("1900-01-01")
	y, _ := DateToNano("2025-01-01")

	if v, ok := r.Form[form]; ok {
		arr := strings.Split(v[0], "~")
		z, ez := DateToNano(arr[0])
		m, em := DateToNano(arr[1])

		if ez != nil && em != nil {
			return []int64{x, y}
		}

		if ez != nil {
			return []int64{x, m}
		}

		if em != nil {
			return []int64{z, y}
		}

		return []int64{z, m}
	}
	return []int64{x, y}
}

func getString(r *http.Request, form string) string {
	if v, ok := r.Form[form]; ok {
		return v[0]
	}
	return ""
}

// 1. Filter Bands by filtering settings

func filterBands(R requests) (Bands, []map[string]string, []map[string]string) {
	var f Bands
	var s []map[string]string
	sugs := make(map[string]string)
	for _, c := range Artists {

		// Hard filters - skips if has no match
		if !inRange(c.Album, R.ByAlbum) {
			continue
		}

		if !inRange(c.Establish, R.ByCreation) {
			continue
		}

		if !inRange(len(c.Members), R.ByMembers) {
			continue
		}

		match := false
		for _, e := range c.Concerts {
			if !cont(e.Place, R.ByLocation) {
				continue
			} else {
				match = true
			}
		}

		if !match {
			continue
		}

		match = false
		for _, e := range c.Concerts {
			for _, d := range e.Dates {
				if !inRange(d, R.ByEvent) {
					continue
				} else {
					match = true
				}
			}
		}

		if !match {
			continue
		}

		// Soft filters - append to variable 'Filtered' if has any match

		matches := map[string]string{"name": "", "member": "", "location": "", "establish": "", "event": "", "album": ""}

		if R.Txt == "" {
			matches["establish"] = "has"
			matches["album"] = "has"
			matches["event"] = "has"
		}

		if has(c.Name, R.Txt) {
			matches["name"] = "has"
			sugs[c.Name] = c.Name + " - artist/band"
		}

		n, e := strconv.Atoi(R.Txt)
		if e == nil {
			if c.Establish == n {
				matches["establish"] = "has"
				x := strconv.Itoa(c.Establish)
				sugs[x] = x + " - establish date"
			}
		}

		z, er := DateToNano(R.Txt)
		if er == nil {
			if z == c.Album {
				matches["album"] = "has"
				x := NanoToDate(c.Album)
				sugs[x] = x + " - first album date"
			}
		}

		for _, m := range c.Members {
			if has(m, R.Txt) {
				matches["member"] = "has"
				sugs[m] = m + " - band member"
			}
		}

		for _, e := range c.Concerts {
			if has(e.Place, R.Txt) {
				matches["location"] = "has"
				sugs[e.Place] = e.Place + " - location"
			}
			for _, d := range e.Dates {
				k, e := DateToNano(R.Txt)
				if e == nil {
					if k == d {
						matches["event"] = "has"
						x := NanoToDate(d)
						sugs[x] = x + " - first album date"
					}
				}
			}
		}

		for _, h := range matches {
			if h == "has" {
				f = append(f, c)
				s = append(s, matches)
				break
			}
		}

	}

	var new []map[string]string
	for key, val := range sugs {
		new = append(new, map[string]string{"label": val, "value": key})
	}

	return f, s, new
}

func inRange(n, rng interface{}) bool {
	switch rng.(type) {
	case []int64:
		x := n.(int64)
		r := rng.([]int64)
		if x >= r[0] && x <= r[1] {
			return true
		}
		return false
	case []int:
		x := n.(int)
		r := rng.([]int)
		if x >= r[0] && x <= r[1] {
			return true
		}
		return false
	}
	return false
}

func cont(s, p string) bool {
	return strings.Contains(strings.ToLower(s), p)
}

func has(s, p string) bool {
	return strings.HasPrefix(strings.ToLower(s), p)
}

// 2. Make array just for api - no need to show all info

type bandsApi []struct {
	Id         int
	Name       string
	Image      string
	MembersNum int
	Establish  int
	Album      string
	Match      map[string]string
}

func BandList(R requests) (bandsApi, []map[string]string) {
	b, m, s := filterBands(R)
	l := make(bandsApi, len(b))
	for i, c := range b {
		l[i].Id = c.Id
		l[i].Name = c.Name
		l[i].Image = c.Image
		l[i].MembersNum = len(c.Members)
		l[i].Establish = c.Establish
		l[i].Album = NanoToDate(c.Album)
		l[i].Match = m[i]
	}
	return l, s
}

// 3. Sort Band by sorting settings
func SortBands(b bandsApi, R requests) bandsApi {
	switch R.Sort {

	// Names by A-Z
	case "az":
		sort.Slice(b, func(i, j int) bool {
			return b[i].Name < b[j].Name
		})
	// Names by Z-A
	case "za":
		sort.Slice(b, func(i, j int) bool {
			return b[i].Name > b[j].Name
		})

	// Members number by Less-Big
	case "lb":
		sort.Slice(b, func(i, j int) bool {
			return b[i].MembersNum < b[j].MembersNum
		})

	// Members number by Big-Less
	case "bl":
		sort.Slice(b, func(i, j int) bool {
			return b[i].MembersNum > b[j].MembersNum
		})

	// Bands by young-old
	case "yo":
		sort.Slice(b, func(i, j int) bool {
			return b[i].Establish > b[j].Establish
		})

	// Bands by old-young
	case "oy":
		sort.Slice(b, func(i, j int) bool {
			return b[i].Establish < b[j].Establish
		})

	case "oya":
		sort.Slice(b, func(i, j int) bool {
			return b[i].Album < b[j].Album
		})

	case "yoa":
		sort.Slice(b, func(i, j int) bool {
			return b[i].Album > b[j].Album
		})
	}
	return b
}

// 4. Paginate Bands by page

const PageSize = 8

func Paginate(b bandsApi, R requests) (bandsApi, int, int, int) {
	var paged bandsApi
	page := R.Page
	l := len(b)

	pages := l / PageSize

	if l%PageSize != 0 {
		pages++
	}

	if page <= 0 || page > pages {
		page = 1
	}

	start := PageSize * (page - 1)
	end := start + PageSize

	if l < PageSize || end > l {
		end = l
	}

	if l != 0 {
		for x := start; x < end; x++ {
			paged = append(paged, b[x])
		}
	}

	return paged, pages, page, l
}

// 5. Cunstruct final array - add page nums, search suggestions

type final struct {
	Found, Page, Pages int
	Suggestions        []map[string]string
	List               bandsApi
}

func Construct(R requests, b bandsApi, suggestions []map[string]string, pages, page, found int) final {
	var res final
	res.Found = found
	res.Page = page
	res.Pages = pages
	if R.Txt != "" {
		res.Suggestions = suggestions
	}
	res.List = b
	return res
}

type bandz []struct {
	Id         int
	Name       string
	Image      string
	MembersNum int
	Members    []string
	Establish  int
	Album      string
	Concerts   []eventz
}

type eventz struct {
	Place  string
	Coords []float64
	Dates  []string
}

func Artist() bandz {
	b := make(bandz, len(Artists))
	for i, v := range Artists {
		b[i].Id = v.Id
		b[i].Name = v.Name
		b[i].Image = v.Image
		b[i].MembersNum = len(v.Members)
		b[i].Members = v.Members
		b[i].Establish = v.Establish
		b[i].Album = NanoToDate(v.Album)
		con := make([]eventz, len(v.Concerts))
		for z, e := range v.Concerts {
			con[z].Place = e.Place
			con[z].Coords = e.Coords
			var dates []string
			for _, d := range e.Dates {
				dates = append(dates, NanoToDate(d))
			}
			con[z].Dates = dates
		}
		b[i].Concerts = con
	}
	return b
}
