package meili

import (
	"github.com/meilisearch/meilisearch-go"
	cmap "github.com/orcaman/concurrent-map/v2"
	"log/slog"
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
	service := &ClientService{Client: client,
		Rule:          ruleServer,
		ChannelClient: ruleServer.PutNewClient(),
		IndexMap:      cmap.New[bool](),
	}
	go service.asyncDataChange()
	slog.Info("MeiliSearchRule启动,监听地址", slog.String("addr", cfg.MeiliSearch.Addr))
	return service
}
