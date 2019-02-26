package goweb

import (
	"errors"
)

type RouterConfig struct {
	DisableAccessLog bool
}

type RouterHub struct {
	BasePattern        string
	PathDepthRouterMap map[int]([]*Router)
}

func NewRouterHub(basePattern string) *RouterHub {
	return &RouterHub{
		BasePattern:        basePattern,
		PathDepthRouterMap: make(map[int]([]*Router), 0),
	}
}

func (rh *RouterHub) AddRouter(r *Router) {
	if r.PathConfig.PatternString() != rh.BasePattern {
		panic(errors.New("Router hub base pattern not match"))
	}
	dp := r.PathConfig.PathDepth
	list, found := rh.PathDepthRouterMap[dp]
	if !found {
		list = make([]*Router, 0)
	}
	list = append(list, r)
	rh.PathDepthRouterMap[dp] = list
}

func (rh *RouterHub) HandleRequest(req *Request, resp *Response, ctx *RequestContext) error {
	path := req.URL.Path
	pd := PathDepth(path)

	list, found := rh.PathDepthRouterMap[pd]
	if found {
		for _, router := range list {
			match, params := router.PathConfig.Match(path)
			if match {
				req.pathParam = params
				router.HandleRequest(req, resp, ctx)
				return nil
			}
		}
	}

	// No Match
	resp.NotFound()
	return nil
}

type Router struct {
	HandlerFunc RequestHandlerFunc
	PathConfig  *PathConfig
	Config      *RouterConfig
	AppServer   *AppServer
}

func (r *Router) HandleRequest(req *Request, resp *Response, ctx *RequestContext) error {
	return r.HandlerFunc(req, resp, ctx)
}
