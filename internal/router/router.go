package router

import (
	"go_blog/pkg/config"
	"go_blog/pkg/transports/httpd/middleware"
	"go_blog/pkg/utils"

	"gopkg.in/antage/eventsource.v1"

	"github.com/gin-gonic/gin"
)

type Router struct {
}

// Router create router handler
func NewRouter(
	conf *config.Config,
	response *utils.Response,
	router *gin.Engine,
	sse eventsource.EventSource,

) {
	// cross-domain middleware
	router.Use(middleware.Cors())
	// Internationalization middleware
	router.Use(middleware.I18nMiddleware(conf))
	// 404 middleware
	router.NoRoute(middleware.NotFoundHandler(response))
	r := &Router{}
	group := router.Group("/api")
	// r.commonRoute(conf, response, group, sse)
	// r.frontRoute(conf, response, group)
	// r.adminRoute(conf, response, group)
}
