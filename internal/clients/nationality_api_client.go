package clients

import (
	"em-task/cfg"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const nationalityApiTimeout = 10

type NationalityAPIClientDeps struct {
	*cfg.Config
}

type NationalityAPIClient struct {
	*cfg.Config
	*http.Client
}

func NewNationalityAPIClient(deps NationalityAPIClientDeps) *NationalityAPIClient {
	return &NationalityAPIClient{
		Config: deps.Config,
		Client: &http.Client{
			Timeout: nationalityApiTimeout * time.Second,
		},
	}
}

func (c *NationalityAPIClient) GetNationalityByName(name string) (*PersonNationalityData, error) {
	resp, err := c.Client.Get("https://api.nationalize.io/" + "?name=" + name)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("API returned non-200 status")
	}

	var data PersonNationalityData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %v", err)
	}

	return &data, nil
}
