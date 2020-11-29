package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/lack-io/cirrus/store"
)

func RegistryGoodController(store *store.Store, handler *gin.RouterGroup) {
	controller := goodController{store: store}
	group := handler.Group("/v1/goods")
	{
		group.GET("", controller.getGoods())
		group.DELETE("/:id", controller.delGroup())
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
		pagination := &store.Pagination{Page: int(page), Size: int(size)}
		if start > 0 || end > 0 {
			goods, err = c.store.GetGoodsByTimeout(start, end, pagination)
		} else {
			goods, err = c.store.GetGoods(pagination)
		}
		if err != nil {
			R().Ctx(ctx).Fail(err)
			return
		}

		R().Ctx(ctx).OK(gin.H{
			"list":       goods,
			"pagination": pagination,
		})
		return
	}
}

func (c *goodController) delGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ids := ctx.Param("id")
		id, _ := strconv.ParseInt(ids, 10, 64)

		good, err := c.store.DelGroup(id)
		if err != nil {
			R().Ctx(ctx).Fail(err)
			return
		}

		R().Ctx(ctx).OK(good)
		return
	}
}
