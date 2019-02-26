package goweb

import (
	"net/http"
)

type RouterAdapter struct {
	RequestHandler   RequestHandler
	AppServer        *AppServer
	DisableAccessLog bool
}

func (r *RouterAdapter) HandleError(err error, resp *Response, context *RequestContext) {
	if r.AppServer.ErrorHandlerFunc != nil {
		r.AppServer.ErrorHandlerFunc(err, resp, context)
	}
}

func (r *RouterAdapter) HandleLog(context *RequestContext) {
	if r.DisableAccessLog {
		return
	}

	if r.AppServer.LogHandlerFunc != nil {
		r.AppServer.LogHandlerFunc(context)
	}
}

func (r *RouterAdapter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	theRequest := NewRequest(req)
	context := NewRequestContext(theRequest)
	resp := &Response{
		Writer:  w,
		Context: context,
	}

	err := theRequest.ParseParam()
	if err != nil {
		r.HandleError(err, resp, context)
	}

	if context.Finished() {
		r.HandleLog(context)
		return
	}

	err = r.RequestHandler.HandleRequest(theRequest, resp, context)
	if err != nil {
		r.HandleError(err, resp, context)
	}

	r.HandleLog(context)
}
