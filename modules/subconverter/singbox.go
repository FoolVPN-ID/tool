package subconverter

import (
	"io"
	"net/http"
	"strings"

	"github.com/FoolVPN-ID/tool/common"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json"
)

func (subconv *subconverterStruct) toSingboxByBaseConfig(configURL string) (option.Options, error) {
	var buf = new(strings.Builder)

	httpClient := common.MakeHTTPClient()
	req, _ := http.NewRequest("GET", configURL, nil)

	res, err := httpClient.Do(req)
	if err != nil {
		return option.Options{}, err
	}

	if res.StatusCode == 200 {
		io.Copy(buf, res.Body)
	}

	var baseConfig = buf.String()

	baseOptions, err := json.UnmarshalExtendedContext[option.Options](globalCtx, []byte(baseConfig))
	if err != nil {
		return option.Options{}, err
	}

	for _, newOutbound := range subconv.Outbounds {
		for i := range baseOptions.Outbounds {
			var oldOutbound = &baseOptions.Outbounds[i]

			switch oldOutbound.Tag {
			case "Internet", "Lock Region ID": // selector
				selectorOptions := oldOutbound.Options.(*option.SelectorOutboundOptions)
				selectorOptions.Outbounds = append(selectorOptions.Outbounds, newOutbound.Tag)
				oldOutbound.Options = selectorOptions
			case "Best Latency": // url-test
				urlTestOptions := oldOutbound.Options.(*option.URLTestOutboundOptions)
				urlTestOptions.Outbounds = append(urlTestOptions.Outbounds, newOutbound.Tag)
				oldOutbound.Options = urlTestOptions
			}
		}
		baseOptions.Outbounds = append(baseOptions.Outbounds, newOutbound)
	}

	return baseOptions, nil
}
