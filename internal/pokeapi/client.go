package pokeapi

import (
	"net/http"

	"github.com/evagant-ops/pokedexcli/internal/pokecache"
)

const BaseUrl string = "https://pokeapi.co/api/v2"

type Client struct {
	cache      *pokecache.Cache
	httpClient http.Client
}

func NewClient(cache *pokecache.Cache) *Client {
	return &Client{httpClient: http.Client{}, cache: cache}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}
