package gins

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/witwoywhy/go-cores/errs"
)

func Error() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		for _, err := range ctx.Errors {
			if e, ok := err.Err.(errs.Error); !ok {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, errs.NewInternalError())
			} else {
				ctx.AbortWithStatusJSON(e.HttpStatus(), e)
			}
		}
	}
}
