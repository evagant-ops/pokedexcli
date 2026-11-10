package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetLocationArea(url string, area string) (LocationArea, error) {
	if url == "" {
		url = fmt.Sprintf("%v/location-area/%v", BaseUrl, area)
	}

	var locationAreaResp LocationArea
	e, ok := c.cache.Get(url)
	if ok {
		err := json.Unmarshal(e, &locationAreaResp)
		if err != nil {
			return locationAreaResp, err
		}

		return locationAreaResp, nil
	}

	var locationAreasResp LocationArea
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return locationAreasResp, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return locationAreasResp, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return locationAreasResp, err
	}

	err = json.Unmarshal(body, &locationAreasResp)
	if err != nil {
		return locationAreasResp, err
	}

	return locationAreasResp, nil
}
