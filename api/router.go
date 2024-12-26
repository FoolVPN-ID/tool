package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FoolVPN-ID/tool/modules/regioncheck"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func buildServer() *http.Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Middlewares
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Hello from gin")
	})
	r.GET("/regionCheck", func(ctx *gin.Context) {
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
