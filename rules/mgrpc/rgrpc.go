package mgrpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"main/rules/mgrpc/api/mycanal"
	"net"
)

func (g *GRPCRuleService) RunGrpcCanal() {
	// 创建tcp监听
	listen, err := net.Listen("tcp", g.cfg.Addr)
	if err != nil {
		grpclog.Fatalf("grpc tcp 错误: %v", err)
	}
	// 创建grpc服务器
	s := grpc.NewServer()
	// 实例化MyCanalServer
	server := MyCanalServer{grpcRuleServer: g}
	// 注册MyCanalServer
	mycanal.RegisterMyCanalServiceServer(s, &server)
	// 启动grpc服务器
	err = s.Serve(listen)
	if err != nil {
		return
	}

}
