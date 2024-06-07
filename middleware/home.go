package middleware

import (
	"log/slog"
	"reflect"

	"github.com/Kitsuya0828/lab-docker-slackbot/view"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func AppHomeOpened(evt *socketmode.Event, clt *socketmode.Client) {
	slog.Info("App Home Opened")
	apiEvt, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		slog.Error("faield to convert socketmode.Event to slackevents.EventsAPIEvent")
	}

	var user string
	openedEvt, ok := apiEvt.InnerEvent.Data.(slackevents.AppHomeOpenedEvent)
	if !ok {
		slog.Error("failed to convert slackevents.EventsAPIEvent.InnerEvent.Data to slackevents.AppHomeOpenedEvent")
		user = reflect.ValueOf(apiEvt.InnerEvent.Data).Elem().FieldByName("User").Interface().(string)
	} else {
		user = openedEvt.User
	}

	view := view.HomeTabView()
	_, err := clt.PublishView(user, view, "")
	if err != nil {
		slog.Error("failed to publish home tab view", "error", err)
	}
}
