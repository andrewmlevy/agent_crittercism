package crittercism

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/jsonq"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const crittercismAPIURL = "https://developers.crittercism.com:443/v1.0"

var accessToken string = ""
var refreshToken string = ""
var accessTokenExpires time.Time

// Request will make a request of the Crittercism API
// This will return a github.com/jmoiron/jsonq JSON query object
func Request(method string, path string, params string, config map[string]string) (jq *jsonq.JsonQuery, err error) {

	// Check for accessToken, getting a new one if needed
	if accessToken == "" || time.Now().After(accessTokenExpires) {
		if token, expires, err := getOAuthToken(config); err == nil {
			accessToken = token
			accessTokenExpires = time.Now().Add(time.Duration(expires) * time.Second)
		} else {
			return nil, err
		}
	}

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

// getOAuthToken fetches a new OAuth Token from the Crittercism API given a username and password
func getOAuthToken(config map[string]string) (token string, expires int, err error) {

	var params = fmt.Sprintf(`{"login": "%s", "password": "%s"}}`, config["login"], config["password"])

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