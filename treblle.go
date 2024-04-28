package treblle_traefik

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	timeout = 2 * time.Second
)

func getBaseUrl() string {
	treblleBaseUrls := []string{
		"https://rocknrolla.treblle.com",
		"https://punisher.treblle.com",
		"https://sicario.treblle.com",
	}
	rand.Seed(time.Now().Unix())
	randomUrlIndex := rand.Intn(len(treblleBaseUrls))

	return treblleBaseUrls[randomUrlIndex]
}

func (t *Treblle) sendToTreblle(info Metadata) {
	baseUrl := getBaseUrl()

	jsonData, err := json.Marshal(info)
	if err != nil {
		return
	}
	os.Stdout.WriteString(string(jsonData))
	req, err := http.NewRequest(http.MethodPost, baseUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", t.ApiKey)

	client := &http.Client{Timeout: timeout}
	res, _ := client.Do(req)
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	os.Stdout.WriteString(string(body))
}
