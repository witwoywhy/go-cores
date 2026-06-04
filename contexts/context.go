package contexts

import "context"

type Context interface {
	GetHeader() Header
	GetContext() context.Context
}

