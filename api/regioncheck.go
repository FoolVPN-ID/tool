package api

import (
	"github.com/FoolVPN-ID/tool/modules/regioncheck"
	"github.com/gin-gonic/gin"
)

func HandleGetRegionCheck(ctx *gin.Context) {
	rawConfig := ctx.Query("config")
	if rawConfig == "" {
		ctx.String(404, "Config not provided!")
		return
	}

	rc := regioncheck.MakeLibrary()
	err := rc.Run(rawConfig)
	if err != nil {
		ctx.String(500, err.Error())
		return
	}

	ctx.JSON(200, rc.Result)
}
