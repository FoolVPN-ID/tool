package config

import (
	"context"
	"errors"
	"net/netip"

	"github.com/FoolVPN-ID/tool/common"
	"github.com/FoolVPN-ID/tool/modules/provider"
	box "github.com/sagernet/sing-box"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/include"
	"github.com/sagernet/sing-box/option"
	CM "github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badoption"
)

func BuildSingboxConfig(rawConfig string) (option.Options, error) {
	ctx := context.Background()
	ctx = box.Context(ctx, include.InboundRegistry(), include.OutboundRegistry(), include.EndpointRegistry(), include.DNSTransportRegistry(), include.ServiceRegistry())

	outbounds, err := provider.Parse(rawConfig)
	if err != nil {
		return option.Options{}, err
	}
	if len(outbounds) == 0 {
		return option.Options{}, errors.New("parsing failed")
	}

	var outboundsAny []any
	outboundsByte, _ := json.Marshal(outbounds)
	json.Unmarshal(outboundsByte, &outboundsAny)

	config := map[string]any{
		"log": map[string]any{
			"disabled": true,
		},
		"dns": map[string]any{
			"servers": []map[string]any{
				{
					"tag":         "default-dns",
					"type":        "udp",
					"server":      "1.1.1.1",
					"server_port": 53,
					"detour":      "direct",
				},
			},
			"final": "default-dns",
		},
		"inbounds": []map[string]any{
			{
				"tag":         "mixed-in",
				"type":        C.TypeMixed,
				"listen":      CM.Ptr(badoption.Addr(netip.IPv4Unspecified())),
				"listen_port": uint16(common.GetFreePort()),
			},
		},
		"outbounds": []map[string]any{
			{
				"tag":  "direct",
				"type": C.TypeDirect,
			},
		},
		"route": map[string]any{
			"rules": []map[string]any{
				{
					"type": C.RuleTypeLogical,
					"mode": "or",
					"rules": []map[string]any{
						{
							"protocol": "dns",
						},
						{
							"port": 53,
						},
					},
					"action": "hijack-dns",
				},
				{
					"type":     C.RuleTypeDefault,
					"network":  "udp",
					"action":   "route",
					"outbound": "direct",
				},
			},
			"final": outbounds[0].Tag,
		},
	}

	for _, outbound := range outboundsAny {
		config["outbounds"] = append([]map[string]any{outbound.(map[string]any)}, config["outbounds"].([]map[string]any)...)
	}

	configByte, _ := json.Marshal(config)
	return json.UnmarshalExtendedContext[option.Options](ctx, configByte)
}
