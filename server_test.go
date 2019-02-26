package goweb

import (
	"log"
	"strconv"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	port := 8081

	appConfig := &AppServerConfig{
		Addr:           ":" + strconv.Itoa(port),
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
		LogHandlerFunc: func(context *RequestContext) {
			log.Println(context)
		},
		ErrorHandlerFunc: func(err error, resp *Response, context *RequestContext) {
			resp.WriteHeader(500)
		},
	}
	server := NewAppServer(appConfig)

	// for health check
	server.ResponseText("/ping", "pong")

	server.AddRouter("/users/", func(req *Request, resp *Response, context *RequestContext) error {
		resp.WriteString("Haha")
		return nil
	}, nil)

	server.AddRouter("/users/:userId", func(req *Request, resp *Response, context *RequestContext) error {
		resp.WriteString(req.PathParam("userId"))
		return nil
	}, nil)

	server.AddRouter("/users/:userId/sites/:siteId", func(req *Request, resp *Response, context *RequestContext) error {
		resp.WriteString(req.PathParam("userId"))
		resp.WriteString(req.PathParam("siteId"))
		return nil
	}, nil)

	hub := NewRouterHub("/hubs/")
	hub.AddRequestFilter(NewRequestFilter(func(resp *Response, ctx *RequestContext) error {
		return nil
	}))
	hub.AddRouter(NewRouter("/hubs/:id", func(req *Request, resp *Response, ctx *RequestContext) error {
		return nil
	}, DefaultRouterConfig))
	server.AddHub(hub)

	server.Start()
	// m.Run()
}
