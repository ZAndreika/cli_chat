package main

import (
	messages "../utils/messages"
	postman "../utils/postman"
	router "../utils/router"
)

func handleMessage(pm *postman.Postman, msg messages.Message) router.STATUS {
	if msg.Author == "server" {
		promises[msg.RequestID].Resolve(msg.RequestID)
	} else {
		printMessage(msg)
	}
	return router.StatusOK
}

func handleExit(pm *postman.Postman, msg messages.Message) router.STATUS {
	return router.StatusEXIT
}
