package goweb

import (
	"log"
	"net/http"
	"time"
)

// AppServer wraps net/http.Server to handle logs, errors, router registration
type AppServer struct {
	Server               *http.Server
	ServeMux             *http.ServeMux
	LogHandlerFunc       LogHandlerFunc
	ErrorHandlerFunc     ErrorHandlerFunc
	basePatternRouterMap map[string]([]*Router)
}

// AppServerConfig config structure for AppServer
type AppServerConfig struct {
	Addr             string        // Address this app server listens to, eg: :80
	ReadTimeout      time.Duration // ReadTimeout for request header
	WriteTimeout     time.Duration // WriteTimeout for response
	MaxHeaderBytes   int           // Max Number of bytes of the request header
	LogHandlerFunc   LogHandlerFunc
	ErrorHandlerFunc ErrorHandlerFunc
}

// NewAppServer create a new AppServer instance
func NewAppServer(config *AppServerConfig) *AppServer {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:           config.Addr,
		Handler:        mux,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}
	appServer := &AppServer{
		Server:               server,
		ServeMux:             mux,
		LogHandlerFunc:       config.LogHandlerFunc,
		ErrorHandlerFunc:     config.ErrorHandlerFunc,
		basePatternRouterMap: make(map[string]([]*Router), 0),
	}
	return appServer
}

// ResponseText echo given text for a given url, this is useful for cases like ping/pong healthcheck
func (server *AppServer) ResponseText(url string, responseStr string) {
	server.AddRouter(url, func(req *Request, resp *Response, context *RequestContext) error {
		resp.WriteString(responseStr)
		return nil
	}, &RouterConfig{
		DisableAccessLog: true,
	})
}

// AddController register controller to this AppServer
func (server *AppServer) AddController(pattern string, ins interface{}, methodMap map[string]string) {
	c, err := WrapController(ins, methodMap)
	if err != nil {
		panic(err)
	}
	server.AddRouter(pattern, c, &RouterConfig{
		DisableAccessLog: false,
	})
}

// AddHub add router hub to the server
func (server *AppServer) AddHub(hub *RouterHub) {
	server.Handle(hub.BasePattern, hub, false)
}

func (server *AppServer) AddRouter(pattern string, handlerFunc RequestHandlerFunc, config *RouterConfig) {
	router := NewRouter(pattern, handlerFunc, config)

	list, found := server.basePatternRouterMap[router.PathConfig.PatternString()]
	if !found {
		list = make([]*Router, 0)
	}
	list = append(list, router)
	server.basePatternRouterMap[router.PathConfig.PatternString()] = list
}

func (server *AppServer) Handle(pattern string, h RequestHandler, disableAccessLog bool) {
	server.ServeMux.Handle(pattern, &RouterAdapter{
		RequestHandler:   h,
		AppServer:        server,
		DisableAccessLog: disableAccessLog,
	})
}

// Start start the AppServer
func (server *AppServer) Start() {
	// process routers
	for basePattern, list := range server.basePatternRouterMap {
		if len(list) == 1 && list[0].PathConfig.IsPlainPath() {
			server.Handle(basePattern, list[0], list[0].Config.DisableAccessLog) // TODO: fix ugly access
		} else {
			// need a hub
			hub := NewRouterHub(basePattern)
			for _, router := range list {
				hub.AddRouter(router)
			}
			server.AddHub(hub)
		}
	}

	log.Fatal(server.Server.ListenAndServe())
}
