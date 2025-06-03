package meili

import (
	"github.com/meilisearch/meilisearch-go"
	"main/common/config"
)

const RuleName = "MeiliSearchRule"

type ClientService struct {
	Client *meilisearch.ServiceManager
}

// NewMeiliService 创建一个新的MeiliSearch客户端服务
func NewMeiliService(cfg *config.Config) *ClientService {
	if !cfg.ExistsRuleName(RuleName) {
		return nil
	}
	client := meilisearch.New(cfg.MeiliSearch.Addr, meilisearch.WithAPIKey(cfg.MeiliSearch.APIKey))
	return &ClientService{Client: &client}
}
