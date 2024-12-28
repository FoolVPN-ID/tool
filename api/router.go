package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/FoolVPN-ID/tool/modules/regioncheck"
	"github.com/FoolVPN-ID/tool/modules/subconverter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sagernet/sing-box/option"
)

func buildServer() *http.Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Middlewares
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Hello from gin")
	})
	r.GET("/api/v1/regioncheck", func(ctx *gin.Context) {
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
	})
	r.POST("/api/v1/convert", func(ctx *gin.Context) {
		type convertAPIFormStruct struct {
			URL      string `json:"url"`
			Format   string `json:"format"`
			Template string `json:"template"` // Only cf for now
		}
		apiForm := convertAPIFormStruct{
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
			err := subconv.ToClash()
			if err != nil {
				ctx.String(500, err.Error())
				return
			}

			ctx.YAML(200, subconv.Result.Clash)
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

			ctx.JSON(200, result)
		}
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: r.Handler(),
	}
}

func RunWithContext(ctx context.Context) {
	srv := buildServer()
	go (func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	})()

	<-ctx.Done()
	fmt.Println("API is shutting down...")

	localCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(localCtx)
}
