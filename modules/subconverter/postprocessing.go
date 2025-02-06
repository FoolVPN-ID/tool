package subconverter

import (
	"fmt"
	"slices"

	"github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badoption"
	"gopkg.in/yaml.v3"
)

const udpAccount = "trojan://t.me%2Ffoolvpn@172.67.73.39:443?path=%2Ftrojan-udp&security=tls&host=id1.foolvpn.me&type=ws&sni=id1.foolvpn.me#Trojan%20UDP"

func (subconv *subconverterStruct) PostTemplateSingBox(template string, singboxConfig option.Options) option.Options {
	var (
		udpSubconv, _      = MakeSubconverterFromConfig(udpAccount)
		udpOutbound        = udpSubconv.Outbounds[0]
		udpOutboundOptions = udpOutbound.Options.(option.TrojanOutboundOptions)
	)

	if template == "cf" {
		// Get used server address
		udpOutboundOptions.Server = subconv.Proxies[len(subconv.Proxies)-1]["server"].(string)
		udpOutbound.Options = udpOutboundOptions
		singboxConfig.Outbounds = append(singboxConfig.Outbounds, udpOutbound)

		// Configure dns
		for i := range singboxConfig.DNS.Servers {
			var dnsServer = &singboxConfig.DNS.Servers[i]

			if dnsServer.Type == "udp" {
				dnsServerOptions := dnsServer.Options.(*option.RemoteDNSServerOptions)
				dnsServerOptions.Detour = udpOutbound.Tag
				dnsServer.Options = dnsServerOptions
			}
		}

		// Configure rules
		for i := range singboxConfig.Route.Rules {
			rule := &singboxConfig.Route.Rules[i]
			if rule.Type == constant.RuleTypeDefault {
				if rule.DefaultOptions.Port != nil || rule.DefaultOptions.PortRange != nil {
					switch rule.DefaultOptions.RouteOptions.Outbound {
					case constant.TypeDirect, constant.TypeBlock, "dns-out":
					default:
						rule.DefaultOptions.Network = badoption.Listable[string]{"tcp"}
					}
				}
			}
		}
		singboxConfig.Route.Rules = append(singboxConfig.Route.Rules, option.Rule{
			Type: constant.RuleTypeDefault,
			DefaultOptions: option.DefaultRule{
				RawDefaultRule: option.RawDefaultRule{
					Network: badoption.Listable[string]{"udp"},
				},
				RuleAction: option.RuleAction{
					Action: "route",
					RouteOptions: option.RouteActionOptions{
						Outbound: udpOutbound.Tag,
					},
				},
			},
		})
	}

	return singboxConfig
}

func (subconv *subconverterStruct) PostTemplateClash(template string, clashConfig map[string]any) map[string]any {
	if template == "cf" {
		var (
			udpSubconv, _ = MakeSubconverterFromConfig(udpAccount)
			udpProxy      = udpSubconv.Proxies[0]
		)

		// Manipulate proxies
		newClashProxies := []map[string]any{}
		newClashProxiesByte, err := yaml.Marshal(clashConfig["proxies"])
		if err != nil {
			panic(err)
		}

		yaml.Unmarshal(newClashProxiesByte, &newClashProxies)
		udpProxy["server"] = newClashProxies[len(newClashProxies)-1]["server"]
		newClashProxies = append(newClashProxies, udpProxy)

		// Manipulate rules
		newClashRules := []string{}
		newClashRulesByte, err := yaml.Marshal(clashConfig["rules"])
		if err != nil {
			panic(err)
		}

		yaml.Unmarshal(newClashRulesByte, &newClashRules)
		newClashRules = slices.Insert(newClashRules, len(newClashRules)-1, []string{fmt.Sprintf("NETWORK,UDP,%s", udpProxy["name"])}...)

		// Overwrite old map
		clashConfig["proxies"] = newClashProxies
		clashConfig["rules"] = newClashRules
	}

	return clashConfig
}
