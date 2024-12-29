package subconverter

import (
	"io"
	"net/http"
	"strings"

	"github.com/FoolVPN-ID/tool/common"
	"github.com/FoolVPN-ID/tool/constant"
	"gopkg.in/yaml.v3"
)

func (subconv *subconverterStruct) ToClash() error {
	var (
		clashConfig      = map[string]any{}
		clashProxyGroups = []map[string]any{
			{
				"name": "Tunnel",
				"type": "select",
				"proxies": []string{
					"Url Test",
					"Selector",
				},
			},
			{
				"name":     "Url Test",
				"type":     "url-test",
				"interval": 300,
				"proxies":  []string{},
			},
			{
				"name":     "Selector",
				"type":     "select",
				"interval": 300,
				"proxies":  []string{},
			},
		}
	)

	httpClient := common.MakeHTTPClient()
	req, _ := http.NewRequest("GET", constant.MIHOMO_BASE_CONFIG, nil)
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == 200 {
		buf := new(strings.Builder)
		io.Copy(buf, res.Body)

		yaml.Unmarshal([]byte(buf.String()), &clashConfig)
		proxyTags := []string{}
		for _, proxy := range subconv.Proxies {
			proxyTags = append(proxyTags, proxy["name"].(string))
		}

		for i := range clashProxyGroups {
			proxyGroup := clashProxyGroups[i]
			switch proxyGroup["name"].(string) {
			case "Tunnel":
			default:
				proxyGroup["proxies"] = proxyTags
			}
		}
	}

	clashConfig["proxy-groups"] = clashProxyGroups
	clashConfig["proxies"] = subconv.Proxies

	subconv.Result.Clash = clashConfig
	return nil
}
