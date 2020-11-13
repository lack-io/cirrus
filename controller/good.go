package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/lack-io/cirrus/store"
)

func RegistryGoodController(store *store.Store, handler *gin.Engine) {
	controller := goodController{store: store}
	group := handler.Group("/api/v1/goods/")
	{
		group.GET("", controller.getGoods())
	}
}

type goodController struct {
	store *store.Store
}

func (c *goodController) getGoods() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		page := DefaultQueryInt64(ctx, "page", 1)
		size := DefaultQueryInt64(ctx, "size", 10)
		start := DefaultQueryInt64(ctx, "start", 0)
		end := DefaultQueryInt64(ctx, "end", 0)

		var err error
		var goods []*store.Good
		if start > 0 || end > 0 {
			goods, err = c.store.GetGoodsByTimeout(start, end, &store.Pagination{Page: int(page), Size: int(size)})
		} else {
			goods, err = c.store.GetGoods(&store.Pagination{Page: int(page), Size: int(size)})
		}
		if err != nil {
			R().Ctx(ctx).Fail(err)
			return
		}

		R().Ctx(ctx).OK(goods)
		return
	}
}
