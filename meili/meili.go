package meili

import (
	"github.com/meilisearch/meilisearch-go"
)

type ClientService struct {
	Client *meilisearch.ServiceManager
}

type ClientConfig struct {
	Addr   string
	APIKey string
}

func NewClient(cfg ClientConfig) *meilisearch.ServiceManager {
	client := meilisearch.New(cfg.Addr, meilisearch.WithAPIKey(cfg.APIKey))
	return &client
}

func NewMeiliService(client *meilisearch.ServiceManager) *ClientService {
	return &ClientService{Client: client}
}
