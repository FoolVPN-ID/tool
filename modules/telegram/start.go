package telegram

import (
	"fmt"

	"github.com/NicoNex/echotron/v3"
)

func (handler *updateHandlers) defaultHandler(bot *botStruct, update *echotron.Update) {
	var message string = fmt.Sprintf("Selamat datang %s\n\n", update.Message.From.FirstName)
	message += "• Kirimkan proxy dalam format <code>1.1.1.1:443</code> untuk melihat status proxy.\n"
	message += "• Kirimkan akun VPN untuk merubahnya ke bentuk sing-box, dan clash\n\n"
	message += "<b>Catatan</b>\n"
	message += "• Region check hanya akan berjalan jika kamu mengirimkan 1 akun VPN\n"
	message += "• Laporan, saran, dan request fitur silahkan ke @d_fordlalatina"

	go bot.SendMessage(message, bot.chatID, &echotron.MessageOptions{
		ParseMode: echotron.HTML,
	})
}
