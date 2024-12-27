package telegram

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/FoolVPN-ID/tool/modules/subconverter"
	"github.com/NicoNex/echotron/v3"
)

func (handler *updateHandlers) configSubconverter(bot *botStruct, _ *echotron.Update) {
	rawConfig := bot.localTemp.matchedText
	subconv, err := subconverter.MakeSubconverterFromConfig(rawConfig)
	if err != nil {
		panic(err)
	}

	// Convert
	subconv.ToSFA()
	subconv.ToBFR()

	// Send result
	var (
		sfaByte, _ = json.MarshalIndent(subconv.Result.SFA, "", "  ")
		sfaFile    = echotron.NewInputFileBytes(fmt.Sprintf("SFA_%v.txt", time.Now().Unix()), sfaByte)

		bfrByte, _ = json.MarshalIndent(subconv.Result.BFR, "", "  ")
		bfrFile    = echotron.NewInputFileBytes(fmt.Sprintf("BFR_%v.txt", time.Now().Unix()), bfrByte)
	)

	bot.SendMediaGroup(bot.chatID, []echotron.GroupableInputMedia{
		echotron.InputMediaDocument{
			Type:  echotron.MediaTypeDocument,
			Media: sfaFile,
		},
		echotron.InputMediaDocument{
			Type:  echotron.MediaTypeDocument,
			Media: bfrFile,
		},
	}, nil)
}
