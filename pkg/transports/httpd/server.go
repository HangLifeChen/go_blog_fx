package httpd

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go_blog/pkg/config"
	"go_blog/pkg/validate"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// newServer creates and configures a new HTTP server using Gin.
// It also sets up lifecycle hooks for starting and stopping the server.
func NewServer(lc fx.Lifecycle, cfg *config.Config) *gin.Engine {
	g := gin.New()
	g.Use(gin.Recovery())

	ginMode(cfg)

	// bind custom validator
	validate.BindValidator()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      g,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Append hooks to the lifecycle for starting and stopping the server.
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Printf("Start the server :%d\n", cfg.Server.Port)
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					fmt.Printf("failed to close http server: %v\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Stopped the server")
			return srv.Shutdown(ctx)
		},
	})
	return g
}

func ginMode(conf *config.Config) {
	// gin mode setting
	switch conf.Server.Mode {
	case "local", "dev":
		gin.SetMode(gin.DebugMode)
	case "pre":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
}
