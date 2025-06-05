package meili

import (
	"github.com/meilisearch/meilisearch-go"
	"github.com/samber/lo"
	"log/slog"
	"main/common/event"
	"time"
)

// asyncDataChange 启动数据监听与同步
func (cs *ClientService) asyncDataChange() {
	ch := cs.channelClient.Chan
	const (
		bufferSize = 1000
		interval   = 2 * time.Second
	)
	buffer := make([]*event.Data, 0, bufferSize)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	flush := func() {
		if len(buffer) == 0 {
			return
		}
		slog.Info("刷新 MeiliSearch 数据", slog.Int("数量", len(buffer)))
		cs.flushToMeiliSearch(buffer)
		buffer = buffer[:0]
	}

	for {
		select {
		case evt, ok := <-ch:
			if !ok {
				slog.Info("MeiliSearch监听通道已关闭")
				return
			}
			buffer = append(buffer, evt)
			if len(buffer) >= bufferSize {
				flush()
			}
		case <-ticker.C:
			flush()
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

	task, err := cs.CreateIndex(&meilisearch.IndexConfig{
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

// flushToMeiliSearch 将数据列表同步到 MeiliSearch
func (cs *ClientService) flushToMeiliSearch(dataList []*event.Data) {
	docsMap := make(map[string][]event.RowData)
	deleteIDsMap := make(map[string][]string)
	// 遍历数据列表，分类处理
	for _, data := range dataList {
		index := data.Table.ObtainTableName()
		switch data.EventType {
		case event.Insert:
			docsMap[index] = append(docsMap[index], lo.Map(data.SaveData, func(item *event.SaveData, _ int) event.RowData {
				return item.RowData
			})...)
		case event.Delete:
			docsMap[index] = append(docsMap[index], lo.Map(data.DeleteData, func(item *event.DeleteData, _ int) event.RowData {
				return item.RowData
			})...)
		case event.Update:
			docsMap[index] = append(docsMap[index], lo.Map(data.EditData, func(item *event.EditData, _ int) event.RowData {
				rowData := item.UnChangeRowData
				for f, v := range item.EditFieldValues {
					rowData[f] = v.After
				}
				return rowData
			})...)
		}
	}

	// 添加文档
	for index, docs := range docsMap {
		if len(docs) == 0 {
			continue
		}
		if _, err := cs.Index(index).AddDocuments(docs); err != nil {
			slog.Error("MeiliSearch添加文档失败", slog.String("index", index), slog.Any("error", err))
		}
	}

	// 删除文档
	for index, delIDs := range deleteIDsMap {
		if len(delIDs) == 0 {
			continue
		}
		if _, err := cs.Index(index).DeleteDocuments(delIDs); err != nil {
			slog.Error("MeiliSearch删除文档失败", slog.String("index", index), slog.Any("error", err))
		}
	}
}

// SettingsAttributes 更新索引的设置
func (cs *ClientService) SettingsAttributes(index string, searchers, filters, sorts []string) error {
	// 设置索引的可搜索属性
	settings := &meilisearch.Settings{
		SearchableAttributes: searchers,
		FilterableAttributes: filters,
		SortableAttributes:   sorts,
	}
	if _, err := cs.Index(index).UpdateSettings(settings); err != nil {
		slog.Error("MeiliSearch更新索引设置失败", slog.String("index", index), slog.Any("error", err))
		return err
	}
	slog.Info("成功更新MeiliSearch索引设置", slog.String("index", index))
	return nil
}

// SetExperimentalFeatures 设置实验性功能
func (cs *ClientService) SetExperimentalFeatures() error {
	// 这里可以添加更多实验性功能的设置
	features := cs.ExperimentalFeatures()
	// 开启CONTAINS 和 START WITH 功能
	features.SetContainsFilter(true)
	// 开启函数编辑文档
	features.SetEditDocumentsByFunction(true)
	return nil
}
