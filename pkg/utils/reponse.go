package utils

import (
	"net/http"

	"go_blog/pkg/config"
	"go_blog/pkg/errorsf"

	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Response struct {
	conf *config.Config
}

func NewResponse(conf *config.Config) *Response {
	return &Response{conf: conf}
}

func (o Response) Success(c *gin.Context, data interface{}) {
	err := errorsf.SUCCESS
	msg := o.GetI18nMessage(c, err)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": err.GetMessageId(),
		"text":    msg,
		"data":    data,
	})
}

func (o Response) Error(c *gin.Context, err errorsf.ErrInter) {
	msg := o.GetI18nMessage(c, err)
	resp := gin.H{
		"code":    1,
		"message": err.GetMessageId(),
		"text":    msg,
	}
	if o.conf.Server.Mode != "prod" {
		resp["error"] = err.Error()
	}
	c.JSON(http.StatusOK, resp)
}

func (Response) GetI18nMessage(c *gin.Context, err errorsf.ErrInter) string {
	msg := ginI18n.MustGetMessage(
		c,
		&i18n.LocalizeConfig{
			MessageID:    err.GetMessageId(),
			TemplateData: err.GetMessageParams(),
		})
	if msg == "" {
		msg = err.GetMessageId()
	}
	return msg
}
