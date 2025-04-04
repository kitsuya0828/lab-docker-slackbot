package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"sync"

	"github.com/Kitsuya0828/lab-docker-slackbot/client/config"
	"github.com/Kitsuya0828/lab-docker-slackbot/client/view"
	pb "github.com/Kitsuya0828/lab-docker-slackbot/proto/stat"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getUserName(evt *socketmode.Event) (string, error) {
	apiEvt, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		return "", fmt.Errorf("failed to convert socketmode.Event to slackevents.EventsAPIEvent")
	}

	var user string
	if openedEvt, ok := apiEvt.InnerEvent.Data.(slackevents.AppHomeOpenedEvent); !ok {
		user = reflect.ValueOf(apiEvt.InnerEvent.Data).Elem().FieldByName("User").Interface().(string)
	} else {
		user = openedEvt.User
	}
	return user, nil
}

func AppHomeOpened(evt *socketmode.Event, clt *socketmode.Client) {
	clt.Ack(*evt.Request)
	user, err := getUserName(evt)
	if err != nil {
		slog.Error("failed to get user name", "error", err)
		return
	}
	slog.Info("app home opened", "user", user)

	ctx := context.Background()
	fsStats := []*pb.FsStat{}
	dockerStats := []*pb.DockerStat{}
	hostnames := []string{}
	mu := sync.Mutex{}

	wg := sync.WaitGroup{}
	for _, host := range config.Cfg.Hosts {
		wg.Add(1)
		target := fmt.Sprintf("%s:%s", host.Address, host.Port)
		go func(ctx context.Context, target string) {
			defer wg.Done()

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
		}(ctx, target)
	}
	wg.Wait()

	v := view.HomeTabView(hostnames, fsStats, dockerStats)
	if _, err := clt.PublishView(user, v, ""); err != nil {
		slog.Error("failed to publish home tab view", "error", err)
	}
}

func AppMention(evt *socketmode.Event, clt *socketmode.Client) {
	clt.Ack(*evt.Request)
	apiEvt, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		slog.Error("failed to convert socketmode.Event to slackevents.EventsAPIEvent")
		return
	}

	mentionedEvt, ok := apiEvt.InnerEvent.Data.(*slackevents.AppMentionEvent)
	if !ok {
		slog.Error("failed to convert slackevents.EventsAPIEvent.InnerEvent.Data to slackevents.AppMentionEvent")
		return
	}

	channel := mentionedEvt.Channel
	user := mentionedEvt.User
	slog.Info("app mentioned", "user", user)

	mu := sync.Mutex{}
	results := []*pb.GetReccomendationResponse{}
	ctx := context.Background()

	wg := sync.WaitGroup{}
	for _, host := range config.Cfg.Hosts {
		wg.Add(1)
		target := fmt.Sprintf("%s:%s", host.Address, host.Port)
		go func(ctx context.Context, target string) {
			defer wg.Done()

			conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				slog.Error("failed to connect to grpc server", "error", err, "target", target)
				return
			}
			defer conn.Close()
			c := pb.NewStatServiceClient(conn)

			resp, err := c.GetReccomendation(ctx, &pb.GetReccomendationRequest{})
			if err != nil {
				slog.Error("failed to get reccomendation", "error", err, "target", target)
				return
			}

			mu.Lock()
			results = append(results, resp)
			mu.Unlock()
		}(ctx, target)
	}
	wg.Wait()

	msg := fmt.Sprintf("<@%s> %v", user, results)
	_, _, err := clt.PostMessage(channel, slack.MsgOptionText(msg, false))
	if err != nil {
		slog.Error("failed to post message", "error", err, "channel", channel, "message", msg)
		return
	}
}
