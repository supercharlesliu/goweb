package goweb

import (
	"context"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

var (
	methodsWithBody = map[string]bool{"POST": true, "PUT": true, "PATCH": true}
)

// Request wraps net/http.Request to provide form parsing ability
type Request struct {
	Req       *http.Request
	URL       *url.URL
	pathParam map[string]string
}

// NewRequest create a request obect from net/http.Request
func NewRequest(req *http.Request) *Request {
	pathParam := make(map[string]string, 0)
	// Get params
	ctx := req.Context()
	params := ctx.Value("params")
	if params != nil {
		if v, ok := params.(map[string]string); ok {
			pathParam = v
		}
	}

	return &Request{
		Req:       req,
		URL:       req.URL,
		pathParam: pathParam,
	}
}

// ParseParam parse http form
func (r *Request) ParseParam() error {
	hasBody, _ := methodsWithBody[r.Req.Method]
	if !hasBody {
		return nil
	}
	contentType := r.Header("Content-Type")
	if contentType == "" {
		return ErrUnknowContentType
	}
	parts := strings.Split(contentType, ";")
	parts[0] = strings.Trim(parts[0], " ")

	if parts[0] == "application/x-www-form-urlencoded" {
		return r.Req.ParseForm()
	} else if parts[0] == "multipart/form-data" {
		return r.Req.ParseMultipartForm(5 * (1 << 20))
	}

	return r.Req.ParseForm()
}

// PathParam get Path Param
func (r *Request) PathParam(name string) string {
	return r.pathParam[name]
}

// Header get request header
func (r *Request) Header(name string) string {
	return r.Req.Header.Get(name)
}

// Cookie get request cookie
func (r *Request) Cookie(name string) string {
	cookie, err := r.Req.Cookie(name)
	if err == http.ErrNoCookie {
		return ""
	}
	return cookie.Value
}

// Param get a param from QueryString or Body
func (r *Request) Param(name string) string {
	return r.Req.FormValue(name)
}

// FillForm populate a form for data validating
func (r *Request) FillForm(m interface{}) error {
	formMeta, err := InspectForm(m)
	if err != nil {
		return err
	}

	form := reflect.Indirect(reflect.New(*formMeta.Type))
	for fieldName, fieldConf := range formMeta.FieldMap {
		var v string
		if fieldConf.FromPath {
			v = r.PathParam(fieldConf.ParamName)
		} else if fieldConf.FromHeader {
			v = r.Header(fieldConf.ParamName)
		} else if fieldConf.FromCookie {
			v = r.Cookie(fieldConf.ParamName)
		} else {
			v = r.Param(fieldConf.ParamName)
		}

		if fieldConf.IsRequired && v == "" {
			return newFormError(FormErrMissingRequired, fieldConf.ParamName, nil)
		}

		value, err := StringCast(v, fieldConf.Type)
		if err != nil {
			if !fieldConf.IsRequired {
				continue
			}
			return newFormError(FormErrTypeCannotCast, fieldConf.ParamName, err)
		}

		fieldV := form.FieldByName(fieldName)
		vv := reflect.ValueOf(value)
		if fieldV.CanSet() && vv.IsValid() {
			fieldV.Set(vv)
		}
	}

	mV := reflect.Indirect(reflect.ValueOf(m))
	mV.Set(form)

	// Validate required
	if vd, ok := m.(FormValidator); ok {
		err = vd.ValidateForm()
		if err != nil {
			return newFormError(FormErrNotAValidForm, "", err)
		}
	}

	return nil
}

// SetContextValue set context value of the current request so that it can be passed down
func (r *Request) SetContextValue(key interface{}, value interface{}) {
	req := r.Req
	ctx := context.WithValue(req.Context(), key, value)
	reqNew := req.WithContext(ctx)
	r.Req = reqNew
}
