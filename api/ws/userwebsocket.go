package ws

import (
	"api/config"
	"api/kafkacli"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
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

type WebSocketContext struct {
	SocketMap map[string]*UserWebSocket
	l         sync.Mutex
}

func (uws *UserWebSocket) addWs(addr string, ws *websocket.Conn) {
	uws.lock.Lock()
	defer uws.lock.Unlock()
	uws.Conn[addr] = ws
	ws.SetReadDeadline(time.Now().Add(config.PongTimePeriod))
	ws.SetPongHandler(
		func(string) error {
			ws.SetReadDeadline(time.Now().Add(config.PongTimePeriod))
			return nil
		})

}

// WsWriteTask write to mq asynchronously such as kafkacli
func (*UserWebSocket) WsWriteTask(ws *websocket.Conn) {

	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			logx.Errorf("error reading message from websocket:%v", err)
			return
		}
		// p will be written to Mq
		msg := string(p)
		err = kafkacli.PushData2Kafka(msg)

		if err != nil {
			logx.Errorf("Error occurs when push data [%s] to kafka:%v", msg, err)
		}

	}
}

// WsReadTask  read from pub server by grpc
func (uws *UserWebSocket) WsReadTask(ws *websocket.Conn) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case msg, ok := <-uws.ReceiveMsgChan:
			{
				if !ok {
					return
				}
				msgByte, err := json.Marshal(msg)
				if err != nil {
					logx.Errorf("Error occurs when parsing data from pub server ")
					break
				}
				ws.WriteMessage(websocket.BinaryMessage, msgByte)
			}
		case <-ticker.C:
			{
				ws.SetWriteDeadline(time.Now().Add(config.PingTimePeriod))
				if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}

}

func (uwsm *WebSocketContext) AddWs(userId string, ws *websocket.Conn) {
	v, ok := uwsm.SocketMap[userId]
	if !ok {
		uwsm.l.Lock()
		defer uwsm.l.Unlock()
		v, ok = uwsm.SocketMap[userId]
		if !ok {
			conMap := make(map[string]*websocket.Conn, 3)
			conMap[ws.RemoteAddr().String()] = ws
			v = &UserWebSocket{
				UserId:         userId,
				Conn:           conMap,
				ReceiveMsgChan: make(chan WebSocketMessage, 10),
				lock:           sync.Mutex{},
			}
			uwsm.SocketMap[userId] = v
			go v.WsReadTask(ws)
			go v.WsWriteTask(ws)
			return
		}

		v.Conn[ws.RemoteAddr().String()] = ws
		return
	}
	v.addWs(ws.RemoteAddr().String(), ws)
	go v.WsReadTask(ws)
	go v.WsWriteTask(ws)
}
