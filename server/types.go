package main

import messages "../utils/messages"

type (
	// User - model of user
	User struct {
		Login    string
		Password string
	}

	message = messages.Message

	dbdata struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		DBname   string `json:"dbname"`
	}

	netdata struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
	}

	config struct {
		Mysql   dbdata  `json:"mysql"`
		Network netdata `json:"network"`
	}
)
