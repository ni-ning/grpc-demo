package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"

	pb "grpc-demo/proto"

	"google.golang.org/grpc"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8000", "启动端口号")
	flag.Parse()
}

// GreeterServer 实现proto定义的接口，go语言常见思路定义 struct 上面附着实现的方法即可
type GreeterServer struct {
	pb.UnimplementedGreeterServer
}

func (s *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello World"}, nil
}

func (s *GreeterServer) SayList(r *pb.HelloRequest, stream pb.Greeter_SayListServer) error {
	for n := 0; n <= 6; n++ {
		_ = stream.Send(&pb.HelloReply{Message: "hello.list"})
	}
	return nil
}

func (s *GreeterServer) SayRecord(stream pb.Greeter_SayRecordServer) error {
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			message := &pb.HelloReply{Message: "say.record"}
			return stream.SendAndClose(message)
		}
		if err != nil {
			return nil
		}
		log.Printf("resp: %v", resp)
	}
	return nil
}

func (s *GreeterServer) SayRoute(stream pb.Greeter_SayRouteServer) error {
	n := 0
	for {
		_ = stream.Send(&pb.HelloReply{Message: "say.route"})
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		n++

		log.Printf("resp: %v", resp)
	}
	return nil
}

func main() {
	server := grpc.NewServer()
	// 将实现好的GreeterServer，注册到gRPC内部的中心，请求时根据内部的"服务发现"，发现该服务端的接口，并进行逻辑处理
	pb.RegisterGreeterServer(server, &GreeterServer{})
	// 创建listen，监听TCP端口号
	lis, _ := net.Listen("tcp", ":"+port)
	// gRPC Server 开始lis.Accept, 直到Stop或GracefulStop
	server.Serve(lis)
}
