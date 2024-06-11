package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sync"

	"github.com/Kitsuya0828/lab-docker-slackbot/disk"
	"github.com/Kitsuya0828/lab-docker-slackbot/docker"
	pb "github.com/Kitsuya0828/lab-docker-slackbot/proto/stat"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
	path = flag.String("path", "/", "The path to get disk usage")
)

type server struct {
	pb.StatServiceServer
	mu sync.Mutex
}

func (s *server) GetFsStat(ctx context.Context, in *pb.GetFsStatRequest) (*pb.GetFsStatResponse, error) {
	du, err := disk.GetDiskUsage(*path)
	if err != nil {
		return nil, err
	}
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	return &pb.GetFsStatResponse{
		FsStat: &pb.FsStat{
			Total: du.Total,
			Used:  du.Used,
			Free:  du.Free,
		},
		Path:     *path,
		Hostname: hostname,
	}, nil
}

func (s *server) GetDockerStat(ctx context.Context, in *pb.GetDockerStatRequest) (*pb.GetDockerStatResponse, error) {
	cli, err := docker.NewClient()
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	// To avoid "rpc error: code = Unknown desc = Error response from daemon: a disk usage operation is already running"
	s.mu.Lock()
	du, err := docker.GetDiskUsage(context.Background(), cli)
	if err != nil {
		return nil, err
	}
	s.mu.Unlock()

	hostname, err := os.Hostname()
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
		Hostname: hostname,
	}, nil
}

func (s *server) GetReccomendation(ctx context.Context, in *pb.GetReccomendationRequest) (*pb.GetReccomendationResponse, error) {
	cli, err := docker.NewClient()
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	rec, err := docker.GetReccomendation(context.Background(), cli)
	if err != nil {
		return nil, err
	}

	reccomendedImages := make([]*pb.ReccomendationItem, 0)
	for _, image := range rec.Images {
		reccomendedImages = append(reccomendedImages, &pb.ReccomendationItem{
			Id:   image.Id,
			User: image.User,
			Name: image.Name,
			Size: uint64(image.Size),
		})
	}
	reccomendedContainers := make([]*pb.ReccomendationItem, 0)
	for _, container := range rec.Containers {
		reccomendedContainers = append(reccomendedContainers, &pb.ReccomendationItem{
			Id:   container.Id,
			User: container.User,
			Name: container.Name,
			Size: uint64(container.Size),
		})
	}
	return &pb.GetReccomendationResponse{
		Images:     reccomendedImages,
		Containers: reccomendedContainers,
	}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		slog.Error("failed to listen", "error", err, "port", *port)
		os.Exit(1)
	}
	s := grpc.NewServer()
	pb.RegisterStatServiceServer(s, &server{})
	slog.Info("server listening", "port", *port)
	if err := s.Serve(lis); err != nil {
		slog.Error("failed to serve", "error", err)
		os.Exit(1)
	}
}
