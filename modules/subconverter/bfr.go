package subconverter

import (
	"github.com/FoolVPN-ID/tool/constant"
)

func (subconv *subconverterStruct) ToBFR() error {
	options, err := subconv.toSingboxByBaseConfig(constant.BFR_BASE_CONFIG)
	if err != nil {
		return err
	}

	subconv.Result.BFR = options
	return nil
}
