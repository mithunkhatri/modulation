package radio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Station struct {
	Name       string `json:"name"`
	URL        string `json:"url_resolved"`
	Homepage   string `json:"homepage"`
	Favicon    string `json:"favicon"`
	Tags       string `json:"tags"`
	Country    string `json:"country"`
	Language   string `json:"language"`
	ClickCount int    `json:"clickcount"`
}

var (
	mirrors = []string{
		"https://de1.api.radio-browser.info",
		"https://at1.api.radio-browser.info",
		"https://nl1.api.radio-browser.info",
	}
	currentMirrorIdx = 0
	mu               sync.RWMutex
	client           = &http.Client{Timeout: 5 * time.Second}
)

func GetCurrentMirror() string {
	mu.RLock()
	defer mu.RUnlock()
	return mirrors[currentMirrorIdx]
}

func RefreshMirrors() error {
	resp, err := client.Get("https://all.api.radio-browser.info/json/servers")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var servers []struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&servers); err != nil {
		return err
	}

	if len(servers) > 0 {
		newMirrors := make([]string, len(servers))
		for i, s := range servers {
			newMirrors[i] = "https://" + s.Name
		}
		mu.Lock()
		mirrors = newMirrors
		currentMirrorIdx = 0
		mu.Unlock()
	}
	return nil
}

func GetStations(query string, tag string, limit int) ([]Station, error) {
	params := url.Values{}
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("hidebroken", "true")
	params.Add("order", "votes")
	params.Add("reverse", "true")
	if query != "" {
		params.Add("name", query)
	}
	if tag != "" {
		params.Add("tag", tag)
	}

	mu.RLock()
	idx := currentMirrorIdx
	mu.RUnlock()

	// Try mirrors starting from the current one
	for i := 0; i < len(mirrors); i++ {
		mirrorIdx := (idx + i) % len(mirrors)
		baseURL := mirrors[mirrorIdx] + "/json/stations/search"

		stations, err := fetchStations(baseURL + "?" + params.Encode())
		if err == nil {
			// Update current mirror if we switched
			if mirrorIdx != idx {
				mu.Lock()
				currentMirrorIdx = mirrorIdx
				mu.Unlock()
			}
			return stations, nil
		}
	}

	return nil, fmt.Errorf("failed to fetch stations from all mirrors")
}

func fetchStations(fullURL string) ([]Station, error) {
	resp, err := client.Get(fullURL)
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
