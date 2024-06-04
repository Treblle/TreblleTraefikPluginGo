package TreblleTraefikPluginGo

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

const (
	timeout = 2 * time.Second
)

var treblleBaseUrls = []string{
	"https://rocknrolla.treblle.com",
	"https://punisher.treblle.com",
	"https://sicario.treblle.com",
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getBaseUrl() string {
	randomUrlIndex := rand.Intn(len(treblleBaseUrls))

	return treblleBaseUrls[randomUrlIndex]
}

func (t *Treblle) sendToTreblle(info Metadata) {
	baseUrl := getBaseUrl()

	jsonData, err := json.Marshal(info)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, baseUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", t.ApiKey)

	client := &http.Client{Timeout: timeout}
	_, _ = client.Do(req)
}
