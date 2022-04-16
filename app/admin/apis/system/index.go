package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
  @Desc :
*/

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
