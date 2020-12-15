package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ResponseWithData(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": Success,
		"data": data,
		"msg":  "",
	})
}

func ResponseWithErr(ctx *gin.Context, code int, err error) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  err.Error(),
		"data": nil,
	})
}
