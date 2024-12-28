package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"strings"
	"sync"

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
	var (
		wg       sync.WaitGroup
		proxyIPs = strings.Split(bot.localTemp.matchedText, "\n")

		checkResults []map[string]any
	)

	for _, proxyip := range proxyIPs {
		wg.Add(1)
		go (func() {
			defer common.RecoverFromPanic()
			defer wg.Done()

			result, err := checkProxyIP(proxyip)
			if err != nil {
				bot.SendMessage(fmt.Sprintf("Error while checking %v: %v", proxyip, err.Error()), bot.chatID, nil)
				return
			}
			checkResults = append(checkResults, result)
		})()
	}

	// Wait for waitgroup
	wg.Wait()

	var message string = "<b>TEST RESULT</b>\n\n"
	for _, resultInJson := range checkResults {
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
	}

	if len(checkResults) == 1 {
		// Build config
		var (
			defaultConfig, _ = url.Parse(configSample)
			defaultQueries   = defaultConfig.Query()
			resultInJson     = checkResults[0]
		)
		defaultConfig.Path = fmt.Sprintf("/%v-%v", resultInJson["proxy"], resultInJson["port"])

		for _, vpn := range []string{"vless", "trojan"} {
			// Select random cf server
			server := servers[rand.Intn(len(servers))]
			defaultQueries.Set("host", server["domain"])
			defaultQueries.Set("sni", server["domain"])
			defaultConfig.Scheme = vpn
			defaultConfig.RawQuery = defaultQueries.Encode()
			defaultConfig.Fragment = fmt.Sprintf("%v %v [%v]", resultInJson["country"], resultInJson["asOrganization"], server["remark"])

			// Resume build config
			message += fmt.Sprintf("<b>%s</b>", strings.ToUpper(vpn))
			message += "<blockquote><code>"
			message += defaultConfig.String()
			message += "</code></blockquote>\n"
			message += fmt.Sprintf("Credit: %v\n\n", server["maintainer"])
		}
	}

	bot.SendMessage(message, bot.chatID, &echotron.MessageOptions{
		ParseMode: echotron.HTML,
	})
}

func checkProxyIP(proxyip string) (map[string]any, error) {
	var (
		resultInJson map[string]any
		buf          = new(strings.Builder)
		client       = common.MakeHTTPClient()
	)

	res, err := client.Get("https://" + constant.PROXYIP_CHECK_DOMAIN + fmt.Sprintf("?ip=%s", proxyip))
	if err != nil {
		return resultInJson, err
	}

	if res.StatusCode == 200 {
		io.Copy(buf, res.Body)
		result := buf.String()

		err := json.Unmarshal([]byte(result), &resultInJson)
		if err != nil {
			return resultInJson, err
		}

		return resultInJson, nil
	}

	return resultInJson, err
}
