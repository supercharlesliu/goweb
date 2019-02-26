package goweb

// RequestHandlerFunc Request handle func
type RequestHandlerFunc func(req *Request, resp *Response, ctx *RequestContext) error

// ErrorHandlerFunc Error handle func
type ErrorHandlerFunc func(err error, resp *Response, ctx *RequestContext)

// LogHandlerFunc Log handle func
type LogHandlerFunc func(ctx *RequestContext)

// FormValidator An interface indecate that a form can be validated be user defined function
type FormValidator interface {
	ValidateForm() error
}

// RequestFilterFunc type of request filter
type RequestFilterFunc func(resp *Response, ctx *RequestContext) error

// RequestFilter a filter to pre-process request before it goes to RequestHandler
type RequestFilter interface {
	FilterRequest(resp *Response, ctx *RequestContext) error
}

// RequestHandler handle request and generate response
type RequestHandler interface {
	HandleRequest(req *Request, resp *Response, ctx *RequestContext) error
}

// LogHandler output log for a given request context
type LogHandler interface {
	HandleLog(ctx *RequestContext)
}
