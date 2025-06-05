package mgrpc

import (
	"github.com/samber/lo"
	"main/common/event"
	"main/rules/mgrpc/api/mycanal"
)

// EventDataToGrpcReply 从行数据中获取更新信息
func EventDataToGrpcReply(data *event.Data) *mycanal.EventTableRowReply {
	table := data.Table
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
		EditData: lo.Map(data.EditData, func(item *event.EditData, _ int) *mycanal.EventUpdateInfoReply {
			editColumns := make(map[string]*mycanal.ColumnValueReply, len(item.EditFieldValues))
			for k, v := range item.EditFieldValues {
				be, af := v.ConvertToPairStr()
				editColumns[k] = &mycanal.ColumnValueReply{
					BeforeValue: be,
					AfterValue:  af,
				}
			}
			return &mycanal.EventUpdateInfoReply{
				UnChangeRowData: item.UnChangeRowData.ConvertMapStrStr(),
				EditColumns:     editColumns,
			}
		}),
	}
}
