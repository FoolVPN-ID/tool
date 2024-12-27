package subconverter

import "strings"

func (subconv *subconverterStruct) ToRaw() {
	subconv.Result.Raw = strings.Split(subconv.rawConfigs, "\n")
}
