package xxl

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-xxl/xxl/log"
	"github.com/go-xxl/xxl/server"
	"go.uber.org/zap"
	"io/ioutil"
)

// convertGinFunc convert func to gin func
func convertGinFunc(f func(ctx *server.Context)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := &server.RunReq{}

		requestBody, err := ctx.GetRawData()
		if err != nil {
			log.Warn("get gin request body data err",
				zap.String("err", err.Error()))
		}
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

		if err = json.Unmarshal(requestBody, &params); err != nil {
			log.Warn("get gin request body data err",
				zap.String("err", err.Error()))
		}

		c := &server.Context{
			Writer:  ctx.Writer,
			Request: ctx.Request,
			Param:   params,
			TraceId: ctx.GetString("traceId"),
		}

		f(c)
	}
}
