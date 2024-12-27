package subconverter

import (
	"errors"
	"log"
	"strings"

	"github.com/LalatinaHub/LatinaSub-go/provider"
	"github.com/sagernet/sing-box/option"
)

type subconverterStruct struct {
	Outbounds  []option.Outbound
	rawConfigs string
	Result     struct {
		Clash string
		Raw   []string
		SFA   option.Options
		BFR   option.Options
	}
}

func MakeSubconverterFromConfig(config string) (subconverterStruct, error) {
	subconv := subconverterStruct{}
	subconv.rawConfigs = strings.ReplaceAll(config, ",", "\n")
	subconv.Outbounds = subconv.parse(subconv.rawConfigs)

	if len(subconv.Outbounds) == 0 {
		return subconv, errors.New("configs not found")
	}

	return subconv, nil
}

func (subconv *subconverterStruct) parse(content string) []option.Outbound {
	out, err := provider.Parse(content)
	if err != nil {
		log.Println(err.Error())
	}

	return out
}
