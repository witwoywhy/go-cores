package contexts

import "context"

type RouteContext struct {
	Header
	Ctx context.Context
}

func (c *RouteContext) GetHeader() Header {
	return c.Header
}

func (c *RouteContext) GetContext() context.Context {
	return c.Ctx
}
