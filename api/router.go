package api

import (
	"github.com/FoolVPN-ID/RegionalCheck/modules/library"
	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Hello from gin")
	})
	r.GET("/regionCheck", func(ctx *gin.Context) {
		rawConfig := ctx.Query("config")
		if rawConfig == "" {
			ctx.String(404, "Config not provided!")
			return
		}

		lb := library.MakeLibrary()
		err := lb.Run(rawConfig)
		if err != nil {
			ctx.String(500, err.Error())
			return
		}

		ctx.JSON(200, lb.Result)
	})

	r.Run() // Listen on 0.0.0.0:8080
}
