package socket

import (
	"cetm/qapi/x/mlog"
)

// Box handles ws request
//

var boxLog = mlog.NewTagLog("box")

type Box struct {
	ID       string
	Clients  *WsClientManager
	handlers map[string]IBoxHandler
	NotFound IBoxHandler
	Join     func(*WsClient) error
	Leave    func(*WsClient)
	Recover  func(*Request, interface{})
}

// NewBox create a new box
func NewBox(ID string) *Box {
	var b = &Box{
		ID:       ID,
		Clients:  NewWsClientManager(),
		handlers: make(map[string]IBoxHandler),
	}
	b.Recover = b.defaultRecover
	// b.NotFound = b.notFound
	b.Join = b.join
	b.Leave = b.leave
	b.Handle("/echo", b.Echo)
	return b
}

// Handle add a handler
func (b *Box) Handle(uri string, handler IBoxHandler) {
	b.handlers[uri] = handler
}

// Serve process the regggo getquest
func (b *Box) Serve(r *Request) {

	defer func() {
		if rc := recover(); rc != nil {
			if nil != b.Recover {
				b.Recover(r, rc)
			}
		}
	}()

	var handler = b.handlers[r.Path()]
	if handler == nil {
		if nil == b.NotFound {
			return
		}
		handler = b.NotFound
	}
	handler(r)
}

// Echo the default echo service
func (b *Box) Echo(r *Request) {
	r.Client.queueForSend(r.Payload)
}

func (b *Box) Broadcast(uri string, v interface{}) {
	b.Clients.SendJson(uri, v)
}

func (b *Box) Destroy() {
	b.Clients.Destroy()
}

func (b *Box) GetStatus() interface{} {
	return map[string]interface{}{
		"active_clients": b.Clients.Count(),
	}
}
