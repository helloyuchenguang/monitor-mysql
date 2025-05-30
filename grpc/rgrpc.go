package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"main/grpc/api/mycanal"
	"net"
)

func RunGrpcCanal(address string) {
	// 创建tcp监听
	listen, err := net.Listen("tcp", address)
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
