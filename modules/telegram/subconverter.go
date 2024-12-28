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
	subconv.ToClash()
	subconv.ToSFA()
	subconv.ToBFR()

	// Send result
	var (
		uniqueID = time.Now().Unix()

		clashByte = []byte(subconv.Result.Clash)
		clashFile = echotron.NewInputFileBytes(fmt.Sprintf("Clash_%v.txt", uniqueID), clashByte)

		sfaByte, _ = json.MarshalIndent(subconv.Result.SFA, "", "  ")
		sfaFile    = echotron.NewInputFileBytes(fmt.Sprintf("SFA_%v.txt", uniqueID), sfaByte)

		bfrByte, _ = json.MarshalIndent(subconv.Result.BFR, "", "  ")
		bfrFile    = echotron.NewInputFileBytes(fmt.Sprintf("BFR_%v.txt", uniqueID), bfrByte)
	)

	bot.SendMediaGroup(bot.chatID, []echotron.GroupableInputMedia{
		echotron.InputMediaDocument{
			Type:  echotron.MediaTypeDocument,
			Media: clashFile,
		},
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
