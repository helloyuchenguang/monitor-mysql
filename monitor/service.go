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

type MonitorService struct {
	SSERule  *web.SSERuleService
	grpcRule *mgrpc.GRPCRuleServer
}

func (m *MonitorService) NewMyEventHandler(cfg *config.Config) *MyEventHandler {
	canalConfig := NewCanalConfig(cfg)
	// 把schema和table正则合成一个正则表达式列表给IncludeTableRegex
	var compiledRegexps []*regexp.Regexp
	// 表格正则对应的监控规则
	rules := make(map[int][]edit.MonitorRuler, len(cfg.WatchHandlers))
	for i, wt := range cfg.WatchHandlers {
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
			for _, ruleName := range wt.Rules {
				switch ruleName {
				case web.RuleName:
					tableRules = append(tableRules, m.SSERule.Rule)
				case mgrpc.RuleName:
					tableRules = append(tableRules, m.grpcRule.Rule)
				default:
					panic("规则 " + ruleName + " 不存在,请检查配置")
				}
			}
		}
	}
	return &MyEventHandler{
		WatchRegexps: compiledRegexps,
	}
}

func NewMonitorService(cfg config.Config, sseRule *web.SSERuleService, grpcRule *mgrpc.GRPCRuleServer) *MonitorService {

	// 初始化SSE服务
	return &MonitorService{
		SSERule:  sseRule,
		grpcRule: grpcRule,
	}
}

// NewCanalConfig 根据配置文件,创建canal.Config
func NewCanalConfig(cfg *config.Config) *canal.Config {
	canalCfg := canal.NewDefaultConfig()
	canalCfg.Addr = cfg.Database.Addr
	canalCfg.User = cfg.Database.User
	canalCfg.Password = cfg.Database.Password
	canalCfg.Flavor = cfg.Database.Flavor
	canalCfg.ServerID = cfg.Database.ServerID
	canalCfg.Dump.ExecutionPath = cfg.Database.DumpExecutionPath
	canalCfg.IncludeTableRegex = cfg.Database.IncludeTableRegex
	return canalCfg
}

// StartCanal 启动canal
func StartCanal(canalCfg *canal.Config, handler *MyEventHandler) {
	// 创建canal实例
	c, err := canal.NewCanal(canalCfg)
	if err != nil {
		slog.Error(fmt.Sprintf("创建canal失败: %v", err))
		return
	}
	// 设置事件处理器
	c.SetEventHandler(handler)

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
