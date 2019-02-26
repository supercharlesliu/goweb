package goweb

import (
	"strconv"
	"strings"
	"time"
)

func fmtDuration(duration time.Duration) string {
	return strconv.FormatFloat(duration.Seconds(), 'f', 6, 64)
}

// RequestContext context for a request
type RequestContext struct {
	StartTime  time.Time
	Method     string
	Proto      string
	Host       string
	URI        string
	RemoteAddr string
	UserAgent  string
	StatusCode int
	Request    *Request
}

// NewRequestContext create request context from a given net/http.Request
func NewRequestContext(req *Request) *RequestContext {
	r := req.Req
	return &RequestContext{
		StartTime:  time.Now(),
		Method:     r.Method,
		Proto:      r.Proto,
		Host:       r.Host,
		URI:        r.RequestURI,
		RemoteAddr: r.RemoteAddr,
		UserAgent:  r.UserAgent(),
		Request:    req,
	}
}

// Finished test if request has finished processing
func (c *RequestContext) Finished() bool {
	return c.StatusCode != 0
}

func (c *RequestContext) String() string {
	fields := []string{strconv.Itoa(c.StatusCode), c.Proto, c.Host, c.Method, c.URI, fmtDuration(time.Since(c.StartTime)), c.RemoteAddr, c.UserAgent}
	return "[" + strings.Join(fields, "] - [") + "]"
}
