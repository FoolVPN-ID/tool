package subconverter

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/FoolVPN-ID/tool/modules/provider"
	"github.com/metacubex/mihomo/common/convert"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/include"
	"github.com/sagernet/sing-box/option"
)

var (
	globalCtx context.Context
)

type subconverterStruct struct {
	Outbounds  []option.Outbound
	Proxies    []map[string]any
	rawConfigs string
	Result     struct {
		Clash map[string]any
		Raw   []string
		SFA   option.Options
		BFR   option.Options
	}
}

func preRun() {
	globalCtx = context.Background()
	globalCtx = box.Context(globalCtx, include.InboundRegistry(), include.OutboundRegistry(), include.EndpointRegistry(), include.DNSTransportRegistry(), include.ServiceRegistry())
}

func MakeSubconverterFromConfig(config string) (subconverterStruct, error) {
	preRun()

	subconv := subconverterStruct{}
	subconv.rawConfigs = strings.ReplaceAll(config, ",", "\n")
	subconv.Outbounds = subconv.parse(subconv.rawConfigs)
	subconv.Proxies, _ = convert.ConvertsV2Ray([]byte(subconv.rawConfigs))

	if len(subconv.Outbounds) == 0 || len(subconv.Proxies) == 0 {
		return subconv, errors.New("configs not found")
	}
	if len(subconv.Outbounds)-len(subconv.Proxies) != 0 {
		return subconv, errors.New("subconverter result is not even")
	}

	return subconv, nil
}

func (subconv *subconverterStruct) parse(content string) []option.Outbound {
	out, err := provider.Parse(content)
	if err != nil {
		log.Println(err.Error())
	}

	return out
}
