package meili

import (
	"github.com/samber/lo"
	"log/slog"
	"main/common/event"
	"time"
)

const (
	bufferSize = 1000
	interval   = 2 * time.Second
)

// asyncDataChange 启动数据监听与同步
func (cs *ClientService) asyncDataChange() {
	slog.Info("MeiliSearchRule启动,监听地址", slog.String("addr", cs.config.addr))
	ch := cs.channelClient.Chan
	// 确保在函数结束时清理资源
	defer cs.Rule.RemoveClientByID(cs.channelClient.ID)
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

// flushToMeiliSearch 将数据列表同步到 MeiliSearch
func (cs *ClientService) flushToMeiliSearch(dataList []*event.Data) {
	docsMap := make(map[string][]event.RowData)
	deleteIDsMap := make(map[string][]string)
	// 遍历数据列表，分类处理
	for _, data := range dataList {
		index, err := cs.CreateIndexOrIgnore(data.Table)
		if err != nil {
			continue
		}
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
