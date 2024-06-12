package main

import (
	"flag"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/Kitsuya0828/lab-docker-slackbot/client/slack"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

var (
	envFile    = flag.String("env", ".env", "path to .env file")
	configFile = flag.String("config", "config.yaml", "path to config file")

	cfg *config
)

type Host struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
}

type config struct {
	Hosts []Host `yaml:"hosts"`
}

func main() {
	flag.Parse()
	if err := godotenv.Load(*envFile); err != nil {
		slog.Error("failed to load .env file", "error", err)
	}

	b, err := os.ReadFile(*configFile)
	if err != nil {
		slog.Error("failed to read config file", "error", err)
	}
	cfg = &config{}
	if err := yaml.Unmarshal(b, cfg); err != nil {
		slog.Error("failed to unmarshal yaml", "error", err)
	}

	client, err := slack.ConnectToSlackViaSocketmode()
	if err != nil {
		slog.Error("failed to connect to slack", "error", err)
		os.Exit(1)
	}

	socketmodeHandler := socketmode.NewSocketmodeHandler(client)

	socketmodeHandler.HandleEvents(slackevents.AppHomeOpened, middlewareAppHomeOpened)
	socketmodeHandler.HandleEvents(slackevents.AppMention, middlewareAppMentioned)

	socketmodeHandler.RunEventLoop()
}
