package svc

import (
	"api/internal/config"
	"api/ws"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config       config.Config
	MysqlDB      *gorm.DB
	WebSocketCtx *ws.WebSocketContext
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.MysqlConfig.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:  c,
		MysqlDB: db,
		WebSocketCtx: &ws.WebSocketContext{
			SocketMap: make(map[string]*ws.UserWebSocket, 1024),
		},
	}
}
