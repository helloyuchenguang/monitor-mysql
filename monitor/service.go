package monitor

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"log/slog"
	"main/common/config"
	"main/common/event/rule"
	"main/rules/meili"
	"main/rules/mgrpc"
	"main/rules/web"
	"math"
	"regexp"
)

type CanalMonitorService struct {
	cfg          *Config
	sseService   *web.SSERuleService
	grpcService  *mgrpc.GRPCRuleService
	meiliService *meili.ClientService
}

type WatchHandler struct {
	tableRegexp *regexp.Regexp
	rules       []string
}

type Database struct {
	addr              string
	user              string
	password          string
	flavor            string
	serverID          uint32
	dumpExecutionPath string
	includeTableRegex []string
}

type Config struct {
	Database      *Database
	WatchHandlers []*WatchHandler
}

// NewMonitorConfig 根据配置文件创建监控配置
func NewMonitorConfig(cfg *config.Config) *Config {
	cfgWatchHandlers := cfg.WatchHandlers
	// 创建配置
	var watchHandlers []*WatchHandler
	for _, handler := range cfgWatchHandlers {
		watchHandlers = append(watchHandlers, &WatchHandler{
			tableRegexp: handler.TableRegexp,
			rules:       handler.Rules,
		})
	}
	database := &Database{
		addr:              cfg.Database.Addr,
		user:              cfg.Database.User,
		password:          cfg.Database.Password,
		flavor:            cfg.Database.Flavor,
		serverID:          cfg.Database.ServerID,
		dumpExecutionPath: cfg.Database.DumpExecutionPath,
		includeTableRegex: cfg.Database.IncludeTableRegex,
	}
	return &Config{
		Database:      database,
		WatchHandlers: watchHandlers,
	}
}

// NewMonitorService 创建一个新的CanalMonitorService实例
func NewMonitorService(cfg *config.Config,
	sseService *web.SSERuleService,
	grpcService *mgrpc.GRPCRuleService,
	meiliService *meili.ClientService) *CanalMonitorService {
	return &CanalMonitorService{
		cfg:          NewMonitorConfig(cfg),
		sseService:   sseService,
		grpcService:  grpcService,
		meiliService: meiliService,
	}
}

// newMyEventHandler 创建自定义事件处理器
func (m *CanalMonitorService) newMyEventHandler() *CustomEventHandler {
	// 把schema和table正则合成一个正则表达式列表给IncludeTableRegex
	watchRegexps := make([]*WatchRegexp, len(m.cfg.WatchHandlers))
	// 表格正则对应的监控规则
	for i, wt := range m.cfg.WatchHandlers {
		rules := make([]rule.MonitorRuler, int(math.Max(float64(len(wt.rules)), 1)))
		for j, ruleName := range wt.rules {
			switch ruleName {
			case web.RuleName:
				if m.sseService == nil {
					panic("SSE规则服务未初始化")
				}
				rules[j] = m.sseService.Rule
			case mgrpc.RuleName:
				if m.grpcService == nil {
					panic("gRPC规则服务未初始化")
				}
				rules[j] = m.grpcService.Rule
			case meili.RuleName:
				if m.meiliService == nil {
					panic("MeiliSearch规则服务未初始化")
				}
				rules[j] = m.meiliService.Rule
			default:
				panic("规则 " + ruleName + " 不存在,请检查配置")
			}

		}
		watchRegexps[i] = &WatchRegexp{
			Regexp: wt.tableRegexp,
			Rules:  rules,
		}
	}
	return &CustomEventHandler{
		WatchRegexps: watchRegexps,
	}
}

// newCanalConfig 根据配置文件,创建canal.Config
func (m *CanalMonitorService) newCanalConfig() *canal.Config {
	canalCfg := canal.NewDefaultConfig()
	canalCfg.Addr = m.cfg.Database.addr
	canalCfg.User = m.cfg.Database.user
	canalCfg.Password = m.cfg.Database.password
	canalCfg.Flavor = m.cfg.Database.flavor
	canalCfg.ServerID = m.cfg.Database.serverID
	canalCfg.Dump.ExecutionPath = m.cfg.Database.dumpExecutionPath
	canalCfg.IncludeTableRegex = m.cfg.Database.includeTableRegex
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
	slog.Info("启动canal....")
	if err := c.RunFrom(pos); err != nil {
		slog.Error(fmt.Sprintf("canal运行失败: %v", err))
	}
}
