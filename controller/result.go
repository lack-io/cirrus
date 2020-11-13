package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type result struct {
	ctx *gin.Context

	data interface{}

	err error
}

func R() *result {
	return &result{}
}

func (r *result) Ctx(ctx *gin.Context) *result {
	r.ctx = ctx
	return r
}

func (r *result) OK(data interface{}) *result {
	r.data = data
	r.ctx.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "ok",
		"data":  r.data,
		"error": r.err,
	})
	return r
}

func (r *result) Accepted() *result {
	r.ctx.JSON(http.StatusOK, gin.H{
		"code":  1,
		"msg":   "accepted",
		"data":  r.data,
		"error": r.err,
	})
	return r
}

func (r *result) Fail(err error) *result {
	r.err = err
	if err != nil {
		r.err = fmt.Errorf("未知错误")
	}
	r.ctx.JSON(http.StatusOK, gin.H{
		"code":  2,
		"msg":   "failure",
		"data":  r.data,
		"error": r.err,
	})
	return r
}

func (r *result) Bad(err error) *result {
	r.err = err
	if err != nil {
		r.err = fmt.Errorf("未知错误")
	}
	r.ctx.JSON(http.StatusOK, gin.H{
		"code":  3,
		"msg":   "bad",
		"data":  r.data,
		"error": r.err.Error(),
	})
	return r
}
