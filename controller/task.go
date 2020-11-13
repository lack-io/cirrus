package controller

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lack-io/cirrus/internal/daemon"
)

type TaskState string

const (
	Free    TaskState = "free"
	Running TaskState = "running"
	Pending TaskState = "pending"
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

func RegistryTaskController(d daemon.Daemon, handler *gin.Engine) {
	task := &Task{State: Free}
	controller := taskController{task: task, lock: &sync.RWMutex{}, d: d}
	group := handler.Group("/api/v1/task/")
	{
		group.GET("", controller.getTask())
		group.POST("action/start", controller.startTask())
		group.POST("action/pause", controller.pauseTask())
	}
}

type taskController struct {
	task *Task

	lock *sync.RWMutex

	d daemon.Daemon
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
		defer c.lock.Unlock()

		if c.task.State == Running {
			R().Ctx(ctx).Bad(fmt.Errorf("task is running"))
			return
		}

		type data struct {
			Root string `json:"root,omitempty"`
		}

		d := data{}
		ctx.BindJSON(&d)
		if d.Root == "" {
			R().Ctx(ctx).Bad(fmt.Errorf("缺少 root 参数"))
			return
		}

		c.task.StartTime = time.Now().Unix()
		c.task.State = Running
		c.d.StartDaemon(d.Root)

		R().Ctx(ctx).Accepted()
		return
	}
}

func (c *taskController) pauseTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c.lock.Lock()
		defer c.lock.Unlock()

		if c.task.State != Running {
			R().Ctx(ctx).Bad(fmt.Errorf("task is not running"))
			return
		}

		c.task.EndTime = time.Now().Unix()
		c.task.State = Pending
		c.d.PauseDaemon()

		R().Ctx(ctx).Accepted()
		return
	}
}
