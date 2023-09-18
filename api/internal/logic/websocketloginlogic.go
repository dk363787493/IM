package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"

	"api/internal/svc"
	"api/internal/types"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocketLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWebSocketLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebSocketLoginLogic {
	return &WebSocketLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WebSocketLoginLogic) WebSocketLogin(w http.ResponseWriter, r *http.Request, req *types.LoginReq) (resp *types.LoginResp, err error) {
	// todo: add your logic here and delete this line
	userId := req.UserId
	pwd := req.Password
	queyPwd := ""
	l.svcCtx.MysqlDB.Raw("select password from userinfo where username=?", userId).Scan(&queyPwd)
	if pwd != queyPwd {
		return nil, errors.New(" The password is incorrect.Please try it again")
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logx.Error(err)
		return
	}
	l.svcCtx.WebSocketCtx.AddWs(userId, conn)
	return
}
