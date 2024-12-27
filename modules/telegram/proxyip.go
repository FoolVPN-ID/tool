package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"strings"

	"github.com/FoolVPN-ID/tool/common"
	"github.com/FoolVPN-ID/tool/constant"
	"github.com/NicoNex/echotron/v3"
)

var (
	configSample = "vless://20eb9c0d-3014-4e45-9866-0fc5d5374558@nautica.foolvpn.me:443?encryption=none&type=ws&host=nautica.foolvpn.me&path=%2F20.233.68.69-2053&security=tls&sni=nautica.foolvpn.me#1%20%F0%9F%87%A6%F0%9F%87%AA%20Microsoft%20Azure%20WS%20TLS%20[nautica]"
	servers      = []map[string]string{
		{
			"name":       "Nautica",
			"remark":     "nautica",
			"domain":     "nautica.foolvpn.me",
			"maintainer": "@d_fordlalatina",
		},
		{
			"name":       "JeelsBoobz",
			"remark":     "nuclear",
			"domain":     "nuclear.us.kg",
			"maintainer": "@BangJeels",
		},
	}
)

func (handler *updateHandlers) proxyipCheck(bot *botStruct, _ *echotron.Update) {
	proxyIP := bot.localTemp.matchedText

	buf := new(strings.Builder)
	client := common.MakeHTTPClient()

	res, err := client.Get("https://" + constant.PROXYIP_CHECK_DOMAIN + fmt.Sprintf("?ip=%s", proxyIP))
	if err != nil {
		bot.SendMessage(fmt.Sprintf("Error while testing proxyip: %s", err.Error()), bot.chatID, nil)
		return
	}

	if res.StatusCode == 200 {
		io.Copy(buf, res.Body)
		result := buf.String()

		var resultInJson map[string]interface{}
		err := json.Unmarshal([]byte(result), &resultInJson)
		if err != nil {
			bot.SendMessage("Failed while parsing json", bot.chatID, nil)
			return
		}

		if resultInJson["proxyip"] == false {
			bot.SendMessage("Proxy is inactive", bot.chatID, nil)
			return
		}

		var message string = "<b>TEST RESULT</b>\n"
		message += "<blockquote><code>"
		message += fmt.Sprintf("Active  : %v\n", resultInJson["proxyip"])
		message += fmt.Sprintf("IP      : %v\n", resultInJson["proxy"])
		message += fmt.Sprintf("Port    : %v\n", resultInJson["port"])
		message += fmt.Sprintf("ORG     : %v\n", resultInJson["asOrganization"])
		message += fmt.Sprintf("Ping    : %v ms\n", resultInJson["delay"])
		message += fmt.Sprintf("Country : %v\n", resultInJson["country"])
		message += fmt.Sprintf("Region  : %v\n", resultInJson["region"])
		message += fmt.Sprintf("City    : %v\n", resultInJson["city"])
		message += fmt.Sprintf("Colo    : %v\n", resultInJson["colo"])
		message += "</code></blockquote>\n\n"

		// Build config
		var (
			defaultConfig, _ = url.Parse(configSample)
			defaultQueries   = defaultConfig.Query()
		)
		defaultConfig.Path = fmt.Sprintf("/%v", proxyIP)

		for _, vpn := range []string{"vless", "trojan"} {
			// Select random cf server
			server := servers[rand.Intn(len(servers))]
			defaultQueries.Set("host", server["domain"])
			defaultQueries.Set("sni", server["domain"])
			defaultConfig.RawQuery = defaultQueries.Encode()
			defaultConfig.Fragment = fmt.Sprintf("%v %v [%v]", resultInJson["country"], resultInJson["asOrganization"], server["remark"])

			// Resume build config
			message += fmt.Sprintf("<b>%s</b>", strings.ToUpper(vpn))
			message += "<blockquote><code>"
			defaultConfig.Scheme = vpn
			message += defaultConfig.String()
			message += "</code></blockquote>\n"
			message += fmt.Sprintf("Credit: %v\n\n", server["maintainer"])
		}

		bot.SendMessage(message, bot.chatID, &echotron.MessageOptions{
			ParseMode: echotron.HTML,
		})
	} else {
		bot.SendMessage(fmt.Sprintf("Error return code %d", res.StatusCode), bot.chatID, nil)
	}
}
