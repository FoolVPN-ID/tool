package telegram

import (
	"fmt"

	"github.com/FoolVPN-ID/tool/modules/regioncheck"
	"github.com/NicoNex/echotron/v3"
)

func (handler *updateHandlers) listenVPNConfigUpdate(bot *botStruct, _ *echotron.Update) {
	rawConfig := bot.localTemp.matchedText
	rc := regioncheck.MakeLibrary()
	err := rc.Run(rawConfig)
	if err != nil {
		bot.SendMessage(fmt.Sprintf("Error while performing region check: %v", err), bot.chatID, nil)
		return
	}

	var message string = "<b>REGION CHECK RESULT</b>\n\n"
	for _, data := range rc.Result {
		message += fmt.Sprintf("<b>%s</b>\n", data.Name)
		message += "<blockquote><code>"
		message += fmt.Sprintf("IATACode : %s\n", data.IATACode)
		message += fmt.Sprintf("Region   : %s\n", data.Region)
		message += fmt.Sprintf("Ping     : %d ms\n", data.Delay)
		message += "</code></blockquote>\n\n"
	}

	bot.SendMessage(message, bot.chatID, &echotron.MessageOptions{
		ParseMode: echotron.HTML,
	})
}
