package meili

import (
	"github.com/meilisearch/meilisearch-go"
	"main/common/config"
	"main/common/event/rule"
)

const RuleName = "MeiliSearchRule"

// NewMeiliService 创建一个新的MeiliSearch客户端服务
func NewMeiliService(cfg *config.Config) *ClientService {
	if !cfg.ExistsRuleName(RuleName) {
		return nil
	}
	client := meilisearch.New(cfg.MeiliSearch.Addr, meilisearch.WithAPIKey(cfg.MeiliSearch.APIKey))
	ruleServer := rule.NewServer()
	return &ClientService{Client: &client,
		Rule:          ruleServer,
		ChannelClient: ruleServer.PutNewClient(),
	}
}
