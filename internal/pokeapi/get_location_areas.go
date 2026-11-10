package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetLocationAreas(url string) (LocationAreas, error) {
	if url == "" {
		url = fmt.Sprintf("%v/location-area/?limit=%d&offset=%d", BaseUrl, 20, 0)
	}

	var locationAreaResp LocationAreas
	e, ok := c.cache.Get(url)
	if ok {
		err := json.Unmarshal(e, &locationAreaResp)
		if err != nil {
			return locationAreaResp, err
		}
		return locationAreaResp, nil

	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return locationAreaResp, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return locationAreaResp, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return locationAreaResp, err
	}

	err = json.Unmarshal(body, &locationAreaResp)
	if err != nil {
		return locationAreaResp, err
	}

	return locationAreaResp, nil
}
