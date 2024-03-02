package server

import (
	"net/http"
	"time"
)

type WSServer interface {
	// Opens the ws connection.
	Open(w http.ResponseWriter, r *http.Request) error
	// Closes the ws connection.
	Close() error
	// Reads a text message from the ws connection.
	Read() (msg []byte, err error)
	// Sends a text message over the ws connection
	Write(msg []byte) error
	// For custom processing in goroutines, e.g. broadcast a message from a chat hub to all ws connections.
	Channel() (channel chan []byte)
	// Sends an empty message over the ws connection, for heartbeat purpose.
	Ping() error
	// Whenever the ws client replies with a pongs the server, execute the callback.
	// For heartbeat purpose.
	OnClientPong(callback func(appData string) error)
	// SetWriteDeadline sets the write deadline on the underlying network connection. After timing out, connection is closed and a closed error is returned.
	// A zero value for t means writes will not time out.
	SetWriteDeadline(t time.Time) error
	// SetWriteDeadline sets the read deadline on the underlying network connection. After timing out, connection is closed and a closed error is returned.
	// A zero value for t means writes will not time out.
	SetReadDeadline(t time.Time) error
	// Check whether websocket connection is closed, using the err returned from reading / writing to the connection.
	IsClosedError(err error) bool
}
