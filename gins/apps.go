package gins

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/witwoywhy/go-cores/logs"
)

type GinApps interface {
	Register(method string, relativePath string, handlers ...gin.HandlerFunc)
	UseMiddleware(middleware ...gin.HandlerFunc)
	ListenAndServe(addr string)
}

type app struct {
	gin *gin.Engine
}

func New() GinApps {
	return &app{
		gin: gin.New(),
	}
}

func (a *app) ListenAndServe(addr string) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := http.ListenAndServe(addr, a.gin); err != nil {
			logs.L.Error(err)
		}
	}()
	wg.Wait()
}

func (a *app) Register(method string, relativePath string, handlers ...gin.HandlerFunc) {
	switch method {
	case http.MethodGet:
		a.gin.GET(relativePath, handlers...)
	case http.MethodPost:
		a.gin.POST(relativePath, handlers...)
	case http.MethodPatch:
		a.gin.PATCH(relativePath, handlers...)
	case http.MethodPut:
		a.gin.PUT(relativePath, handlers...)
	case http.MethodDelete:
		a.gin.DELETE(relativePath, handlers...)
	}
}

func (a *app) UseMiddleware(middleware ...gin.HandlerFunc) {
	a.gin.Use(middleware...)
}
