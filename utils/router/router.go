package router

import (
	"../messages"
	"../postman"
)

// STATUS - Required handlers return type
type STATUS = uint64

// Required handlers return statuses
const (
	StatusOK    uint64 = iota // StatusOK - expected function's return
	StatusERROR               // StatusERROR - if some went wrong
	StatusEXIT                // StatusEXIT - if user sent exit request
)

// RouteHandler type of handler for route
type RouteHandler func(pm *postman.Postman, msg messages.Message) STATUS

// Router - class for route messages
type Router struct {
	handlers map[string]RouteHandler
}

// New creates router
func New() Router {
	var r Router
	r.handlers = make(map[string]RouteHandler)
	return r
}

// GetHandlerByRoute returns handler by route if it exists
func (r Router) GetHandlerByRoute(route string) (RouteHandler, bool) {
	command, ok := r.handlers[route]
	return command, ok
}

// RegisterRoute register route
func (r Router) RegisterRoute(route string, action RouteHandler) {
	r.handlers[route] = action
}

// RegisterRoutesMap register map [route]Handler
func (r Router) RegisterRoutesMap(routesMap map[string]RouteHandler) {
	for route, handler := range routesMap {
		r.handlers[route] = handler
	}
}
