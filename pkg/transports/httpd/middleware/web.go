package middleware

import (
	"net/http"

	"go_blog/pkg/config"
	"go_blog/pkg/errorsf"
	"go_blog/pkg/utils"

	"github.com/gin-gonic/gin"
)

// NotFoundHandler not found api router
func NotFoundHandler(response *utils.Response) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response.Error(ctx, errorsf.NOT_FOUND)
	}
}

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, X-Requested-With, XMLHttpRequest,Unique-Finger")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type,Unique-Finger")
		context.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusOK)
		}
		context.Next()
	}
}

// func JwtMiddleware(conf *config.Config, response *utils.Response) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		authHeader := ctx.GetHeader("Authorization")
// 		if len(authHeader) == 0 {
// 			response.Error(ctx, errorsf.TOKEN_INVALID)
// 			ctx.Abort()
// 			return
// 		}
// 		auths := strings.Split(authHeader, " ")

// 		bearer := auths[0]

// 		if len(auths) != 2 || bearer != "Bearer" {
// 			response.Error(ctx, errorsf.TOKEN_INVALID)
// 			ctx.Abort()
// 			return
// 		}
// 		token := auths[1]
// 		user, err := utils.ParseJwtToken(token)

// 		if err != nil {
// 			response.Error(ctx, errorsf.TOKEN_INVALID)
// 			ctx.Abort()
// 			return
// 		}

// 		ctx.Set("uid", user.Uid)
// 		ctx.Set("email", user.Email)
// 		ctx.Set("token", token)
// 		ctx.Next()
// 	}
// }

// func NotRequiredJwtMiddleware(conf *config.Config, response *utils.Response) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		authHeader := ctx.GetHeader("Authorization")
// 		if len(authHeader) == 0 {
// 			ctx.Next()
// 			return
// 		}
// 		auths := strings.Split(authHeader, " ")

// 		bearer := auths[0]

// 		if len(auths) != 2 || bearer != "Bearer" {
// 			ctx.Next()
// 			return
// 		}
// 		token := auths[1]
// 		user, err := utils.ParseJwtToken(token)

// 		if err != nil {
// 			ctx.Next()
// 			return
// 		}

// 		ctx.Set("uid", user.Uid)
// 		ctx.Set("email", user.Email)
// 		ctx.Next()
// 	}
// }

// func AdminJwtMiddleware(conf *config.Config, response *utils.Response) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		authHeader := ctx.GetHeader("Authorization")
// 		if len(authHeader) == 0 {
// 			response.Error(ctx, errorsf.TOKEN_INVALID)
// 			ctx.Abort()
// 			return
// 		}
// 		auths := strings.Split(authHeader, " ")

// 		bearer := auths[0]

// 		if len(auths) != 2 || bearer != "Bearer" {
// 			response.Error(ctx, errorsf.TOKEN_INVALID)
// 			ctx.Abort()
// 			return
// 		}
// 		token := auths[1]
// 		user, err := utils.ParseJwtToken(token)

// 		if err != nil {
// 			response.Error(ctx, errorsf.TOKEN_INVALID)
// 			ctx.Abort()
// 			return
// 		}
// 		if user.UserType != utils.UserTypeAdmin {
// 			response.Error(ctx, errorsf.TOKEN_INVALID)
// 			ctx.Abort()
// 			return
// 		}
// 		ctx.Set("userId", int64(user.UserId))
// 		ctx.Next()
// 	}
// }

func WhitelistMiddleware(conf *config.Config, response *utils.Response) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := "0.0.0.0"
		allow, ok := conf.Server.Whitelist[ip]
		if ok {
			if allow {
				ctx.Next()
				return
			} else {
				response.Error(ctx, errorsf.UN_AUTH)
				ctx.Abort()
				return
			}
		}
		ip = ctx.ClientIP()
		allow, ok = conf.Server.Whitelist[ip]
		if !ok || !allow {
			response.Error(ctx, errorsf.UN_AUTH)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
