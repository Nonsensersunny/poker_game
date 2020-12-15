package util

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func ExtractIntsFromQuery(ctx *gin.Context, key string) ([]int, error) {
	intsS := ctx.Query(key)
	intsArr := strings.Split(intsS, ",")

	var result []int
	for _, v := range intsArr {
		if i, err := strconv.Atoi(v); err != nil {
			return result, err
		} else {
			result = append(result, i)
		}
	}

	return result, nil
}
