package api

import (
	"strings"

	"github.com/FoolVPN-ID/tool/modules/subconverter"
	"github.com/gin-gonic/gin"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json"
)

type ConvertAPIFormStruct struct {
	URL      string `json:"url"`
	Format   string `json:"format"`
	Template string `json:"template"` // Only cf for now
}

func HandlePostConvert(ctx *gin.Context) {
	apiForm := ConvertAPIFormStruct{
		Format: "raw",
	}
	if err := ctx.ShouldBindBodyWithJSON(&apiForm); err != nil {
		ctx.String(400, err.Error())
		return
	}

	subconv, err := subconverter.MakeSubconverterFromConfig(apiForm.URL)
	if err != nil {
		ctx.String(500, err.Error())
		return
	}

	switch apiForm.Format {
	case "raw":
		subconv.ToRaw()
		ctx.String(200, strings.Join(subconv.Result.Raw, "\n"))
	case "clash":
		var (
			result map[string]any
			err    = subconv.ToClash()
		)
		if err != nil {
			ctx.String(500, err.Error())
			return
		}

		result = subconv.Result.Clash
		if apiForm.Template == "cf" {
			result = subconv.PostTemplateClash(apiForm.Template, result)
		}

		ctx.YAML(200, result)
	case "bfr", "sfa":
		var (
			result option.Options
			err    error
		)

		if apiForm.Format == "bfr" {
			err = subconv.ToBFR()
			result = subconv.Result.BFR
		} else {
			err = subconv.ToSFA()
			result = subconv.Result.SFA
		}

		if err != nil {
			ctx.String(500, err.Error())
			return
		}

		if apiForm.Template == "cf" {
			result = subconv.PostTemplateSingBox(apiForm.Template, result)
		}

		var (
			bufResult, _ = json.Marshal(result)
			mapResult    any
		)

		json.Unmarshal(bufResult, &mapResult)
		ctx.JSON(200, mapResult)
	}
}
