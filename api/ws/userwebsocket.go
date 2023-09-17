package ws

import (
	"github.com/gorilla/websocket"
	"sync"
)

type UserWebSocket struct {
	UserId         string
	Conn           map[string]*websocket.Conn
	ReceiveMsgChan chan WebSocketMessage // Receive Message from other Websockets
	lock           sync.Mutex
	//	MqClient
}

type WebSocketMessage struct {
	From   string
	To     string
	MsgTxt string
}

type UserWebSocketMap struct {
	socketMap map[string]*UserWebSocket
	l         sync.Mutex
}

func (uws *UserWebSocket) addWs(addr string, ws *websocket.Conn) {
	uws.lock.Lock()
	defer uws.lock.Unlock()
	uws.Conn[addr] = ws
}

func (uwsm *UserWebSocketMap) AddWs(userId string, ws *websocket.Conn) {
	v, ok := uwsm.socketMap[userId]
	if !ok {
		uwsm.l.Lock()
		defer uwsm.l.Unlock()
		v, ok = uwsm.socketMap[userId]
		if !ok {
			conMap := make(map[string]*websocket.Conn, 3)
			conMap[ws.RemoteAddr().String()] = ws
			v = &UserWebSocket{
				UserId:         userId,
				Conn:           conMap,
				ReceiveMsgChan: make(chan WebSocketMessage, 10),
				lock:           sync.Mutex{},
			}
			uwsm.socketMap[userId] = v
			return
		}

		v.Conn[ws.RemoteAddr().String()] = ws
		return
	}
	v.addWs(ws.RemoteAddr().String(), ws)
}
