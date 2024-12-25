package config

import (
	"errors"
	"fmt"
	"net/netip"

	"github.com/LalatinaHub/LatinaSub-go/helper"
	"github.com/LalatinaHub/LatinaSub-go/provider"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

func BuildSingboxConfig(rawConfig string) (option.Options, error) {
	fmt.Println(rawConfig)
	outbound, err := provider.Parse(rawConfig)
	if err != nil {
		return option.Options{}, err
	}
	if len(outbound) == 0 {
		return option.Options{}, errors.New("PARSING_FAILED")
	}

	config := option.Options{
		Log: &option.LogOptions{
			Disabled: true,
		},
		DNS: &option.DNSOptions{
			Servers: []option.DNSServerOptions{
				{
					Tag:     "default-dns",
					Address: "1.1.1.1",
					Detour:  "direct",
				},
			},
			Final: "default-dns",
		},
		Inbounds: []option.Inbound{
			{
				Type: C.TypeMixed,
				MixedOptions: option.HTTPMixedInboundOptions{
					ListenOptions: option.ListenOptions{
						Listen:     option.NewListenAddress(netip.IPv4Unspecified()),
						ListenPort: uint16(helper.GetFreePort()),
					},
				},
			},
		},
		Outbounds: []option.Outbound{
			outbound[0],
			{
				Tag:  "direct",
				Type: C.TypeDirect,
			},
		},
		Route: &option.RouteOptions{
			Rules: []option.Rule{
				{
					Type: C.RuleTypeDefault,
					DefaultOptions: option.DefaultRule{
						Protocol: option.Listable[string]{"dns"},
						Outbound: "direct",
					},
				},
				{
					Type: C.RuleTypeDefault,
					DefaultOptions: option.DefaultRule{
						Network:  option.Listable[string]{"udp"},
						Outbound: "direct",
					},
				},
			},
			Final: outbound[0].Tag,
		},
	}

	return config, nil
}
