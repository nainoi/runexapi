package tracer

import (
	//"fmt"

	"fmt"

	"github.com/gin-gonic/gin"
	stdopentracing "github.com/opentracing/opentracing-go"
	Tracelog "github.com/opentracing/opentracing-go/log"
	"thinkdev.app/think/runex/runexapi/logger"
)

func CreateTracerAndSpan(spanName string, c *gin.Context) (stdopentracing.Span, error) {

	tracer := stdopentracing.GlobalTracer()

	userSpanCtx, err := tracer.Extract(stdopentracing.HTTPHeaders, stdopentracing.HTTPHeadersCarrier(c.Request.Header))

	if err != nil {
		logger.Logger.Infof(err.Error())
		return nil, err
	}

	userSpan := tracer.StartSpan(spanName, stdopentracing.ChildOf(userSpanCtx))
	defer userSpan.Finish()
	return userSpan, nil
}

func OnErrorLog(receivedSpan stdopentracing.Span, err error) {
	fmt.Println(receivedSpan)
	receivedSpan.LogFields(
		Tracelog.String("event", "error"),
		Tracelog.String("message", err.Error()),
	)
	logger.Logger.Infof(err.Error())
}