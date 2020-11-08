package controller

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/xingyys/cirrus/cdiscount"
)

type TaskState string

const (
	Free    TaskState = "free"
	Running TaskState = "running"
	//Pending TaskState = "pending"
)

// 抓取任务
type Task struct {
	Root string `json:"root,omitempty"`
	// 任务状态
	State TaskState `json:"state,omitempty"`
	// 任务开始时间
	StartTime int64 `json:"startTime,omitempty"`
	// 任务结束时间
	EndTime int64 `json:"endTime,omitempty"`
}

func RegistryTaskController(cds *cdiscount.Cdiscount, handler *gin.Engine) {
	task := &Task{State: Free}
	controller := taskController{task: task, lock: &sync.RWMutex{}, cds: cds}
	group := handler.Group("/api/v1/task/")
	{
		group.GET("", controller.getTask())
		group.POST("action/start", controller.startTask())
	}
}

type taskController struct {
	task *Task

	lock *sync.RWMutex

	cds *cdiscount.Cdiscount
}

func (c *taskController) getTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c.lock.RLock()
		defer c.lock.RUnlock()

		R().Ctx(ctx).OK(c.task)
		return
	}
}

func (c *taskController) startTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c.lock.Lock()
		defer c.lock.RUnlock()

		if c.task.State == Running {
			R().Ctx(ctx).Bad(fmt.Errorf("task is running"))
			return
		}

		data := &Task{}
		if err := ctx.ShouldBind(data); err != nil {
			R().Ctx(ctx).Bad(err)
			return
		}

		c.task.StartTime = time.Now().Unix()
		c.task.State = Running
		// TODO: start task


		R().Ctx(ctx).Accepted()
		return
	}
}
