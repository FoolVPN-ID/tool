package subconverter

import (
	"github.com/FoolVPN-ID/tool/constant"
)

func (subconv *subconverterStruct) ToSFA() error {
	options, err := subconv.toSingboxByBaseConfig(constant.SFA_BASE_CONFIG)
	if err != nil {
		return err
	}

	subconv.Result.SFA = options
	return nil
}
