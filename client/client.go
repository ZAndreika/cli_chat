package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"runtime"
	"time"

	logger "../utils/logger"
	messages "../utils/messages"
	postman "../utils/postman"
	router "../utils/router"

	"github.com/fanliao/go-promise"
)

var (
	currentRequestID   uint64
	login              string
	promises           []*promise.Promise
	sentMessages       []messages.Message
	sentMessagesStatus []bool
)

func init() {
	if len(os.Args) < 3 {
		logger.Error("Need login and password")
		os.Exit(1)
	}
	currentRequestID = 0
	login = os.Args[1]
	promises = make([]*promise.Promise, 100, 100)
	sentMessages = make([]messages.Message, 100, 100)
	sentMessagesStatus = make([]bool, 100, 100)
}

func main() {
	conn, err := net.Dial("tcp", "0.0.0.0:8080")
	if err != nil {
		logger.Error(err.Error())
		return
	}
	pm := postman.New(conn)
	defer pm.Dismiss()

	password := os.Args[2]
	if !authenticate(pm, login, password) {
		logger.Error("Authentication error")
		return
	}
	logger.Info(fmt.Sprintf("Authentication success with login: %v", login))

	r := router.New()
	r.RegisterRoute("/message", handleMessage)
	r.RegisterRoute("/exit", handleExit)

	input := make(chan string)
	serverChan := make(chan messages.Message)

	go updateMessages(input)
	go getServerMessages(pm, serverChan)

	for {
		select {
		case i := <-input:
			if i == "/exit" {
				sendExitCommand(pm)
			} else {
				sendMessage(pm, i)
			}
		case msg := <-serverChan:
			if handler, ok := r.GetHandlerByRoute(msg.Route); ok {
				if status := handler(pm, msg); status == router.StatusEXIT {
					return
				}
			}
		case <-time.After(500 * time.Millisecond):
			continue
		}
	}
}

func sendExitCommand(pm *postman.Postman) {
	msg := messages.Message{Route: "/exit", Author: login}
	pm.Send(msg)
}

func sendMessage(pm *postman.Postman, msg string) {
	msgStruct := messages.Message{
		Route: "/message",
		Data: messages.Data{
			Text: msg,
			Time: time.Now().String(),
		},
		RequestID: currentRequestID,
		Author:    login,
	}

	sentMessages[currentRequestID] = msgStruct
	promises[currentRequestID] = promise.NewPromise()
	promises[currentRequestID].OnSuccess(func(v interface{}) {
		id := v.(uint64)
		sentMessagesStatus[id] = true
		printMyMessage(sentMessages[id])
	})
	sentMessagesStatus[currentRequestID] = false

	currentRequestID = currentRequestID + 1

	pm.Send(msgStruct)
}

func updateMessages(c chan string) {
	for {
		in := bufio.NewReader(os.Stdin)
		result, err := in.ReadString('\n')
		if err != nil {
			logger.Error(err.Error())
			break
		}

		if runtime.GOOS == "windows" {
			result = result[:len(result)-2]
		}

		if len(result) != 0 {
			c <- result
		}
	}
}

func getServerMessages(pm *postman.Postman, c chan messages.Message) {
getServerMessagesLoop:
	for {
		pm.Conn.SetReadDeadline(time.Now().Add(time.Second * 5))
		msg, err := pm.Receive()
		if err != nil {
			if e, ok := err.(interface{ Timeout() bool }); ok && e.Timeout() {
			} else if err.Error() == postman.PostmanEmptyMessageError {
			} else {
				logger.Error(err.Error())
				break getServerMessagesLoop
			}
		}
		if len(msg.Data.Text) > 0 {
			c <- msg
		}
	}
}

func authenticate(pm *postman.Postman, login string, password string) bool {
	aData := messages.AuthData{Login: login, Password: password}
	res, err := json.Marshal(aData)
	if err != nil {
		logger.Warning(err.Error())
	}
	d := messages.Data{Text: string(res), Time: time.Now().String()}
	msg := messages.Message{Route: "/auth/login", Data: d, RequestID: currentRequestID, Author: login}
	currentRequestID = currentRequestID + 1
	pm.Send(msg)
	recvMsg, _ := pm.Receive()
	sentMessagesStatus[0] = true

	if recvMsg.Data.Text == "error" {
		return false
	}

	return true
}

func printMessage(msg messages.Message) {
	fmt.Printf("\t\t=============================\n")
	fmt.Printf("\t\t%s\n", msg.Data.Time[:19])
	fmt.Printf("\t\t%s:%s\n", msg.Author, msg.Data.Text)
	fmt.Printf("\t\t=============================\n")
}

func printMyMessage(msg messages.Message) {
	fmt.Println("=============================")
	fmt.Println(msg.Data.Time[:19])
	fmt.Println(msg.Data.Text)
	fmt.Println("=============================")
}
