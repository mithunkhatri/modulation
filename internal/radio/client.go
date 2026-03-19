package radio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Station struct {
	Name        string `json:"name"`
	URL         string `json:"url_resolved"`
	Homepage    string `json:"homepage"`
	Favicon     string `json:"favicon"`
	Tags        string `json:"tags"`
	Country     string `json:"country"`
	Language    string `json:"language"`
	ClickCount  int    `json:"clickcount"`
}

const (
	baseURL = "https://de1.api.radio-browser.info/json/stations/search"
)

func GetStations(query string, tag string, limit int) ([]Station, error) {
	params := url.Values{}
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("hidebroken", "true")
	params.Add("order", "votes") // Votes often reflect popularity better than clickcount
	params.Add("reverse", "true")
	if query != "" {
		params.Add("name", query)
	}
	if tag != "" {
		params.Add("tag", tag)
	}

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var stations []Station
	if err := json.NewDecoder(resp.Body).Decode(&stations); err != nil {
		return nil, err
	}

	return stations, nil
}
