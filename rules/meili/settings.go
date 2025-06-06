package meili

import (
	"github.com/meilisearch/meilisearch-go"
	"github.com/samber/lo"
	"log/slog"
	"main/common/event"
)

// CreateIndexOrIgnore 创建索引（若不存在）
func (cs *ClientService) CreateIndexOrIgnore(table *event.Table) (string, error) {
	tableName := table.ObtainTableName()

	if index, ok := cs.IndexMap.Load(tableName); ok {
		return index.(string), nil
	}

	// 将索引名称存入IndexMap
	index, err := cs.SetAttributesByIndexConfig(tableName)
	if err != nil {
		return "", err
	}
	// 存储索引名称到IndexMap
	cs.IndexMap.Store(tableName, index)
	return index, nil
}

func (cs *ClientService) SetAttributesByIndexConfig(tableName string) (string, error) {
	if indexCfg, ok := lo.Find(cs.config.indexConfigs, func(item *IndexConfig) bool {
		return item.tableRegex.MatchString(tableName)
	}); ok {
		index := indexCfg.index
		task, err := cs.CreateIndex(&meilisearch.IndexConfig{
			Uid:        index,
			PrimaryKey: "id",
		})
		if err != nil {
			slog.Error("创建MeiliSearch索引失败", slog.String("index", tableName), slog.Any("error", err))
			cs.IndexMap.Delete(tableName) // 创建失败清理
			return "", err
		}
		slog.Info("成功创建MeiliSearch索引", slog.String("index", task.IndexUID))
		settings := &meilisearch.Settings{
			SearchableAttributes: indexCfg.searchers,
			FilterableAttributes: indexCfg.filters,
			SortableAttributes:   indexCfg.sorts,
		}
		if _, err := cs.Index(index).UpdateSettings(settings); err != nil {
			slog.Error("MeiliSearch更新索引设置失败", slog.String("index", index), slog.Any("error", err))
			return "", err
		}
		slog.Info("成功更新MeiliSearch索引设置", slog.String("index", index))
		return index, nil
	} else {
		slog.Error("未找到匹配的索引配置", slog.String("tableName", tableName))
		return "", nil
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

// DeleteAllIndexes 删除所有索引
func (cs *ClientService) DeleteAllIndexes() error {
	// 删除所有索引
	indexes, err := cs.ListIndexes(&meilisearch.IndexesQuery{
		Limit: 1000, // 设置一个合理的限制
	})
	if err != nil {
		slog.Error("获取MeiliSearch索引失败", slog.Any("error", err))
		return err
	}
	for _, index := range indexes.Results {
		if _, err := cs.DeleteIndex(index.UID); err != nil {
			slog.Error("删除MeiliSearch索引失败", slog.String("index", index.UID), slog.Any("error", err))
			return err
		}
		slog.Info("成功删除MeiliSearch索引", slog.String("index", index.UID))
	}
	return nil
}
