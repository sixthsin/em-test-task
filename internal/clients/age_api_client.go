package clients

import (
	"em-task/cfg"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const ageApiTimeout = 10

type AgeAPIClientDeps struct {
	*cfg.Config
}

type AgeAPIClient struct {
	*cfg.Config
	*http.Client
}

func NewAgeAPIClient(deps AgeAPIClientDeps) *AgeAPIClient {
	return &AgeAPIClient{
		Config: deps.Config,
		Client: &http.Client{
			Timeout: ageApiTimeout * time.Second,
		},
	}
}

func (c *AgeAPIClient) GetAgeByName(name string) (*PersonAgeData, error) {
	resp, err := c.Client.Get("https://api.agify.io/" + "?name=" + name)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("API returned non-200 status")
	}

	var data PersonAgeData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}
