package meili

import (
	"github.com/meilisearch/meilisearch-go"
	"log/slog"
	"main/common/config"
	"main/common/event/rule"
	"sync"
)

const RuleName = "meili"

type Config struct {
	Enable bool   // 是否启用MeiliSearch规则服务
	Addr   string // MeiliSearch服务地址
	APIKey string // MeiliSearch API密钥
}

type ClientService struct {
	meilisearch.ServiceManager
	config        *Config
	Rule          *rule.Server
	channelClient *rule.ChannelClient
	IndexMap      sync.Map // map[string]struct{}
}

func NewMeiliConfig(cfg *config.Config) *Config {
	meiliCfg := cfg.SubscribeServerConfig.Meili
	// 创建MeiliSearch服务配置
	return &Config{
		Enable: meiliCfg.Enable,
		Addr:   meiliCfg.Addr,
		APIKey: meiliCfg.APIKey,
	}
}

// NewMeiliService 创建一个新的MeiliSearch客户端服务
func NewMeiliService(cfg *Config) *ClientService {
	if !cfg.Enable {
		slog.Info("MeiliSearch规则服务未启用", slog.String("addr", cfg.Addr))
		return nil
	}
	// 初始化MeiliSearch客户端
	client := meilisearch.New(cfg.Addr, meilisearch.WithAPIKey(cfg.APIKey))
	// 创建规则服务器实例
	ruleServer := rule.NewServer()
	service := &ClientService{
		ServiceManager: client,
		config:         cfg,
		Rule:           ruleServer,
		channelClient:  ruleServer.PutNewClient(),
		IndexMap:       sync.Map{},
	}
	// 开启MeiliSearch预览功能
	go func() {
		err := service.SetExperimentalFeatures()
		if err != nil {
			slog.Error("MeiliSearchConfig 开启预览功能失败设置失败", slog.String("error", err.Error()))
			return
		}
	}()
	// 启动异步数据变更监听
	go service.asyncDataChange()
	slog.Info("MeiliSearchRule启动,监听地址", slog.String("addr", cfg.Addr))
	return service
}
