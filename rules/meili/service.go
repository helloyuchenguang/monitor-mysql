package meili

import (
	"github.com/meilisearch/meilisearch-go"
	"github.com/samber/lo"
	"log/slog"
	"main/common/config"
	"main/common/event/rule"
	"regexp"
	"sync"
)

const RuleName = "meili"

type Config struct {
	enable       bool           // 是否启用MeiliSearch规则服务
	addr         string         // MeiliSearch服务地址
	apiKey       string         // MeiliSearch API密钥
	indexConfigs []*IndexConfig // 索引配置列表
}

type IndexConfig struct {
	tableRegex *regexp.Regexp // 表名正则表达式
	index      string         // MeiliSearch索引名称
	searchers  []string       // 搜索字段
	filters    []string       // 过滤字段
	sorts      []string       // 排序字段
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
	// 构建索引配置
	indexConfigs := lo.Map(cfg.WatchHandlers, func(item config.WatchHandler, _ int) *IndexConfig {
		// 创建索引配置
		indexConfig := &IndexConfig{
			tableRegex: item.TableRegexp,
			index:      item.MeiliSearchIndex.Index,
			searchers:  item.MeiliSearchIndex.Searchers,
			filters:    item.MeiliSearchIndex.Filters,
			sorts:      item.MeiliSearchIndex.Sorts,
		}
		return indexConfig
	})

	// 创建MeiliSearch服务配置
	return &Config{
		enable:       meiliCfg.Enable,
		addr:         meiliCfg.Addr,
		apiKey:       meiliCfg.APIKey,
		indexConfigs: indexConfigs,
	}
}

// NewMeiliService 创建一个新的MeiliSearch客户端服务
func NewMeiliService(cfg *Config) *ClientService {
	if !cfg.enable {
		slog.Info("MeiliSearch规则服务未启用", slog.String("addr", cfg.addr))
		return nil
	}
	// 初始化MeiliSearch客户端
	client := meilisearch.New(cfg.addr, meilisearch.WithAPIKey(cfg.apiKey))
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
	slog.Info("MeiliSearchRule启动,监听地址", slog.String("addr", cfg.addr))
	return service
}
