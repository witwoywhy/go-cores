package echos

import (
	"net/http"

	"github.com/aymerick/raymond"
	"github.com/labstack/echo/v5"
	"github.com/witwoywhy/go-cores/errs"
)

func Error(mapping errs.ErrorCodeMapping) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			err := next(ctx)
			if err == nil {
				return nil
			}

			ers, ok := err.(*errs.Errs)
			if !ok {
				return ctx.JSON(http.StatusInternalServerError, errs.NewInternalError(err))
			}

			status := http.StatusInternalServerError
			for i, e := range ers.Errors {
				status = e.HttpStatus()

				mappingLang, ok := mapping[e.ErrorCode]
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

			return ctx.JSON(status, ers)
		}
	}
}
