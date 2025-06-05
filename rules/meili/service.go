package meili

import (
	"github.com/meilisearch/meilisearch-go"
	"log/slog"
	"main/common/config"
	"main/common/event/rule"
	"sync"
)

const RuleName = "MeiliSearchRule"

// NewMeiliService 创建一个新的MeiliSearch客户端服务
func NewMeiliService(cfg *config.Config) *ClientService {
	// 检查配置中是否存在MeiliSearch规则
	if !cfg.ExistsRuleName(RuleName) {
		return nil
	}

	// 初始化MeiliSearch客户端
	client := meilisearch.New(cfg.MeiliSearch.Addr, meilisearch.WithAPIKey(cfg.MeiliSearch.APIKey))
	// 创建规则服务器实例
	ruleServer := rule.NewServer()
	service := &ClientService{
		ServiceManager: client,
		Rule:           ruleServer,
		channelClient:  ruleServer.PutNewClient(),
		IndexMap:       sync.Map{},
	}
	// 开启MeiliSearch预览功能
	go func() {
		err := service.SetExperimentalFeatures()
		if err != nil {
			slog.Error("MeiliSearch 开启预览功能失败设置失败", slog.String("error", err.Error()))
			return
		}
	}()
	// 启动异步数据变更监听
	go service.asyncDataChange()
	slog.Info("MeiliSearchRule启动,监听地址", slog.String("addr", cfg.MeiliSearch.Addr))
	return service
}
