package main

import (
	"github.com/gin-gonic/gin"
	gxxl "github.com/go-xxl/gin-xxl"
	xxl "github.com/go-xxl/xxl"
	"github.com/go-xxl/xxl/admin"
	"github.com/go-xxl/xxl/job"
	"github.com/go-xxl/xxl/server"
	"github.com/go-xxl/xxl/utils"
	"time"
)

func main() {
	g := gin.New()

	g.Use(func(context *gin.Context) {
		context.Set("traceId", utils.Uuid())
	})

	router := gxxl.NewGinRouter(xxl.Options{
		AdmAddresses:  []string{"http://127.0.0.1:8080/xxl-job-admin/"},
		Timeout:       10 * time.Second,
		ExecutorIp:    utils.GetLocalIp(),
		ExecutorPort:  "8081",
		RegistryKey:   "rule-engine",
		RegistryValue: utils.BuildEndPoint(utils.GetLocalIp(), "8081"),
	})

	router.Job("/demo", func(ctx *server.Context) job.Resp {
		param := ctx.Param

		time.Sleep(time.Second * 30)

		return job.Resp{
			LogID:       param.LogID,
			LogDateTime: time.Now().Unix(),
			HandleCode:  admin.SuccessCode,
			HandleMsg:   "get result",
		}
	})
	router.Router(&g.RouterGroup)
	router.Register()

	g.Run(":8081")
}
