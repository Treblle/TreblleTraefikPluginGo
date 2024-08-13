package TreblleTraefikPluginGo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"
)

type Config struct {
	ApiKey                 string
	ProjectId              string
	AdditionalFieldsToMask []string
	RoutesToBlock          []string
	RoutesRegex            string
	DebugMode              bool
}

func CreateConfig() *Config {
	return &Config{}
}

type Treblle struct {
	next          http.Handler
	name          string
	ApiKey        string
	ProjectId     string
	FieldsMap     map[string]bool
	RoutesToBlock []string
	RoutesRegex   *regexp.Regexp
	serverInfo    ServerInfo
	languageInfo  LanguageInfo
	DebugMode     bool
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
	if len(config.RoutesToBlock) > 0 {
		t.RoutesToBlock = config.RoutesToBlock
	}
	if config.DebugMode {
		t.DebugMode = true
	}
	if config.RoutesRegex != "" {
		re, err := regexp.Compile(config.RoutesRegex)
		if err != nil {
			return nil, err
		}

		t.RoutesRegex = re
	}

	t.serverInfo = getServerInfo()
	t.languageInfo = getLanguageInfo()

	return t, nil
}

func (t *Treblle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(t.RoutesToBlock) > 0 {
		for _, route := range t.RoutesToBlock {
			if r.RequestURI == route {
				t.next.ServeHTTP(w, r)
				return
			}
		}
	}

	if t.RoutesRegex != nil {
		if t.RoutesRegex.MatchString(r.RequestURI) {
			t.next.ServeHTTP(w, r)
			return
		}
	}

	startTime := time.Now()
	reqInfo, _ := t.getRequestInfo(r, startTime)
	rec := httptest.NewRecorder()

	t.next.ServeHTTP(rec, r)

	for k, v := range rec.Header() {
		w.Header()[k] = v
	}

	w.WriteHeader(rec.Code)

	_, _ = w.Write(rec.Body.Bytes())

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

	// don't block execution while sending data to Treblle
	go t.sendToTreblle(ti)
}
