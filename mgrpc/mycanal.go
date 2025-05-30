package mgrpc

import (
	"log/slog"
	"main/common"
	"main/mgrpc/api/mycanal"
)

type MyCanalServer struct {
	mycanal.UnimplementedMyCanalServiceServer
}

func (s *MyCanalServer) SubscribeRegexTable(req *mycanal.SubscribeTableRequest, stream mycanal.MyCanalService_SubscribeRegexTableServer) error {
	rs := common.GetRule[*mycanal.EventTableRowReply](RuleName)
	slog.Info("接收到订阅的表名regex", slog.String("regex", req.TableNameRegex))

	client := rs.PutNewClient()
	clientID := client.ID
	defer rs.RemoveClientByID(clientID)

	for {
		select {
		case <-stream.Context().Done():
			slog.Warn("grpc客户端断开连接", slog.String("clientID", clientID))
			return nil
		case evt, ok := <-client.Chan:
			if !ok {
				slog.Info("通道关闭", slog.String("clientID", clientID))
				return nil
			}
			if err := stream.Send(evt); err != nil {
				slog.Error("grpc推送消息失败", slog.String("clientID", clientID), slog.Any("error", err))
				return err
			}
		}
	}
}
