package server

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Websocket Service.
// Following https://tutorialedge.net/golang/go-websocket-tutorial/
type wsServer struct {
	conn *websocket.Conn
	ch   chan []byte
	mu   sync.Mutex
}

// New web socket server connection.
func New() *wsServer {
	return &wsServer{
		ch: make(chan []byte),
	}
}

// Opens the ws connection.
func (ws *wsServer) Open(w http.ResponseWriter, r *http.Request) error {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // skip CORS check for now
	// upgrade this connection to a WebSocket
	// connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	ws.conn = conn
	return nil
}

// Closes the ws connection.
func (ws *wsServer) Close() error {
	if ws.conn == nil {
		return errors.New("conn not initialized with Open()")
	}
	return ws.conn.Close()
}

// Reads a text message from the ws connection.
func (ws *wsServer) Read() (msg []byte, err error) {
	if ws.conn == nil {
		return msg, errors.New("conn not initialized with Open()")
	}
	// do not use mutex lock here as the read is blocking, and lock up other processes
	_, msg, err = ws.conn.ReadMessage()
	return msg, err
}

// Sends a text message over the ws connection
func (ws *wsServer) Write(msg []byte) error {
	if ws.conn == nil {
		return errors.New("conn not initialized with Open()")
	}
	ws.mu.Lock()
	defer ws.mu.Unlock()
	return ws.conn.WriteMessage(websocket.TextMessage, msg)
}

// Sends an empty message over the ws connection, for heartbeat purpose.
func (ws *wsServer) Ping() error {
	if ws.conn == nil {
		return errors.New("conn not initialized with Open()")
	}
	ws.mu.Lock()
	defer ws.mu.Unlock()
	return ws.conn.WriteMessage(websocket.PingMessage, nil)
}

// Whenever the ws client sends a message or pongs the server, execute the callback.
// For heartbeat purpose.
func (ws *wsServer) OnClientPong(callback func(appData string) error) {
	if ws.conn == nil {
		panic("conn not initialized with Open()")
	}
	ws.conn.SetPongHandler(callback)
}

// SetWriteDeadline sets the write deadline on the underlying network connection. After timing out, connection is closed and a closed error is returned.
// A zero value for t means writes will not time out.
func (ws *wsServer) SetWriteDeadline(t time.Time) error {
	if ws.conn == nil {
		return errors.New("conn not initialized with Open()")
	}
	ws.mu.Lock()
	defer ws.mu.Unlock()
	return ws.conn.SetWriteDeadline(t)
}

// SetWriteDeadline sets the read deadline on the underlying network connection. After timing out, connection is closed and a closed error is returned.
// A zero value for t means writes will not time out.
func (ws *wsServer) SetReadDeadline(t time.Time) error {
	if ws.conn == nil {
		return errors.New("conn not initialized with Open()")
	}
	ws.mu.Lock()
	defer ws.mu.Unlock()
	return ws.conn.SetReadDeadline(t)
}

// For custom processing in goroutines, e.g. broadcast a message from a chat hub to all ws connections.
func (ws *wsServer) Channel() (channel chan []byte) {
	return ws.ch
}

// Check whether websocket connection is closed, using the err returned from reading / writing to the connection.
func (ws *wsServer) IsClosedError(err error) bool {
	return websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err)
}
