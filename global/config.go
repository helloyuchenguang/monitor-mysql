package global

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

// Config 配置文件结构体
type Config struct {
	Database struct {
		Addr              string   `yaml:"addr"`
		User              string   `yaml:"user"`
		Password          string   `yaml:"password"`
		Flavor            string   `yaml:"flavor"`
		ServerID          uint32   `yaml:"server_id"`
		DumpExecutionPath string   `yaml:"dump_execution_path"`
		IncludeTableRegex []string `yaml:"include_table_regex"`
	} `yaml:"database"`
	WatchHandlers []struct {
		TableRegex string   `yaml:"table_regex"`
		Rules      []string `yaml:"rules"`
	} `yaml:"watch_handlers"`
	Web struct {
		Addr string `yaml:"addr"`
	}
	GRPC struct {
		Addr string `yaml:"addr"`
	}
}

// LoadConfig 加载配置文件
func LoadConfig(cfgFile string) (Config, error) {
	f, err := os.Open(cfgFile)
	if err != nil {
		slog.Error(fmt.Sprintf("打开配置文件失败: %v", err))
		return Config{}, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			slog.Error(fmt.Sprintf("关闭配置文件失败: %v", err))
		}
	}(f)

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		slog.Error(fmt.Sprintf("解析配置文件失败: %v", err))
		return Config{}, err
	}
	slog.Info(fmt.Sprintf("读取配置文件: %-v", cfg))
	return cfg, nil
}
