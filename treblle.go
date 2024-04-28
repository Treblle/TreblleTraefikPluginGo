package treblle_traefik

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type Config struct {
	ApiKey    string
	ProjectId string
}

func CreateConfig() *Config {
	return &Config{}
}

type Treblle struct {
	next      http.Handler
	ApiKey    string
	ProjectId string
	name      string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.ApiKey) == 0 || len(config.ProjectId) == 0 {
		return nil, nil
	}

	return &Treblle{
		next:      next,
		ApiKey:    config.ApiKey,
		ProjectId: config.ProjectId,
		name:      name,
	}, nil
}

func (t *Treblle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	traceId := uuid.New()
	w.Header().Add("X-Traefik-Id", traceId.String())

	t.next.ServeHTTP(w, r)
}
