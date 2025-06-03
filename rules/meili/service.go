package meili

import (
	"github.com/meilisearch/meilisearch-go"
	"main/common/config"
)

const RuleName = "meili"

type ClientService struct {
	Client *meilisearch.ServiceManager
}

func NewMeiliService(cfg *config.Config) *ClientService {
	if !cfg.ExistsRuleName(RuleName) {
		return nil
	}
	client := meilisearch.New(cfg.MeiliSearch.Addr, meilisearch.WithAPIKey(cfg.MeiliSearch.APIKey))
	return &ClientService{Client: &client}
}
