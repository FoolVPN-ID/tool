package telegram

import (
	"fmt"

	"github.com/NicoNex/echotron/v3"
)

func (handler *updateHandlers) cmdStartHandler(bot *botStruct, update *echotron.Update) {
	var message string = fmt.Sprintf("Selamat datang %s\n\n", update.Message.From.FirstName)
	message += "Kirimkan proxy dalam format <code>1.1.1.1:443</code> untuk melihat status proxy.\n\n"
	message += "Bot ini masih dalam pengembangan, fitur lain akan hadir segera!\n"
	message += "Berikut adalah fitur yang direncanakan hadir:\n"
	message += "• Check proxy\n"
	message += "• Check streaming region\n"
	message += "• Subconverter\n"

	go bot.SendMessage(message, bot.chatID, &echotron.MessageOptions{
		ParseMode: echotron.HTML,
	})
}
