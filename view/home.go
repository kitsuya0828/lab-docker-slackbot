package view

import (
	"github.com/slack-go/slack"
)

func HomeTabView() slack.HomeTabViewRequest {
	return slack.HomeTabViewRequest{
		Type: slack.VTHomeTab,
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewSectionBlock(
					slack.NewTextBlockObject(slack.MarkdownType, "Hello! I'm docker-bot, your friendly neighborhood bot.", false, false),
					nil,
					nil,
				),
				slack.NewDividerBlock(),
			},
		},
	}
}
