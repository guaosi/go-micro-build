package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"math/rand"
	"net/http"
	"time"
)

const contextTracerKey = "Tracer-context"

// sf sampling frequency
var sf = 100

func init() {
	rand.Seed(time.Now().Unix())
}

// SetSamplingFrequency 设置采样频率
// 0 <= n <= 100
func SetSamplingFrequency(n int) {
	sf = n
}
func SetTracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		sp := opentracing.GlobalTracer().StartSpan(c.Request.URL.Path)
		tracer := opentracing.GlobalTracer()
		md := make(map[string]string)
		nsf := sf
		spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err == nil {
			sp = opentracing.GlobalTracer().StartSpan(c.Request.URL.Path, opentracing.ChildOf(spanCtx))
			tracer = sp.Tracer()
			nsf = 100
		}
		defer sp.Finish()

		if err := tracer.Inject(sp.Context(),
			opentracing.TextMap,
			opentracing.TextMapCarrier(md)); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, "init tracer failed")
		}

		ctx := context.TODO()
		ctx = opentracing.ContextWithSpan(ctx, sp)
		ctx = metadata.NewContext(ctx, md)
		c.Set(contextTracerKey, ctx)

		c.Next()

		statusCode := c.Writer.Status()
		ext.HTTPStatusCode.Set(sp, uint16(statusCode))
		ext.HTTPMethod.Set(sp, c.Request.Method)
		ext.HTTPUrl.Set(sp, c.Request.URL.EscapedPath())
		if statusCode >= http.StatusInternalServerError {
			ext.Error.Set(sp, true)
		} else if rand.Intn(100) > nsf {
			ext.SamplingPriority.Set(sp, 0)
		}
	}
}

func GetTracerContext(c *gin.Context) context.Context {
	return c.MustGet(contextTracerKey).(context.Context)
}
