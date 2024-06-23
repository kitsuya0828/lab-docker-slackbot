package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/Kitsuya0828/lab-docker-slackbot/client/config"
	"github.com/Kitsuya0828/lab-docker-slackbot/client/middleware"
	"github.com/Kitsuya0828/lab-docker-slackbot/client/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

var configPath = flag.String("config", "config.yaml", "path to config file")

func main() {
	flag.Parse()
	if err := config.LoadConfig(*configPath); err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	client, err := slack.NewSocketmodeClient()
	if err != nil {
		slog.Error("failed to connect to slack", "error", err)
		os.Exit(1)
	}

	handler := socketmode.NewSocketmodeHandler(client)

	handler.HandleEvents(slackevents.AppHomeOpened, middleware.AppHomeOpened)
	handler.HandleEvents(slackevents.AppMention, middleware.AppMention)

	handler.RunEventLoop()
}
