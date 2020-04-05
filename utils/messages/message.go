package messages

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Data - part of message
type Data struct {
	Text string `json:"text"`
	Time string `json:"time"`
}

// Message - specification of message between client and server
type Message struct {
	Route     string `json:"route"`
	Data      Data   `json:"data"`
	RequestID uint64 `json:"requestId"`
	Author    string `json:"author"`
}

// AuthData - data for user's authentication
type AuthData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Deserialize - Transforms byte slice to Message struct
func Deserialize(msg []byte) (Message, error) {
	var res Message
	err := json.Unmarshal(msg, &res)
	if err != nil {
		err = errors.New("cannot deserialize byte array to Message")
	}
	return res, err
}

// Serialize - Transforms Message struct to byte slice
func Serialize(msg Message) ([]byte, error) {
	res, err := json.Marshal(msg)
	if err != nil {
		err = errors.New("cannot serialize Message to byte array")
	}
	return res, err
}

func (m *Message) String() string {
	return fmt.Sprintf("{ Route: %v, Data: { Text: %v, Time: %v }, RequestID: %v, Author: %v }", m.Route, m.Data.Text, m.Data.Time, m.RequestID, m.Author)
}
