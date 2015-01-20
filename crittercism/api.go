package crittercism

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/jsonq"
	"github.com/telemetryapp/gotelemetry"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const crittercismAPIURL = "https://developers.crittercism.com:443/v1.0"

type CrittercismAPIClient struct {
	accessToken  string
	refreshToken string
	appId        string
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

func NewCrittercismAPIClient(accessToken, appId string) (*CrittercismAPIClient, error) {
	// var err error

	// if accessToken == "" {
	// 	accessToken, accessTokenExpires, err = getCrittercismOAuthToken(login, password)

	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return &CrittercismAPIClient{
		accessToken: accessToken,
		appId:       appId,
	}, nil
}

type CrittercismAPIParams struct {
	GroupBy  string `json:"groupBy,omitempty"`
	Graph    string `json:"graph"`
	Duration int    `json:"duration"`
	AppID    string `json:"appId"`
}

// Request will make a request of the Crittercism API
// This will return a github.com/jmoiron/jsonq JSON query object

type CrittercismAPIClientResponseHandler func(res *http.Response, err error)

func (c *CrittercismAPIClient) RawRequest(method, path string, params *CrittercismAPIParams, handler CrittercismAPIClientResponseHandler) {
	// Construct REST Request
	url := fmt.Sprintf("%s/%s", crittercismAPIURL, path)

	log.Printf("REQUESTING %s, %s", method, url)

	var buffer *bytes.Buffer

	if params != nil {
		p, err := json.Marshal(map[string]interface{}{"params": params})

		if err != nil {
			handler(nil, err)
			return
		}

		buffer = bytes.NewBuffer(p)
	} else {
		buffer = bytes.NewBuffer([]byte{})
	}

	client := &http.Client{}
	req, _ := http.NewRequest(method, url, buffer)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Set("Content-Type", "application/json")

	// Make Request
	resp, err := client.Do(req)

	handler(resp, err)
}

func (c *CrittercismAPIClient) Request(method, path string, params *CrittercismAPIParams) (jq *jsonq.JsonQuery, outErr error) {
	c.RawRequest(
		method,
		path,
		params,
		func(resp *http.Response, err error) {
			if resp != nil && resp.Body != nil {
				defer resp.Body.Close()
			}

			if err != nil {
				outErr = err
				return
			}

			// Parse Body
			if body, outErr := ioutil.ReadAll(resp.Body); outErr == nil {
				data := map[string]interface{}{}
				dec := json.NewDecoder(strings.NewReader(string(body)))
				dec.Decode(&data)
				jq = jsonq.NewQuery(data)
			}
		},
	)

	return jq, outErr
}

func (c *CrittercismAPIClient) NewCrittercismAPIParams(groupBy, graph string, duration int) CrittercismAPIParams {
	return CrittercismAPIParams{
		GroupBy:  groupBy,
		Graph:    graph,
		Duration: duration,
		AppID:    c.appId,
	}
}

type CrittercismCrashStatusDetail struct {
	Reason             string  `json:"reason"`
	DisplayReason      *string `json:"displayReason"`
	Name               *string `json:"name"`
	UniqueSessionCount int     `json:"uniqueSessionCount"`
	SessionCount       int     `json:"sessionCount"`
}

type CrittercismCrashSlice []CrittercismCrashStatusDetail

func (c CrittercismCrashSlice) Aggregate() CrittercismCrashSlice {
	result := CrittercismCrashSlice{}

	for _, crash := range c {
		found := false

		for _, newCrash := range result {
			if newCrash.Reason == crash.Reason {
				newCrash.SessionCount += crash.SessionCount
				found = true
				break
			}
		}

		if !found {
			result = append(
				result,
				CrittercismCrashStatusDetail{
					Reason:       crash.Reason,
					SessionCount: crash.SessionCount,
				},
			)
		}
	}

	return result
}

func (c CrittercismCrashSlice) Len() int           { return len(c) }
func (c CrittercismCrashSlice) Less(i, j int) bool { return c[i].Reason < c[j].Reason }
func (c CrittercismCrashSlice) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func (c *CrittercismAPIClient) FetchCrashStatus() (result CrittercismCrashSlice, outErr error) {
	result = CrittercismCrashSlice{}

	c.RawRequest(
		"GET",
		"app/"+c.appId+"/crash/summaries?status=unresolved&sortBy=sessionCount&sortOrder=DESC",
		nil,
		func(resp *http.Response, err error) {
			if resp != nil && resp.Body != nil {
				defer resp.Body.Close()
			}

			if err != nil {
				outErr = err
				return
			}

			data, outErr := ioutil.ReadAll(resp.Body)

			if outErr != nil {
				return
			}

			outErr = json.Unmarshal(data, &result)
		},
	)

	return result, outErr
}

func (c *CrittercismAPIClient) FetchGraphRaw(path, name, groupBy string, duration int) (*jsonq.JsonQuery, error) {
	params := c.NewCrittercismAPIParams(groupBy, name, duration)

	return c.Request("POST", path, &params)
}

func (c *CrittercismAPIClient) FetchGraph(path, name string, duration int) ([]float64, error) {
	jq, err := c.FetchGraphRaw(path, name, "", duration)

	if err != nil {
		return []float64{}, err
	}

	return jq.ArrayOfFloats("data", "series", "0", "points")
}

func (c *CrittercismAPIClient) FetchGraphIntoFlow(path, name string, duration int, scale int, f *gotelemetry.Flow) error {
	if data, found := f.GraphData(); found == true {
		series, err := c.FetchGraph(path, name, duration)

		if err != nil {
			return err
		}

		// Eliminate last value
		if len(series) > 1 {
			series = series[:len(series)-1]
		}

		data.Series[0].Values = series
		data.StartTime = time.Now().Add(-time.Duration(scale) * time.Second).Unix()

		return nil
	}

	return gotelemetry.NewError(400, "Cannot extract value data from flow"+f.Tag)
}

func (c *CrittercismAPIClient) FetchLastValueOfGraph(path, name string, interval int) (float64, error) {
	series, err := c.FetchGraph(path, name, interval)

	if err != nil || len(series) == 0 {
		return -1, err
	}

	return series[len(series)-1], err
}

func (c *CrittercismAPIClient) FetchSumOfGraph(path, name string, interval int) (float64, error) {
	series, err := c.FetchGraph(path, name, interval)

	if err != nil || len(series) == 0 {
		return -1, err
	}

	sum := 0.0

	for _, value := range series {
		sum += value
	}

	return sum, err
}

func (c *CrittercismAPIClient) FetchLastValueOfGraphIntoFlow(path, name string, interval int, f *gotelemetry.Flow) error {
	if data, found := f.ValueData(); found == true {
		value, err := c.FetchLastValueOfGraph(path, name, interval)

		if err != nil {
			return err
		}

		data.Value = value

		return nil
	}

	return gotelemetry.NewError(400, "Cannot extract value data from flow"+f.Tag)
}

func (c *CrittercismAPIClient) FetchSumOfGraphIntoFlow(path, name string, interval int, f *gotelemetry.Flow) error {
	if data, found := f.ValueData(); found == true {
		value, err := c.FetchSumOfGraph(path, name, interval)

		if err != nil {
			return err
		}

		data.Value = value

		return nil
	}

	return gotelemetry.NewError(400, "Cannot extract value data from flow"+f.Tag)
}
