package mgrpc

import (
	"github.com/samber/lo"
	"main/common/event"
	"main/rules/mgrpc/api/mycanal"
)

// EventDataToGrpcReply 从行数据中获取更新信息
func EventDataToGrpcReply(data *event.Data) *mycanal.EventTableRowReply {
	table := data.Table
	editData := data.EditData
	editColumns := make(map[string]*mycanal.ColumnValueReply, len(editData.EditFieldValues))
	for k, v := range editData.EditFieldValues {
		be, af := v.ConvertToPairStr()
		editColumns[k] = &mycanal.ColumnValueReply{
			BeforeValue: be,
			AfterValue:  af,
		}
	}
	return &mycanal.EventTableRowReply{
		EventType: mycanal.CanalEventType_UPDATE,
		Table: &mycanal.TableStructReply{
			Schema: table.Schema,
			Table:  table.Table,
			Columns: lo.Map(table.Columns, func(item *event.Column, _ int) *mycanal.ColumnInfoReply {
				return &mycanal.ColumnInfoReply{
					ColumnName: item.Name,
					ColumnType: item.RowType,
				}
			}),
		},
		EditData: &mycanal.EventUpdateInfoReply{
			UnChangeRowData: editData.UnChangeRowData.ConvertMapStrStr(),
			EditColumns:     editColumns,
		},
	}
}
