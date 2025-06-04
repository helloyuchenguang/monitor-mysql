package meili

import (
	"github.com/meilisearch/meilisearch-go"
	cmap "github.com/orcaman/concurrent-map/v2"
	"log/slog"
	"main/common/event"
	"main/common/event/rule"
)

type ClientService struct {
	Client        meilisearch.ServiceManager
	Rule          *rule.Server
	ChannelClient *rule.ChannelClient
	IndexMap      cmap.ConcurrentMap[string, bool]
}

// asyncDataChange 同步数据变化
func (cs *ClientService) asyncDataChange() {
	ch := cs.ChannelClient.Chan
	for {
		select {
		case evt, ok := <-ch:
			if !ok {
				return // 通道已关闭
			}
			// 创建索引或忽略
			index, err := cs.CreateIndexOrIgnore(evt.Table)
			if err != nil {
				slog.Error("创建索引失败", slog.String("index", evt.Table.ObtainTableName()), slog.Any("error", err))
				return
			}
			// 根据事件类型处理数据变化
			switch evt.EventType {
			case event.Insert:
				err = cs.insertDataToMeili(index, evt.SaveData)
				if err != nil {
					slog.Error("插入数据到MeiliSearch失败", slog.Any("error", err))
				}
			case event.Update:
				err := cs.updateDataToMeili(index, evt.EditData)
				if err != nil {
					slog.Error("更新数据到MeiliSearch失败", slog.Any("error", err))
				}
			case event.Delete:
				err := cs.deleteDataToMeili(evt)
				if err != nil {
					slog.Error("删除数据到MeiliSearch失败", slog.Any("error", err))
				}
			}
		}
	}
}

func (cs *ClientService) CreateIndexOrIgnore(table *event.Table) (string, error) {
	tableName := table.ObtainTableName()
	if cs.IndexMap.Has(tableName) {
		return tableName, nil
	}
	indexConfig := &meilisearch.IndexConfig{
		Uid:        tableName,
		PrimaryKey: "id",
	}
	// 创建索引
	task, err := cs.Client.CreateIndex(indexConfig)
	if err != nil {
		slog.Error("创建MeiliSearch索引失败", slog.String("index", tableName), slog.Any("error", err))
		return "", err
	}
	slog.Info("成功创建MeiliSearch索引", slog.String("index", task.IndexUID))
	cs.IndexMap.Set(tableName, true)
	return task.IndexUID, nil
}

func (cs *ClientService) insertDataToMeili(index string, data *event.SaveData) error {
	_, err := cs.Client.Index(index).AddDocuments([]any{data.RowData})
	if err != nil {
		return err
	}
	return nil
}

func (cs *ClientService) updateDataToMeili(index string, data *event.EditData) error {
	rowData := data.UnChangeRowData
	for f, v := range data.EditFieldValues {
		rowData[f] = v.After
	}
	_, err := cs.Client.Index(index).UpdateDocuments([]any{rowData})
	if err != nil {
		return err
	}
	return nil
}

func (cs *ClientService) deleteDataToMeili(evt *event.Data) error {
	return nil
}
