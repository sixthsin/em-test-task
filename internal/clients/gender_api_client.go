package clients

import (
	"em-task/cfg"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const genderApiTimeout = 10

type GenderAPIClientDeps struct {
	*cfg.Config
}

type GenderAPIClient struct {
	*cfg.Config
	*http.Client
}

func NewGenderAPIClient(deps GenderAPIClientDeps) *GenderAPIClient {
	return &GenderAPIClient{
		Config: deps.Config,
		Client: &http.Client{
			Timeout: genderApiTimeout * time.Second,
		},
	}
}

func (c *GenderAPIClient) GetGenderByName(name string) (*PersonGenderData, error) {
	resp, err := c.Client.Get("https://api.genderize.io/" + "?name=" + name)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("API returned non-200 status")
	}

	var data PersonGenderData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %v", err)
	}

	return &data, nil
}
