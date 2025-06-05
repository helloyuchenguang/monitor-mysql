package config

import (
	"fmt"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"regexp"
)

// Config 配置文件结构体
type Config struct {
	Database              DatabaseConfig        `yaml:"database"`
	WatchHandlers         []WatchHandler        `yaml:"watchHandlers"`
	SubscribeServerConfig SubscribeServerConfig `yaml:"subscribeServerConfig"` // 订阅服务器配置
}

// SubscribeServerConfig 订阅服务器配置
type SubscribeServerConfig struct {
	Grpc  GRPCConfig        `yaml:"grpc"`  // gRPC规则
	SSE   SSEConfig         `yaml:"sse"`   // SSE规则
	Meili MeiliSearchConfig `yaml:"meili"` // MeiliSearch规则
}

type WatchHandler struct {
	TableRegexp      *regexp.Regexp
	Table            string           `yaml:"table"`
	MeiliSearchIndex MeiliSearchIndex `yaml:"meiliSearchIndex"`
	Rules            []string         `yaml:"rules"`
}

// BuildRegexp 构建正则表达式
func (w *WatchHandler) buildRegexp() {
	if w.Table == "" {
		panic("tableRegex不能为空")
	}
	reg, err := regexp.Compile(w.Table)
	if err != nil {
		panic("无法编译正则表达式: " + w.Table)
	}
	w.TableRegexp = reg
}

type DatabaseConfig struct {
	Addr              string   `yaml:"addr"`
	User              string   `yaml:"user"`
	Password          string   `yaml:"password"`
	Flavor            string   `yaml:"flavor"`
	ServerID          uint32   `yaml:"serverId"`
	DumpExecutionPath string   `yaml:"dumpExecutionPath"`
	IncludeTableRegex []string `yaml:"includeTableRegex"`
}

type SSEConfig struct {
	Enable bool   `yaml:"enable"`
	Addr   string `yaml:"addr"` // Web服务地址
}

type GRPCConfig struct {
	Enable bool   `yaml:"enable"`
	Addr   string `yaml:"addr"` // gRPC服务地址
}

type MeiliSearchConfig struct {
	Enable bool   `yaml:"enable"`
	Addr   string `yaml:"addr"`   // MeiliSearch服务地址
	APIKey string `yaml:"apiKey"` // MeiliSearchConfig API Key
}

type MeiliSearchIndex struct {
	Index     string   `yaml:"index"`     // MeiliSearch索引名称
	Searchers []string `yaml:"searchers"` // 搜索字段
	Filters   []string `yaml:"filters"`   // 过滤字段
	Sorts     []string `yaml:"sorts"`     // 排序字段
}

// LoadConfig 加载配置文件
func LoadConfig(cfgFile string) *Config {
	f, err := os.Open(cfgFile)
	if err != nil {
		slog.Error(fmt.Sprintf("打开配置文件失败: %v", err))
		panic("无法打开配置文件: " + cfgFile)
	}
	// 确保在函数结束时关闭文件
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			slog.Error(fmt.Sprintf("关闭配置文件失败: %v", err))
		}
	}(f)

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		panic("无法解析配置文件: " + cfgFile)
	}
	for _, handler := range cfg.WatchHandlers {
		// 编译正则表达式
		handler.buildRegexp()
		// rules 去重
		handler.Rules = lo.Uniq(handler.Rules)
	}
	slog.Info(fmt.Sprintf("读取配置文件: %-v", cfg))
	return &cfg
}
