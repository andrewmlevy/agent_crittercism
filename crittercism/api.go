package crittercism

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/jsonq"
	"io/ioutil"
	"net/http"
	"strings"
)

const crittercismAPIURL = "https://developers.crittercism.com:443/v1.0"

var accessToken string = "cg299DBJDbZfLJq621uFewRlzvwDjKZM"
var refreshToken string = ""
var accessTokenExpires int

type CrittercismAPIClient struct {
	accessToken    string
	refreshToken   string
	expirationDate int
	appId          string
}

// getCrittercismOAuthToken fetches a new OAuth Token from the Crittercism API given a username and password
func getCrittercismOAuthToken(login, password string) (token string, expires int, err error) {
	var params = fmt.Sprintf(`{"login": "%s", "password": "%s"}}`, login, password)

	// Construct REST Request
	url := fmt.Sprintf("%s/token", crittercismAPIURL)
	p := []byte(params)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(p))
	req.Header.Set("Content-Type", "application/json")

	// Make Request
	if resp, err := client.Do(req); err == nil {
		defer resp.Body.Close()

		// Parse JSON
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			data := map[string]interface{}{}
			dec := json.NewDecoder(strings.NewReader(string(body)))
			dec.Decode(&data)
			jq := jsonq.NewQuery(data)

			// Parse out the token
			token, _ := jq.String("access_token")
			expires, _ := jq.Int("expires")
			return token, expires, nil

		} else {
			return "", 0, err // Parse Error
		}

	} else {
		return "", 0, err // Request Error
	}
}

func NewCrittercismAPIClient(login, password, appId string) (*CrittercismAPIClient, error) {
	var err error

	if accessToken == "" {
		accessToken, accessTokenExpires, err = getCrittercismOAuthToken(login, password)

		if err != nil {
			return nil, err
		}
	}

	return &CrittercismAPIClient{
		accessToken:    accessToken,
		expirationDate: accessTokenExpires,
		appId:          appId,
	}, nil
}

// Request will make a request of the Crittercism API
// This will return a github.com/jmoiron/jsonq JSON query object
func (c *CrittercismAPIClient) Request(method, path, params string) (jq *jsonq.JsonQuery, err error) {
	// Construct REST Request
	url := fmt.Sprintf("%s/%s", crittercismAPIURL, path)
	p := []byte(params)
	client := &http.Client{}
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(p))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	// Make Request
	if resp, err := client.Do(req); err == nil {
		defer resp.Body.Close()

		// Parse Body
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			data := map[string]interface{}{}
			dec := json.NewDecoder(strings.NewReader(string(body)))
			dec.Decode(&data)
			jq = jsonq.NewQuery(data)
			return jq, nil

		} else {
			return nil, err // Parse Error
		}

	} else {
		return nil, err // Request Error
	}
}

func (c *CrittercismAPIClient) FetchGraph(path, name string, duration int) ([]float64, error) {
	params := fmt.Sprintf(`{"params":{"graph": "%s", "duration": 86400, "appId": "%s"}}`, name, c.appId)

	// Get the data from Crittercism

	jq, err := c.Request("POST", path, params)

	if err != nil {
		return []float64{}, err
	}

	return jq.ArrayOfFloats("data", "series", "0", "points")
}
