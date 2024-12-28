package subconverter

import (
	"regexp"
	"slices"

	"github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

const udpAccount = "trojan://t.me%2Ffoolvpn@172.67.73.39:443?path=%2Ftrojan-udp&security=tls&host=id1.foolvpn.me&type=ws&sni=id1.foolvpn.me#Trojan%20UDP"

func (subconv *subconverterStruct) PostTemplateSingBox(template string, singboxConfig option.Options) option.Options {
	var (
		udpSubconv, _ = MakeSubconverterFromConfig(udpAccount)
		udpOutbound   = udpSubconv.Outbounds[0]
	)

	if template == "cf" {
		// Get used server address
		udpOutbound.TrojanOptions.Server = subconv.Proxies[len(subconv.Proxies)-1]["server"].(string)
		singboxConfig.Outbounds = append(singboxConfig.Outbounds, udpOutbound)

		// Configure dns
		for i := range singboxConfig.DNS.Servers {
			dnsServer := &singboxConfig.DNS.Servers[i]
			if regexp.MustCompile(`^\d`).MatchString(dnsServer.Address) {
				if dnsServer.Detour != constant.TypeDirect {
					dnsServer.Detour = udpOutbound.Tag
				}
			}
		}

		// Configure rules
		for i := range singboxConfig.Route.Rules {
			rule := &singboxConfig.Route.Rules[i]
			if rule.Type == constant.RuleTypeDefault {
				if slices.Contains(rule.DefaultOptions.Network, "udp") {
					switch rule.DefaultOptions.Outbound {
					case constant.TypeBlock, constant.TypeDirect:
					default:
						rule.DefaultOptions.Outbound = udpOutbound.Tag
					}
				}
			}
		}
		singboxConfig.Route.Rules = append(singboxConfig.Route.Rules, option.Rule{
			Type: constant.RuleTypeDefault,
			DefaultOptions: option.DefaultRule{
				Network:  option.Listable[string]{"udp"},
				Outbound: udpOutbound.Tag,
			},
		})
	}

	return singboxConfig
}
