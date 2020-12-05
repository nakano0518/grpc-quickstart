package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/golang/protobuf/ptypes/duration"
	pb "github.com/nakano0518/grpc-quickstart"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	time.Sleep(3 * time.Second) //server側で時間のかかる処理を想定(client側でキャンセル処理を実装)
	//return &pb.HelloReply{Message: "Hello " + in.Name}, nil
	//return nil, status.New(codes.NotFound, "resource not found").Err() //Headerにgrpc-status,grpc-messageとして格納
	st, _ := status.New(codes.Aborted, "aborted").WithDetails(&errdetails.RetryInfo{ //grpc-status,grpc-message以外にエラー時の詳細情報をHeaderのgrpc-status-details-binに追加(ここでは3秒後にリトライ)//エラーになったコードの箇所なども情報に含めることもできる
		RetryDelay: &duration.Duration{
			Seconds: 3,
			Nanos:   0,
		},
	})
	return nil, st.Err()
}

func main() {
	addr := ":50051"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("gRPC server listening on" + addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
