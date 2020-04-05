package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	logger "../utils/logger"
	messages "../utils/messages"
	postman "../utils/postman"
	router "../utils/router"
)

var (
	r          router.Router
	postmans   []*postman.Postman
	serverName string
	db         *gorm.DB
)

func init() {
	serverName = "server"
}

func main() {
	defer panicСatch()

	cfg := readConfig("config.json")

	ln := createNetListener(cfg)

	db = createDbConnection(cfg)
	defer db.Close()
	prepareDB(db)

	r = router.New()
	r.RegisterRoutesMap(routesHandlers)

	go kickNoAuthenticated()
	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		pm := postman.New(conn)
		postmans = append(postmans, pm)
		go handleConnection(pm)
	}
}

func handleConnection(pm *postman.Postman) {
	defer closeConnection(pm)

	logger.Info(fmt.Sprint("New client: ", pm.Conn.RemoteAddr()))

	for {
		msg, err := pm.Receive()
		if err != nil {
			if err.Error() == postman.PostmanEmptyMessageError {
				logger.Warning(fmt.Sprintf("%v, conn: %v", err.Error(), pm.Conn.RemoteAddr()))
			}
			continue
		}
		logger.Info(msg.String())

		if handler, ok := r.GetHandlerByRoute(msg.Route); ok {
			if status := handler(pm, msg); status != router.StatusOK {
				break
			}
		} else {
			msgData := messages.Data{Text: "No such route", Time: time.Now().String()}
			msgSend := message{Route: "/error", Data: msgData, RequestID: msg.RequestID, Author: serverName}
			pm.Send(msgSend)
		}
	}
}

func closeConnection(pm *postman.Postman) {
	for i, p := range postmans {
		if p == pm {
			postmans = append(postmans[:i], postmans[i+1:]...) // remove from slice
			break
		}
	}
	pm.Dismiss()
}

func createNetListener(cfg config) net.Listener {
	listenAddress := cfg.Network.Address + ":" + strconv.Itoa(cfg.Network.Port)
	ln, err := net.Listen("tcp", listenAddress)
	if err != nil {
		panic(err)
	}

	logger.Debug(fmt.Sprint("Start server on: ", ln.Addr()))
	return ln
}

func createDbConnection(cfg config) *gorm.DB {
	connectionString := cfg.Mysql.Login + ":" + cfg.Mysql.Password + "@/" + cfg.Mysql.DBname + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}
	return db
}

func prepareDB(db *gorm.DB) {
	if !db.HasTable(&User{}) {
		logger.Debug("No table Users, creating")
		db.CreateTable(&User{})
	}
}

func readConfig(filename string) (cfg config) {
	configFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	configBytes, _ := ioutil.ReadAll(configFile)

	err = json.Unmarshal(configBytes, &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

func kickNoAuthenticated() {
	for {
		time.Sleep(5 * time.Second)
		for _, p := range postmans {
			if !p.IsAuthenticated() {
				closeConnection(p)
			}
		}
	}
}

func panicСatch() {
	if err := recover(); err != nil {
		logger.Error(fmt.Sprintf("PANIC: %v", err))
	}
}
