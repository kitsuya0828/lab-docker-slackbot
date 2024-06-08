package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/Kitsuya0828/lab-docker-slackbot/proto/stat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewStatServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r1, err := c.GetFsStat(ctx, &pb.GetFsStatRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Filesystem stats: %s", r1.String())

	r2, err := c.GetDockerStat(ctx, &pb.GetDockerStatRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Docker stats: %s", r2.String())
}
