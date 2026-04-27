package gins

import (
	"net/http"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
	"github.com/witwoywhy/go-cores/errs"
)

func (a *app) Error() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		var status int
		for _, err := range ctx.Errors {
			ers, ok := err.Err.(*errs.Errs)
			if !ok {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, errs.NewInternalError())
				return
			}

			for i, e := range ers.Errors {
				status = e.HttpStatus()

				mappingLang, ok := a.errorMapping[e.ErrorCode]
				if !ok {
					continue
				}

				lang := GetLanguage(ctx)
				mapping, ok := mappingLang[lang]
				if !ok {
					continue
				}

				if e.Data != nil {
					message, err := raymond.Render(mapping.Message, e.Data)
					if err == nil {
						ers.Errors[i].Message = message
						continue
					}
				}

				ers.Errors[i].Message = mapping.Message
			}

			ctx.AbortWithStatusJSON(status, ers)
		}
	}
}
