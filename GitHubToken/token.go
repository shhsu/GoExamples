package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"azure.com/acr/acr-build-runner/runner/domain"
	jwt "github.com/dgrijalva/jwt-go"
)

type IntegrationAppClient struct {
	config *domain.GitHubIntegrationConfig
}

func NewIntegrationAppClient(config *domain.GitHubIntegrationConfig) *IntegrationAppClient {
	return &IntegrationAppClient{config: config}
}

func (c *IntegrationAppClient) GetAccessToken(key interface{}, installationID int) (*domain.GitHubAccessTokenResponse, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * 5).Unix(),
		"iss": c.config.AppID,
	})
	signed, err := token.SignedString(key)
	if err != nil {
		return nil, err
	}
	return c.requestAccessToken(signed, installationID)
}

func (c *IntegrationAppClient) requestAccessToken(signedToken string, installationID int) (*domain.GitHubAccessTokenResponse, error) {
	tokenEndpoint := fmt.Sprintf("https://api.github.com/installations/%d/access_tokens", installationID)
	req, err := http.NewRequest("POST", tokenEndpoint, nil)
	req.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", signedToken))
	req.Header.Add("User-Agent", "acr")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error getting token form github: %s", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Unexpected status code from github token request: %d", resp.StatusCode)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading github token response: %s", err)
	}
	var result domain.GitHubAccessTokenResponse
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, fmt.Errorf("Error parsing github response: %s", err)
	}
	return &result, nil
}
