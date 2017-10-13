package socket

import (
	"runtime/debug"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type WsServer struct{}

func (s *WsServer) WriteError(ws *websocket.Conn, err error) {
	data := BuildErrorMessage("/system", err)
	ws.WriteMessage(websocket.TextMessage, data)
}

func (s *WsServer) Recover(ws *websocket.Conn) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			if _, ok = err.(IWebError); ok {
				s.WriteError(ws, err)
				return
			}
			glog.Error(err, string(debug.Stack()))
			s.WriteError(ws, errInternalServer)
		} else {
			glog.Error(r, string(debug.Stack()))
			s.WriteError(ws, errInternalServer)
		}
	}
}
