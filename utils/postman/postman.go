package postman

import (
	"errors"
	"net"

	"../messages"
)

// MaxMsgSize - max size of message
const MaxMsgSize = 1024

// PostmanEmptyMessageError - specific string for empty message error
var PostmanEmptyMessageError = "Empty message"

// Postman - struct for send and recieve messages
type Postman struct {
	Conn   net.Conn
	isAuth bool
}

// New - constructor for Postman
func New(c net.Conn) *Postman {
	return &Postman{Conn: c, isAuth: false}
}

// SetAuthenticate - set authentication status
func (p *Postman) SetAuthenticate(value bool) {
	p.isAuth = value
}

// IsAuthenticated - check is postman authenticated
func (p *Postman) IsAuthenticated() bool {
	return p.isAuth
}

// Receive - read bytes and deserialize it to message
func (p *Postman) Receive() (messages.Message, error) {
	var resMsg messages.Message

	readBytes := make([]byte, MaxMsgSize)

	n, err := p.Conn.Read(readBytes)
	if err != nil {
		return resMsg, err
	}

	if n <= 0 {
		return resMsg, errors.New(PostmanEmptyMessageError)
	}

	recvMsg := make([]byte, MaxMsgSize)
	recvMsg = readBytes[:n]

	return messages.Deserialize(recvMsg)
}

// Send - serialize message and send bytes array
func (p *Postman) Send(msg messages.Message) {
	serializedMsg, err := messages.Serialize(msg)
	if err != nil {
		return // bad! maybe need to cache message and try to send again
	}
	p.Conn.Write(serializedMsg)
}

// SendBytes - send bytes array
func (p *Postman) SendBytes(bytes []byte) {
	p.Conn.Write(bytes)
}

// Dismiss - close connection
func (p *Postman) Dismiss() {
	p.Conn.Close()
}
