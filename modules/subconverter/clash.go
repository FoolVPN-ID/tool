package subconverter

import (
	"github.com/metacubex/mihomo/common/convert"
	"gopkg.in/yaml.v3"
)

type clashProxiesMapping map[string]any

func (subconv *subconverterStruct) ToClash() error {
	var clashProxies = clashProxiesMapping{}
	results, err := convert.ConvertsV2Ray([]byte(subconv.rawConfigs))
	if err != nil {
		return err
	}

	clashProxies["proxies"] = results
	out, err := yaml.Marshal(&clashProxies)
	if err != nil {
		return err
	}

	subconv.Result.Clash = string(out)
	return nil
}
