package transports

import (
	"go_blog/pkg/transports/httpd"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		httpd.NewServer,
	),
	fx.Invoke(
		func(r *gin.Engine) {},
	),
)
