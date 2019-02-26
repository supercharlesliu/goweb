package goweb

// RouterConfig config for a router
type RouterConfig struct {
	DisableAccessLog bool
}

var DefaultRouterConfig *RouterConfig

func init() {
	DefaultRouterConfig = &RouterConfig{DisableAccessLog: false}
}

// Router represent a router rule
type Router struct {
	HandlerFunc RequestHandlerFunc
	PathConfig  *PathConfig
	Config      *RouterConfig
}

// HandleRequest implements the standard HandlerFunc interface
func (r *Router) HandleRequest(req *Request, resp *Response, ctx *RequestContext) error {
	return r.HandlerFunc(req, resp, ctx)
}

// NewRouter create a router object
func NewRouter(pattern string, handlerFunc RequestHandlerFunc, config *RouterConfig) *Router {
	pathConfig := ParsePathParam(pattern)
	return &Router{
		HandlerFunc: handlerFunc,
		PathConfig:  pathConfig,
		Config:      config,
	}
}
