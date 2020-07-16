package middleware

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetHystrixBreaker() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Request.Method + "-" + c.Request.RequestURI
		hystrix.Do(name, func() error {

			c.Next()

			if c.Writer.Status() >= http.StatusInternalServerError {
				return fmt.Errorf("status code %d", c.Writer.Status())
			}
			return nil
		}, func(e error) error {
			if e == hystrix.ErrCircuitOpen {
				c.AbortWithStatusJSON(http.StatusAccepted, gin.H{
					"code": -1,
					"msg":  "请稍后重试",
				})
			}
			return e
		})
	}
}
