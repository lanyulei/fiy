package app

import (
	"fiy/common/global"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"

	"fiy/tools"
)

// 失败数据处理
func Error(c *gin.Context, code int, err error, msg string) {
	var (
		res       Response
		errs      validator.ValidationErrors
		ok        bool
		errList   []string
		errString string
	)

	if err != nil {
		if reflect.TypeOf(err).String() == "validator.ValidationErrors" {
			if errs, ok = err.(validator.ValidationErrors); !ok {
				c.AbortWithStatusJSON(http.StatusOK, res.ReturnError(code))
				return
			}

			errList = make([]string, 0)
			for _, s := range errs.Translate(global.Trans) {
				errList = append(errList, s)
			}
			errString = strings.Join(errList, " | ")
		} else {
			errString = err.Error()
		}

		if msg != "" {
			res.Msg = fmt.Sprintf("%s, error: %s", msg, errString)
		} else {
			res.Msg = errString
		}
	} else {
		res.Msg = msg
	}

	//if err != nil && msg != "" {
	//	res.Msg = fmt.Sprintf("%s, error: %s", msg, err.Error())
	//} else if err != nil {
	//	res.Msg = err.Error()
	//} else if msg != "" {
	//	res.Msg = msg
	//}
	res.RequestId = tools.GenerateMsgIDFromContext(c)
	c.AbortWithStatusJSON(http.StatusOK, res.ReturnError(code))
}

// 通常成功数据处理
func OK(c *gin.Context, data interface{}, msg string) {
	var res Response
	res.Data = data
	if msg != "" {
		res.Msg = msg
	}
	res.RequestId = tools.GenerateMsgIDFromContext(c)
	c.AbortWithStatusJSON(http.StatusOK, res.ReturnOK())
}

// 分页数据处理
func PageOK(c *gin.Context, result interface{}, count int, pageIndex int, pageSize int, msg string) {
	var res Page
	res.List = result
	res.Count = count
	res.PageIndex = pageIndex
	res.PageSize = pageSize
	OK(c, res, msg)
}

// 兼容函数
func Custum(c *gin.Context, data gin.H) {
	c.AbortWithStatusJSON(http.StatusOK, data)
}
