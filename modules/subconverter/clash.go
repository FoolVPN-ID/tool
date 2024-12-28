package subconverter

import (
	"gopkg.in/yaml.v3"
)

type clashProxiesMapping map[string]any

func (subconv *subconverterStruct) ToClash() error {
	var clashProxies = clashProxiesMapping{}

	clashProxies["proxies"] = subconv.Proxies
	out, err := yaml.Marshal(&clashProxies)
	if err != nil {
		return err
	}

	/**
	TODO
	- Build full clash meta config
	- Test config
	*/

	subconv.Result.Clash = string(out)
	return nil
}
