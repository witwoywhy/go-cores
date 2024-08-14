package iface

import (
	"github.com/witwoywhy/go-cores/contexts"
	"github.com/witwoywhy/go-cores/errs"
	"github.com/witwoywhy/go-cores/logger"
)

type Service[Request, Response any] interface {
	Execute(request *Request, rctx *contexts.RouteContext, l logger.Logger) (*Response, errs.Error)
}
