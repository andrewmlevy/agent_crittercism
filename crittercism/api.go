package crittercism

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/jsonq"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var token string
var tokenExpires time.Time

func Request(method string, path string, params string) (*jsonq.JsonQuery, error) {

	url := fmt.Sprintf("https://developers.crittercism.com:443/v1.0/%s", path)

	p := []byte(params)

	client := &http.Client{}
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(p))
	req.Header.Set("Authorization", "Bearer cg299DBJDbZfLJq621uFewRlzvwDjKZM")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(string(body)))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)

	return jq, err
}
