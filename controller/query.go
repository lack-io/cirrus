package controller

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DefaultQueryInt64(ctx *gin.Context, key string, de int64) int64 {
	page := ctx.DefaultQuery(key, fmt.Sprintf("%d", de))
	n, _ := strconv.ParseInt(page, 10, 64)
	return n
}

