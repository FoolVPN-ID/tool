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
	if PROXY_IP_REGEXP.MatchString(messageText) {
		bot.localTemp.matchedText = messageText
		bot.handlers.proxyipCheck(bot, update)
	} else if CONFIG_VPN_REGEXP.MatchString(messageText) {
		bot.localTemp.matchedText = messageText
		bot.handlers.configSubconverter(bot, update)
		bot.handlers.configRegioncheck(bot, update)
	} else {
		bot.handlers.defaultHandler(bot, update)
	}
}
