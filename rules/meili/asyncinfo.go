package meili

import (
	"github.com/meilisearch/meilisearch-go"
	"github.com/samber/lo"
	"log/slog"
	"main/common/event"
	"main/common/event/rule"
	"sync"
)

type ClientService struct {
	Client        meilisearch.ServiceManager
	Rule          *rule.Server
	ChannelClient *rule.ChannelClient
	IndexMap      sync.Map // map[string]struct{}
}

// 启动数据监听与同步
func (cs *ClientService) asyncDataChange() {
	ch := cs.ChannelClient.Chan
	for {
		select {
		case evt, ok := <-ch:
			if !ok {
				slog.Info("MeiliSearch监听通道已关闭")
				return
			}
			index, err := cs.CreateIndexOrIgnore(evt.Table)
			if err != nil {
				slog.Error("创建或获取MeiliSearch索引失败", slog.String("table", evt.Table.ObtainTableName()), slog.Any("error", err))
				return
			}
			switch evt.EventType {
			case event.Insert:
				_, err := cs.Client.Index(index).AddDocuments(lo.Map(evt.SaveData, func(item *event.SaveData, _ int) event.RowData {
					return item.RowData
				}))
				if err != nil {
					slog.Error("插入MeiliSearch文档失败", err)
					return
				}
			case event.Delete:
				_, err := cs.Client.Index(index).DeleteDocuments(lo.Map(evt.DeleteData, func(item *event.DeleteData, _ int) string {
					return item.RowData.PrimaryKey()
				}))
				if err != nil {
					slog.Error("删除MeiliSearch文档失败", err)
					return
				}
			case event.Update:
				_, err := cs.Client.Index(index).UpdateDocuments(lo.Map(evt.EditData, func(item *event.EditData, _ int) event.RowData {
					rowData := item.UnChangeRowData
					for f, v := range item.EditFieldValues {
						rowData[f] = v.After
					}
					return rowData
				}))
				if err != nil {
					slog.Error("修改MeiliSearch文档失败", err)
					return
				}
			}
		}
	}
}

// CreateIndexOrIgnore 创建索引（若不存在）
func (cs *ClientService) CreateIndexOrIgnore(table *event.Table) (string, error) {
	tableName := table.ObtainTableName()

	if _, ok := cs.IndexMap.Load(tableName); ok {
		return tableName, nil
	}

	// double check
	_, loaded := cs.IndexMap.LoadOrStore(tableName, struct{}{})
	if loaded {
		return tableName, nil
	}

	task, err := cs.Client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        tableName,
		PrimaryKey: "id",
	})
	if err != nil {
		slog.Error("创建MeiliSearch索引失败", slog.String("index", tableName), slog.Any("error", err))
		cs.IndexMap.Delete(tableName) // 创建失败清理
		return "", err
	}

	slog.Info("成功创建MeiliSearch索引", slog.String("index", task.IndexUID))
	return task.IndexUID, nil
}

// 插入文档
func (cs *ClientService) insertDataToMeili(index string, data *event.SaveData) error {
	_, err := cs.Client.Index(index).AddDocuments([]any{data.RowData})
	return err
}

// 更新文档
func (cs *ClientService) updateDataToMeili(index string, data *event.EditData) error {
	rowData := data.UnChangeRowData
	for f, v := range data.EditFieldValues {
		rowData[f] = v.After
	}
	_, err := cs.Client.Index(index).UpdateDocuments([]any{rowData})
	return err
}

// 删除文档
func (cs *ClientService) deleteDataToMeili(index string, data *event.DeleteData) error {
	_, err := cs.Client.Index(index).DeleteDocument(data.RowData.PrimaryKey())
	if err != nil {
		slog.Error("删除MeiliSearch文档失败", slog.String("index", index), slog.Any("error", err))
	}
	return err
}
