package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/nakano0518/grpc-quickstart"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func unaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	//RPCの前後で処理を挟むInterceptor(ログ出力)(他、共通処理や認証処理などを記述するのに便利)
	log.Printf("before call: %s, request: %+V", method, req)
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("after call: %s, response: %+v", method, reply)
	return err
}

func main() {
	addr := "localhost:50051"

	/*
		//TLS認証で接続
			creds, err := credentials.NewClientTLSFromFile("server.crt", "")
			if err != nil {
				log.Fatal(err)
			}
			conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	*/

	/*
		//認証せずに接続
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
	*/

	//認証せずに接続時Interceptorを使用
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(unaryInterceptor))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	name := os.Args[1]
	md := metadata.Pairs("timestamp", time.Now().Format(time.Stamp)) //metadata(HTTP/2のHeaderに格納し運ばれる)//key-value方式(文字列型)//圧縮形式の指定や認証情報の受け渡しなどに使用
	//ctx := context.Background() //キャンセルやタイムアウトを考慮せずずっと待ち続ける(Goの仕組み)
	ctx, cancel := context.WithCancel(context.Background()) //キャンセル可能に
	defer cancel()
	ctx = metadata.NewOutgoingContext(ctx, md)
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name}, grpc.Trailer(&md))
	if err != nil {
		//エラーの詳細情報(grpc-status-details-bin)
		s, ok := status.FromError(err) //エラーからgrpcのstatusを取り出す
		if ok {
			log.Printf("gRPC Error (message: %s)", s.Message())
			for _, d := range s.Details() { //エラーの詳細情報は、status.Details()に格納
				switch info := d.(type) {
				case *errdetails.RetryInfo:
					log.Printf(" RetryInfo: %v", info)
				}
			}
			os.Exit(1)
		} else {
			log.Fatalf("could not greet: %v", err)
		}
	}
	log.Printf("Greeting: %s", r.Message)
}
