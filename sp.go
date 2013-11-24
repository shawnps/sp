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

type ExternalId struct {
	Type string
	Id   IntString
}

type Artist struct {
	Href       string
	Name       string
	Popularity FloatString `json:"omitempty"`
}

type Album struct {
	Name         string
	Released     string       `json:"omitempty"`
	Popularity   FloatString  `json:"omitempty"`
	ExternalIds  []ExternalId `json:"external-ids"`
	Length       float64      `json:"omitempty"`
	Href         string
	Artists      []Artist `json:"artists,omitempty"`
	Availability struct {
		Territories string
	}
}

type Info struct {
	NumResults int `json:"num_results"`
	Limit      int
	Offset     int
	Query      string
	Type       string
	Page       int
}

type SearchAlbumsResponse struct {
	Info   Info
	Albums []Album
}

type SearchArtistsResponse struct {
	Info    Info
	Artists []Artist
}

type Track struct {
	Album       Album
	Name        string
	ExternalIds []ExternalId `json:"external-ids"`
	Popularity  FloatString
	Explicit    bool
	Length      float64
	Href        string
	Artists     []Artist
	TrackNumber IntString `json:"track-number"`
}

type SearchTracksResponse struct {
	Info   Info
	Tracks []Track
}

func (r *Spotify) getRequest(params map[string]string, endpoint string) ([]byte, error) {
	v := url.Values{}
	for key, val := range params {
		v.Set(key, val)
	}
	u := apiURL + endpoint + "?" + v.Encode()
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
		return SearchAlbumsResponse{}, err
	}
	var s SearchAlbumsResponse
	err = json.Unmarshal(resp, &s)
	if err != nil {
		return SearchAlbumsResponse{}, err
	}
	return s, nil
}

func (r *Spotify) SearchArtists(q string) (SearchArtistsResponse, error) {
	p := map[string]string{"q": q}
	e := "/search/1/artist.json"
	resp, err := r.getRequest(p, e)
	if err != nil {
		return SearchArtistsResponse{}, err
	}
	var s SearchArtistsResponse
	err = json.Unmarshal(resp, &s)
	if err != nil {
		return SearchArtistsResponse{}, err
	}
	return s, nil
}

func (r *Spotify) SearchTracks(q string) (SearchTracksResponse, error) {
	p := map[string]string{"q": q}
	e := "/search/1/track.json"
	resp, err := r.getRequest(p, e)
	if err != nil {
		return SearchTracksResponse{}, err
	}
	var s SearchTracksResponse
	err = json.Unmarshal(resp, &s)
	if err != nil {
		return SearchTracksResponse{}, err
	}
	return s, nil
}
