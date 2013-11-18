package sp

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const (
	apiURL = "http://ws.spotify.com"
)

type Spotify struct{}

// The API returns strings for some things that one would expect
// to be numbers (for example, popularity))
// FloatString takes a string and converts it to a float64.
type FloatString string

func (f *FloatString) UnmarshalJSON(i interface{}) float64 {
	s := i.(string)
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatal(err)
	}
	return n
}

// Same as FloatString, but for integers.
type IntString string

func (f *IntString) UnmarshalJSON(i interface{}) int64 {
	s := i.(string)
	n, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		log.Fatal(err)
	}
	return n
}

type Album struct {
	Name        string
	Popularity  FloatString
	ExternalIds []struct {
		Type string
		Id   IntString
	} `json:"external-ids"`
	Href    string
	Artists []struct {
		Href string
		Name string
	}
	Availability struct {
		Territories string
	}
}

type SearchAlbumsResponse struct {
	Info struct {
		NumResults int
		Limit      int
		Offset     int
		Query      string
		Type       string
		Page       int
	}
	Albums []Album
}

func (r *Spotify) getRequest(params map[string]string, endpoint string) ([]byte, error) {
	v := url.Values{}
	for key, val := range params {
		v.Set(key, val)
	}
	u := apiURL + endpoint + "?" + v.Encode()
	println(u)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (r *Spotify) SearchAlbums(q string) (SearchAlbumsResponse, error) {
	p := map[string]string{"q": q}
	e := "/search/1/album.json"
	resp, err := r.getRequest(p, e)
	if err != nil {
		return SearchAlbumsResponse{}, nil
	}
	var s SearchAlbumsResponse
	err = json.Unmarshal(resp, &s)
	if err != nil {
		return SearchAlbumsResponse{}, nil
	}
	return s, nil
}
