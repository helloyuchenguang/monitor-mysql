package monitormysql

import (
	"net/http"
)

// Package monitor-mysql 监控 mysql 数据库变更

// MonitorHandler 监控处理器接口
type MonitorHandler interface {
	// OnChange 处理更新信息
	OnChange(ui UpdateInfo) error
	// AddSseClient 添加 SSE 客户端
	AddSseClient(w http.ResponseWriter, r *http.Request)
}

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

type TplNodeHandler struct {
	SseServer *SSEServer
}

func (h *TplNodeHandler) OnChange(ui UpdateInfo) error {
	// 处理更新信息
	//fmt.Printf("%-v\n", ui)
	//oldRaw, newRaw, err := ConvertByUpdateInfo[TplNodeModel](ui)
	//if err != nil {
	//	slog.Error(fmt.Sprintf("转换更新信息失败: %v", err))
	//	return err
	//}
	h.SseServer.Broadcast(ui)
	return nil
}

func (h *TplNodeHandler) AddSseClient(w http.ResponseWriter, r *http.Request) {
	h.SseServer.AddClient(w, r)
}
