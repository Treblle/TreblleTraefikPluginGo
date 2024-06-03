package treblle_traefik

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"time"
)

type Config struct {
	ApiKey                 string
	ProjectId              string
	AdditionalFieldsToMask []string
}

func CreateConfig() *Config {
	return &Config{}
}

type Treblle struct {
	next         http.Handler
	name         string
	ApiKey       string
	ProjectId    string
	FieldsMap    map[string]bool
	serverInfo   ServerInfo
	languageInfo LanguageInfo
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	t := &Treblle{
		next: next,
		name: name,
	}

	if config.ApiKey != "" {
		t.ApiKey = config.ApiKey
	}
	if config.ProjectId != "" {
		t.ProjectId = config.ProjectId
	}
	if len(config.AdditionalFieldsToMask) > 0 {
		t.FieldsMap = generateFieldsToMask(config.AdditionalFieldsToMask)
	}

	t.serverInfo = getServerInfo()
	t.languageInfo = getLanguageInfo()

	return t, nil
}

func (t *Treblle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	reqInfo, err := t.getRequestInfo(r, startTime)
	rec := httptest.NewRecorder()

	t.next.ServeHTTP(rec, r)

	for k, v := range rec.Header() {
		w.Header()[k] = v
	}

	w.WriteHeader(rec.Code)

	_, err = w.Write(rec.Body.Bytes())
	if err != nil {
		return
	}

	if !errors.Is(err, ErrNotJson) {
		ti := Metadata{
			ApiKey:    t.ApiKey,
			ProjectID: t.ProjectId,
			Version:   "0.0.1",
			Sdk:       "go",
			Data: DataInfo{
				Server:   t.serverInfo,
				Language: t.languageInfo,
				Request:  reqInfo,
				Response: t.getResponseInfo(rec, startTime),
			},
		}
		os.Stdout.WriteString("Sending data to treblle...")
		// don't block execution while sending data to Treblle
		go t.sendToTreblle(ti)
	}
}
