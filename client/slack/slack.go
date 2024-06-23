package slack

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/Kitsuya0828/lab-docker-slackbot/client/config"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func NewSocketmodeClient() (*socketmode.Client, error) {
	if !strings.HasPrefix(config.Cfg.AppToken, "xapp-") {
		return nil, errors.New("SLACK_APP_TOKEN must have the prefix \"xapp-\"")
	}

	if !strings.HasPrefix(config.Cfg.BotToken, "xoxb-") {
		return nil, errors.New("SLACK_BOT_TOKEN must have the prefix \"xoxb-\"")
	}

	api := slack.New(
		config.Cfg.BotToken,
		slack.OptionDebug(true),
		slack.OptionAppLevelToken(config.Cfg.AppToken),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)

	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	return client, nil
}
