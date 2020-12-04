package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/nakano0518/grpc-quickstart"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	addr := "localhost:50051"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	name := os.Args[1]
	md := metadata.Pairs("timestamp", time.Now().Format(time.Stamp)) //metadata(HTTP/2のHeaderに格納し運ばれる)//key-value方式(文字列型)//圧縮形式の指定や認証情報の受け渡しなどに使用
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name}, grpc.Trailer(&md))
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}
