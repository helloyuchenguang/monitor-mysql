package mrpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log/slog"
	"monitormysql/mrpc/api/mycanal"
	"net"
)

type MyCanalServer struct {
	mycanal.UnimplementedMyCanalServiceServer
}

func RunGrpcCanal() {
	// 创建tcp监听
	listen, err := net.Listen("tcp", ":18081")
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	// 创建grpc服务器
	s := grpc.NewServer()
	// 实例化MyCanalServer
	server := MyCanalServer{}
	// 注册MyCanalServer
	mycanal.RegisterMyCanalServiceServer(s, &server)
	// 启动grpc服务器
	err = s.Serve(listen)
	if err != nil {
		return
	}

}

func (s *MyCanalServer) SubscribeRegexTable(req *mycanal.SubscribeTableRequest,
	stream mycanal.MyCanalService_SubscribeRegexTableServer) error {
	// 这里可以实现具体的逻辑
	// 例如，发送一些事件到客户端
	slog.Info("接收到订阅请求", req.TableNameRegex)
	return nil
}
