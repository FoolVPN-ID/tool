package subconverter

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/FoolVPN-ID/tool/common"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json"
)

var localConfig = make(map[string]string)

func (subconv *subconverterStruct) toSingboxByBaseConfig(configURL string) (option.Options, error) {
	// Make md5 of configURL
	var (
		configKeyHash = md5.Sum([]byte(configURL))
		configKey     = hex.EncodeToString(configKeyHash[:])

		buf = new(strings.Builder)
	)

	if localConfig[configKey] == "" {
		httpClient := common.MakeHTTPClient()
		req, _ := http.NewRequest("GET", configURL, nil)

		res, err := httpClient.Do(req)
		if err != nil {
			return option.Options{}, err
		}

		if res.StatusCode == 200 {
			reqBuf := new(strings.Builder)
			io.Copy(reqBuf, res.Body)
			localConfig[configKey] = reqBuf.String()
		} else {
			return option.Options{}, errors.New(res.Status)
		}
	}

	io.Copy(buf, strings.NewReader(localConfig[configKey]))
	baseConfig := buf.String()
	baseOptions, err := json.UnmarshalExtended[option.Options]([]byte(baseConfig))
	if err != nil {
		return option.Options{}, err
	}

	for _, newOutbound := range subconv.Outbounds {
		for i := range baseOptions.Outbounds {
			var oldOutbound = &baseOptions.Outbounds[i]

			switch oldOutbound.Tag {
			case "Internet", "Lock Region ID": // selector
				oldOutbound.SelectorOptions.Outbounds = append(oldOutbound.SelectorOptions.Outbounds, newOutbound.Tag)
			case "Best Latency": // url-test
				oldOutbound.URLTestOptions.Outbounds = append(oldOutbound.URLTestOptions.Outbounds, newOutbound.Tag)
			}
		}
		baseOptions.Outbounds = append(baseOptions.Outbounds, newOutbound)
	}

	ctx, cancel := context.WithCancel(context.Background())
	instance, err := box.New(box.Options{
		Context: ctx,
		Options: baseOptions,
	})
	if err == nil {
		instance.Close()
	}
	cancel()

	return baseOptions, nil
}
