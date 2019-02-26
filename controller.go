package goweb

import (
	"reflect"

	"github.com/pkg/errors"
)

type JSONResponse struct {
	StatusCode int
	BodyData   interface{}
}

func (r *JSONResponse) Response(resp *Response) {
	resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp.WriteHeader(r.StatusCode)
	resp.WriteJSON(r.BodyData)
}

type Controller struct {
	MethodFuncMap map[string]*reflect.Value
}

// Invoke invoke controller method
func (cm *Controller) Invoke(req *Request, resp *Response, context *RequestContext) error {
	methodRef, found := cm.MethodFuncMap[req.Req.Method]
	if !found {
		resp.WriteHeader(405)
		return nil
	}

	inValues := make([]reflect.Value, 0)
	inValues = append(inValues, reflect.ValueOf(req))
	inValues = append(inValues, reflect.ValueOf(resp))
	inValues = append(inValues, reflect.ValueOf(context))

	outValues := (*methodRef).Call(inValues)
	if outValues == nil || len(outValues) < 1 {
		return nil
	}
	if !outValues[0].IsValid() {
		return nil
	}

	first := outValues[0].Interface()
	if j, ok := first.(*JSONResponse); ok {
		j.Response(resp)
		return nil
	}
	// TODO: check return type
	return nil
}

func WrapController(ins interface{}, methodMap map[string]string) (RequestHandlerFunc, error) {
	v := reflect.ValueOf(ins)
	methodFuncMap := make(map[string]*reflect.Value, 0)

	for httpMethod, methodName := range methodMap {
		controllerMethod := v.MethodByName(methodName)
		if !controllerMethod.IsValid() {
			return nil, errors.WithMessage(ErrMethodNotFound, methodName)
		}
		methodFuncMap[httpMethod] = &controllerMethod
	}

	mapper := &Controller{
		MethodFuncMap: methodFuncMap,
	}

	return mapper.Invoke, nil
}
