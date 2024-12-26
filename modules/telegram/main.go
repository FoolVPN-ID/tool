package telegram

import (
	"context"
	"fmt"
	"os"

	"github.com/FoolVPN-ID/tool/common"
	"github.com/NicoNex/echotron/v3"
)

type updateHandlers struct{}

type botStruct struct {
	chatID    int64
	handlers  updateHandlers
	localTemp struct {
		matchedText string
	}
	echotron.API
}

var token = os.Getenv("TELEGRAM_BOT_TOKEN")

func newBot(chatID int64) echotron.Bot {
	return &botStruct{
		chatID: chatID,
		API:    echotron.NewAPI(token),
	}
}

func RunWithContext(ctx context.Context) {
	dsp := makeDispatcher()
	go dsp.Poll()

	<-ctx.Done()
	fmt.Println("Shutting down telegram bot...")
}

func makeDispatcher() *echotron.Dispatcher {
	return echotron.NewDispatcher(token, newBot)
}

func (bot *botStruct) Update(update *echotron.Update) {
	// Error handler
	defer common.RecoverFromPanic()

	// Defers
	defer bot.SetMessageReaction(bot.chatID, update.Message.ID, &echotron.MessageReactionOptions{
		Reaction: []echotron.ReactionType{
			{
				Type:  "emoji",
				Emoji: "👍",
			},
		},
	})

	// Middlewares
	go bot.SendChatAction(echotron.Typing, bot.chatID, nil)

	// Update handlers
	var messageText = update.Message.Text
	if messageText == "/start" {
		bot.handlers.cmdStartHandler(bot, update)
	} else if proxyIP := PROXY_IP_REGEXP.FindString(messageText); proxyIP != "" {
		bot.localTemp.matchedText = proxyIP
		bot.handlers.listenProxyIPUpdate(bot, update)
	} else if rawConfig := CONFIG_VPN_REGEXP.FindString(messageText); rawConfig != "" {
		bot.localTemp.matchedText = rawConfig
		fmt.Println(rawConfig)
	} else {
		bot.handlers.cmdStartHandler(bot, update)
	}
}
