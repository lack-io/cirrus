package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/lack-io/cirrus/proxy"
)

func RegistryProxyController(p proxy.Proxy, handler *gin.RouterGroup) {
	controller := proxyController{p: p}
	group := handler.Group("/v1/proxy")
	{
		group.GET("", controller.getProxy())
	}
}

type proxyController struct {
	p proxy.Proxy
}

func (c *proxyController) getProxy() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		bl, _ := c.p.GetBalance(ctx)

		R().Ctx(ctx).OK(gin.H{
			"proxy":   c.p.GetJSON(),
			"balance": bl,
		})
		return
	}
}
