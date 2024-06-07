package main

import (
	"log/slog"
	"os"

	"github.com/Kitsuya0828/lab-docker-slackbot/driver"
	"github.com/Kitsuya0828/lab-docker-slackbot/middleware"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	if err := godotenv.Load("./.env"); err != nil {
		slog.Error("failed to load .env file", "error", err.Error())
	}

	client, err := drivers.ConnectToSlackViaSocketmode()
	if err != nil {
		slog.Error("failed to connect to slack", "error", err.Error())
		os.Exit(1)
	}

	socketmodeHandler := socketmode.NewSocketmodeHandler(client)

	socketmodeHandler.HandleEvents(slackevents.AppHomeOpened, middleware.AppHomeOpened)

	socketmodeHandler.RunEventLoop()
}
