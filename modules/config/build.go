package config

import (
	"errors"
	"net/netip"

	"github.com/FoolVPN-ID/tool/common"
	"github.com/FoolVPN-ID/tool/modules/provider"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	CM "github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/json/badoption"
)

func BuildSingboxConfig(rawConfig string) (option.Options, error) {
	outbound, err := provider.Parse(rawConfig)
	if err != nil {
		return option.Options{}, err
	}
	if len(outbound) == 0 {
		return option.Options{}, errors.New("parsing failed")
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
				Options: option.HTTPMixedInboundOptions{
					ListenOptions: option.ListenOptions{
						Listen:     CM.Ptr(badoption.Addr(netip.IPv4Unspecified())),
						ListenPort: uint16(common.GetFreePort()),
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
					LogicalOptions: option.LogicalRule{
						RawLogicalRule: option.RawLogicalRule{
							Mode:   "or",
							Invert: false,
							Rules: []option.Rule{
								{
									DefaultOptions: option.DefaultRule{
										RawDefaultRule: option.RawDefaultRule{
											Protocol: badoption.Listable[string]{"dns"},
										},
									},
								},
								{
									DefaultOptions: option.DefaultRule{
										RawDefaultRule: option.RawDefaultRule{
											Port: badoption.Listable[uint16]{53},
										},
									},
								},
							},
						},
						RuleAction: option.RuleAction{
							Action: "hijack-dns",
						},
					},
				},
				{
					Type: C.RuleTypeDefault,
					DefaultOptions: option.DefaultRule{
						RawDefaultRule: option.RawDefaultRule{
							Network: badoption.Listable[string]{"udp"},
						},
						RuleAction: option.RuleAction{
							Action: "route",
							RouteOptions: option.RouteActionOptions{
								Outbound: "direct",
							},
						},
					},
				},
			},
			Final: outbound[0].Tag,
		},
	}

	return config, nil
}
