package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	logger "../utils/logger"
	messages "../utils/messages"
	postman "../utils/postman"
	router "../utils/router"
)

func handleBroadcastMessage(pm *postman.Postman, msg message) router.STATUS {
	serializedMsg, err := messages.Serialize(msg)
	if err != nil {
		logger.Error(err.Error())
		return router.StatusERROR
	}

	for _, p := range postmans {
		if p == pm {
			msgData := messages.Data{Text: "ok", Time: time.Now().String()}
			msgSend := message{Route: "/message", Data: msgData, RequestID: msg.RequestID, Author: serverName}
			p.Send(msgSend)
			continue
		}
		p.SendBytes(serializedMsg)
	}
	return router.StatusOK
}

func handleAuthLogin(pm *postman.Postman, msg message) router.STATUS {
	var aData messages.AuthData
	err := json.Unmarshal([]byte(msg.Data.Text), &aData)
	if err != nil {
		logger.Error(fmt.Sprintf("cannot deserialize %v to authData", msg.Data.Text))
		return router.StatusERROR
	}

	hash := md5.Sum([]byte(aData.Password))
	passHash := hex.EncodeToString(hash[:])

	var user User
	db.Where("login = ?", aData.Login).First(&user)

	if user.Login != aData.Login { // if user is empty after db query
		user = User{Login: aData.Login, Password: passHash}
		db.Create(&user)
	}

	status := "success"
	retStatus := router.StatusOK
	pm.SetAuthenticate(true)

	if user.Password != passHash {
		status = "error"
		retStatus = router.StatusERROR
		pm.SetAuthenticate(false)
		logger.Warning(fmt.Sprint("failed authentication for user: ", user.Login, " from ", pm.Conn.RemoteAddr()))
	}

	msgData := messages.Data{Text: status, Time: time.Now().String()}
	msgSend := message{Route: msg.Route, Data: msgData, RequestID: msg.RequestID, Author: serverName}
	pm.Send(msgSend)

	return retStatus
}

func handleExit(pm *postman.Postman, msg message) router.STATUS {
	msgData := messages.Data{Text: "bye, " + msg.Author, Time: time.Now().String()}
	msgSend := message{Route: "/exit", Data: msgData, RequestID: msg.RequestID, Author: serverName}

	pm.Send(msgSend)
	logger.Info(fmt.Sprint("Close connection with ", pm.Conn.RemoteAddr()))

	return router.StatusEXIT
}
