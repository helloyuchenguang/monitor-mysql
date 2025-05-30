package monitor

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"log/slog"
	"main/common/config"
	"main/common/event/edit"
	"main/mgrpc"
	"main/web"
	"regexp"
)

type CanalMonitorService struct {
	cfg      *Config
	SSERule  *web.SSERuleService
	grpcRule *mgrpc.GRPCRuleService
}

type WatchHandlers []struct {
	TableRegex string
	Rules      []string
}

type Database struct {
	Addr              string
	User              string
	Password          string
	Flavor            string
	ServerID          uint32
	DumpExecutionPath string
	IncludeTableRegex []string
}

type Config struct {
	Database      Database
	WatchHandlers WatchHandlers
}

func NewMonitorConfig(cfg *config.Config) *Config {
	cfgWatchHandlers := cfg.WatchHandlers
	// 创建配置
	watchHandlers := make(WatchHandlers, len(cfgWatchHandlers))
	for i, handler := range cfgWatchHandlers {
		watchHandlers[i] = struct {
			TableRegex string
			Rules      []string
		}{
			TableRegex: handler.TableRegex,
			Rules:      handler.Rules,
		}
	}
	return &Config{
		Database: Database{
			Addr:              cfg.Database.Addr,
			User:              cfg.Database.User,
			Password:          cfg.Database.Password,
			Flavor:            cfg.Database.Flavor,
			ServerID:          cfg.Database.ServerID,
			DumpExecutionPath: cfg.Database.DumpExecutionPath,
			IncludeTableRegex: cfg.Database.IncludeTableRegex,
		},
		WatchHandlers: watchHandlers,
	}
}

func NewMonitorService(cfg config.Config, sseRule *web.SSERuleService, grpcRule *mgrpc.GRPCRuleService) *CanalMonitorService {
	return &CanalMonitorService{
		cfg:      NewMonitorConfig(&cfg),
		SSERule:  sseRule,
		grpcRule: grpcRule,
	}
}

func (m *CanalMonitorService) newMyEventHandler() *CustomEventHandler {
	// 把schema和table正则合成一个正则表达式列表给IncludeTableRegex
	var compiledRegexps []*regexp.Regexp
	// 表格正则对应的监控规则
	rules := make(map[int][]edit.MonitorRuler, len(m.cfg.WatchHandlers))
	for i, wt := range m.cfg.WatchHandlers {
		r, err := regexp.Compile(wt.TableRegex)
		if err != nil {
			panic("编译正则失败: " + err.Error())
		}
		compiledRegexps = append(compiledRegexps, r)
		// 如果没有规则,使用默认规则
		if len(wt.Rules) == 0 {
			slog.Error(fmt.Sprintf("表 %s 没有监控规则,使用默认监控规则", wt.TableRegex))
			rules[i] = []edit.MonitorRuler{m.SSERule.Rule}
			continue
		} else {
			tableRules := make([]edit.MonitorRuler, len(wt.Rules))
			for j, ruleName := range wt.Rules {
				switch ruleName {
				case web.RuleName:
					tableRules[j] = m.SSERule.Rule
				case mgrpc.RuleName:
					tableRules[j] = m.grpcRule.Rule
				default:
					panic("规则 " + ruleName + " 不存在,请检查配置")
				}
			}
			rules[i] = tableRules
		}
	}
	return &CustomEventHandler{
		WatchRegexps: compiledRegexps,
		Rules:        rules,
	}
}

// newCanalConfig 根据配置文件,创建canal.Config
func (m *CanalMonitorService) newCanalConfig() *canal.Config {
	canalCfg := canal.NewDefaultConfig()
	canalCfg.Addr = m.cfg.Database.Addr
	canalCfg.User = m.cfg.Database.User
	canalCfg.Password = m.cfg.Database.Password
	canalCfg.Flavor = m.cfg.Database.Flavor
	canalCfg.ServerID = m.cfg.Database.ServerID
	canalCfg.Dump.ExecutionPath = m.cfg.Database.DumpExecutionPath
	canalCfg.IncludeTableRegex = m.cfg.Database.IncludeTableRegex
	return canalCfg
}

// StartCanal 启动canal
func (m *CanalMonitorService) StartCanal() {
	canalCfg := m.newCanalConfig()
	// 创建canal实例
	c, err := canal.NewCanal(canalCfg)
	if err != nil {
		slog.Error(fmt.Sprintf("创建canal失败: %v", err))
		return
	}
	// 设置事件处理器
	c.SetEventHandler(m.newMyEventHandler())

	// 获取当前的binlog位置
	pos, err := c.GetMasterPos()
	if err != nil {
		slog.Error(fmt.Sprintf("获取masterPos失败: %v", err))
		return
	}

	if err := c.RunFrom(pos); err != nil {
		slog.Error(fmt.Sprintf("canal运行失败: %v", err))
	}
}
