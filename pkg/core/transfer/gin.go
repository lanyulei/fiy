package transfer

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/**
e.g.

package main

import (
	"github.com/gin-gonic/gin"
	"fiy/pkg/core/transfer"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	r := gin.Default()
	r.GET("/monitor", transfer.Handler(promhttp.Handler()))
	r.Run(":9999")
}
**/

// Handler http.Handler 转换成 gin.HandlerFunc
func Handler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
