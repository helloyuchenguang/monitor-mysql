package web

import (
	"log/slog"
	"monitormysql/global"
	"monitormysql/global/mevent"
	"net/http"
)

const ruleName = "SSERule"

// TplNodeModel 模板节点模型
type TplNodeModel struct {
	Id                        string `column:"id"`
	AccKey                    string `column:"acc_key"`
	BusinessCode              string `column:"business_code"`
	CreateTime                string `column:"create_time"`
	UpdateTime                string `column:"update_time"`
	IsDeleted                 int    `column:"is_deleted"`
	OwnerUser                 string `column:"owner_user"`
	Title                     string `column:"title"`
	Executor                  string `column:"executor"`
	Participator              string `column:"participator"`
	NodeState                 string `column:"node_state"`
	EstimateStartDate         string `column:"estimate_start_date"`
	EstimateEndDate           string `column:"estimate_end_date"`
	EstimatePeriod            string `column:"estimate_period"`
	ActualStartDate           string `column:"actual_start_date"`
	ActualEndDate             string `column:"actual_end_date"`
	ActualPeriod              string `column:"actual_period"`
	SerialNumber              string `column:"serial_number"`
	ParentNodeBusinessCode    string `column:"parent_node_business_code"`
	ParentNodeDataId          string `column:"parent_node_data_id"`
	ParentNodeDataTitle       string `column:"parent_node_data_title"`
	ParentSerialNumbers       string `column:"parent_serial_numbers"`
	BelongProjectBusinessCode string `column:"belong_project_business_code"`
	BelongProjectDataId       string `column:"belong_project_data_id"`
	BelongProjectDataTitle    string `column:"belong_project_data_title"`
	BelongProjectBusinessName string `column:"belong_project_business_name"`
	ProjectState              string `column:"project_state"`
	ProjectChargeUser         string `column:"project_charge_user"`
	ProjectSubChargeUser      string `column:"project_sub_charge_user"`
	ProjectMembers            string `column:"project_members"`
	IsTplNode                 int    `column:"is_tpl_node"`
}

// SSERule SSE 规则处理器
type SSERule struct {
	SseServer *SSEServer
}

func (h *SSERule) OnChange(ui mevent.UpdateInfo) error {
	// 发布更新信息到所有 SSE 客户端
	h.SseServer.Broadcast(ui)
	return nil
}

func (h *SSERule) AddSseClient(w http.ResponseWriter, r *http.Request) {
	h.SseServer.AddClient(w, r)
}

func GetSSERule() *SSERule {
	if rule, ok := global.GetRuleByName(ruleName); ok {
		return (*rule).(*SSERule)
	}
	return nil
}

// 自动注册
func init() {
	slog.Info("自动注册 SSE 规则处理器")

	global.RegisterRule(ruleName, &SSERule{
		SseServer: NewSSEServer(),
	})
}
