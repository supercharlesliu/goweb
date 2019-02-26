package goweb

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

// Response encapsulate http.ResponseWriter object to provide a simple api
// Original http.ResponseWriter can be accessed via Writer field.
type Response struct {
	Writer  http.ResponseWriter
	Context *RequestContext
}

// WriteHeader write http response code
func (resp *Response) WriteHeader(statusCode int) {
	resp.Context.StatusCode = statusCode
	resp.Writer.WriteHeader(statusCode)
}

// Write write http body data
func (resp *Response) Write(data []byte) (int, error) {
	if resp.Context.StatusCode == 0 {
		resp.WriteHeader(http.StatusOK)
	}
	return resp.Writer.Write(data)
}

// WriteString write a string as http response body
func (resp *Response) WriteString(data string) error {
	if resp.Context.StatusCode == 0 {
		resp.WriteHeader(http.StatusOK)
	}
	_, err := resp.Write([]byte(data))
	return err
}

// WriteJSON write a JSON string as response body
// Note that param data will be encoded as JSON string automatically
func (resp *Response) WriteJSON(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if resp.Context.StatusCode == 0 {
		resp.WriteHeader(http.StatusOK)
	}
	_, err = resp.Write(jsonData)
	return err
}

// WriteXML write a XML string as response body
// param data will be encoded as xml string automatically
func (resp *Response) WriteXML(data interface{}) error {
	resp.Header().Set("Content-Type", "text/xml; charset=utf-8")
	enc := xml.NewEncoder(resp)
	err := enc.Encode(data)
	return err
}

// Header get HTTP header
// This function is usefull if you want to set response header by yourself.
func (resp *Response) Header() http.Header {
	return resp.Writer.Header()
}

// NotFound shorthand for return 404
func (resp *Response) NotFound() {
	http.NotFound(resp.Writer, resp.Context.Request.Req)
}

// Redirect redirect to a given url
func (resp *Response) Redirect(url string, code int) {
	http.Redirect(resp.Writer, resp.Context.Request.Req, url, code)
}
