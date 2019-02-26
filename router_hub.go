package goweb

import "errors"

type RequestFilterWrapper struct {
	f RequestFilterFunc
}

func (w *RequestFilterWrapper) FilterRequest(resp *Response, ctx *RequestContext) error {
	return w.f(resp, ctx)
}

func NewRequestFilter(f RequestFilterFunc) RequestFilter {
	return &RequestFilterWrapper{f: f}
}

// RouterHub a hub for a group of routers which share the same configuration
type RouterHub struct {
	BasePattern        string
	requestFilters     []RequestFilter
	PathDepthRouterMap map[int]([]*Router)
}

// NewRouterHub create a new routerhub
func NewRouterHub(basePattern string) *RouterHub {
	return &RouterHub{
		BasePattern:        basePattern,
		requestFilters:     make([]RequestFilter, 0),
		PathDepthRouterMap: make(map[int]([]*Router), 0),
	}
}

// AddRouter add router to the hub
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

// AddController add controller to the hub
func (rh *RouterHub) AddController(pattern string, ins interface{}, methodMap map[string]string) {
	c, err := WrapController(ins, methodMap)
	if err != nil {
		panic(err)
	}
	router := NewRouter(pattern, c, DefaultRouterConfig)
	rh.AddRouter(router)
}

// AddRequestFilter add new request filter to the hub
func (rh *RouterHub) AddRequestFilter(filter RequestFilter) {
	rh.requestFilters = append(rh.requestFilters, filter)
}

// HandleRequest implements the standard HandlerFunc interface
func (rh *RouterHub) HandleRequest(req *Request, resp *Response, ctx *RequestContext) error {
	// apply request filters
	if len(rh.requestFilters) > 0 {
		for _, filter := range rh.requestFilters {
			// check error in filters
			if err := filter.FilterRequest(resp, ctx); err != nil {
				return err
			}

			// if one filter has send response to the client, the whole process is done.
			if ctx.Finished() {
				return nil
			}
		}
	}

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
