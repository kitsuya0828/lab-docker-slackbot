package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/Kitsuya0828/lab-docker-slackbot/disk"
	"github.com/Kitsuya0828/lab-docker-slackbot/docker"
	pb "github.com/Kitsuya0828/lab-docker-slackbot/proto/stat"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.StatServiceServer
}

func (s *server) GetFsStat(ctx context.Context, in *pb.GetFsStatRequest) (*pb.GetFsStatResponse, error) {
	log.Printf("Received: %v", in.String())
	du, err := disk.GetDiskUsage("/")
	if err != nil {
		return nil, err
	}
	return &pb.GetFsStatResponse{
		FsStat: &pb.FsStat{
			Total: du.Total,
			Used:  du.Used,
			Free:  du.Free,
		},
	}, nil
}

func (s *server) GetDockerStat(ctx context.Context, in *pb.GetDockerStatRequest) (*pb.GetDockerStatResponse, error) {
	log.Printf("Received: %v", in.String())
	cli, err := docker.NewClient()
	if err != nil {
		return nil, err
	}
	du, err := docker.GetDiskUsage(context.Background(), cli)
	if err != nil {
		return nil, err
	}
	return &pb.GetDockerStatResponse{
		DockerStat: &pb.DockerStat{
			Images: &pb.DockerStat_Item{
				Active:      du.Images.Active,
				Size:        float32(du.Images.Size),
				Reclaimable: float32(du.Images.Reclaimable),
				TotalCount:  du.Images.TotalCount,
			},
			Containers: &pb.DockerStat_Item{
				Active:      du.Containers.Active,
				Size:        float32(du.Containers.Size),
				Reclaimable: float32(du.Containers.Reclaimable),
				TotalCount:  du.Containers.TotalCount,
			},
			LocalVolumes: &pb.DockerStat_Item{
				Active:      du.LocalVolumes.Active,
				Size:        float32(du.LocalVolumes.Size),
				Reclaimable: float32(du.LocalVolumes.Reclaimable),
				TotalCount:  du.LocalVolumes.TotalCount,
			},
			BuildCache: &pb.DockerStat_Item{
				Active:      du.BuildCache.Active,
				Size:        float32(du.BuildCache.Size),
				Reclaimable: float32(du.BuildCache.Reclaimable),
				TotalCount:  du.BuildCache.TotalCount,
			},
		},
	}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterStatServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
