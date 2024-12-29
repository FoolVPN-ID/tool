package subconverter

import (
	"fmt"
	"regexp"

	"github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"gopkg.in/yaml.v3"
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
				if rule.DefaultOptions.Port != nil || rule.DefaultOptions.PortRange != nil {
					switch rule.DefaultOptions.Outbound {
					case constant.TypeDirect, constant.TypeBlock, "dns-out":
					default:
						rule.DefaultOptions.Network = option.Listable[string]{"tcp"}
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
		newClashRules = append([]string{fmt.Sprintf("NETWORK,UDP,%s", udpProxy["name"])}, newClashRules...)

		// Overwrite old map
		clashConfig["proxies"] = newClashProxies
		clashConfig["rules"] = newClashRules
	}

	return clashConfig
}
