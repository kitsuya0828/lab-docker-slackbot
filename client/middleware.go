package main

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"sync"

	pb "github.com/Kitsuya0828/lab-docker-slackbot/proto/stat"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func middlewareAppHomeOpened(evt *socketmode.Event, clt *socketmode.Client) {
	apiEvt, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		slog.Error("faield to convert socketmode.Event to slackevents.EventsAPIEvent")
	}

	var user string
	if openedEvt, ok := apiEvt.InnerEvent.Data.(slackevents.AppHomeOpenedEvent); !ok {
		user = reflect.ValueOf(apiEvt.InnerEvent.Data).Elem().FieldByName("User").Interface().(string)
	} else {
		user = openedEvt.User
	}
	slog.Info("app home opened", "user", user)

	ctx := context.Background()
	fsStats := []*pb.FsStat{}
	dockerStats := []*pb.DockerStat{}
	hostnames := []string{}
	mu := sync.Mutex{}

	wg := sync.WaitGroup{}
	for _, host := range cfg.Hosts {
		wg.Add(1)
		go func(ctx context.Context, host Host) {
			defer wg.Done()

			target := fmt.Sprintf("%s:%s", host.Address, host.Port)
			conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				slog.Error("failed to connect to grpc server", "error", err, "target", target)
			}
			defer conn.Close()
			c := pb.NewStatServiceClient(conn)

			fs, err := c.GetFsStat(ctx, &pb.GetFsStatRequest{})
			if err != nil {
				slog.Error("failed to get fs stat", "error", err, "target", target)
			}
			ds, err := c.GetDockerStat(ctx, &pb.GetDockerStatRequest{})
			if err != nil {
				slog.Error("failed to get docker stat", "error", err, "target", target)
			}

			mu.Lock()
			fsStats = append(fsStats, fs.FsStat)
			dockerStats = append(dockerStats, ds.DockerStat)
			hostnames = append(hostnames, fs.Hostname)
			mu.Unlock()
		}(ctx, host)
	}
	wg.Wait()

	v := HomeTabView(hostnames, fsStats, dockerStats)
	if _, err := clt.PublishView(user, v, ""); err != nil {
		slog.Error("failed to publish home tab view", "error", err)
	}
}
