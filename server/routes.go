package main

import (
	"../utils/router"
)

var routesHandlers = map[string]router.RouteHandler{
	"/auth/login": handleAuthLogin,
	"/message":    handleBroadcastMessage,
	"/exit":       handleExit,
}
