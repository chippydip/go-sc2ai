package client

import (
	"fmt"
	"net/url"
	"time"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
)

type connection struct {
	Status api.Status

	urlStr  string
	timeout time.Duration

	conn     *websocket.Conn
	requests chan request
}

type request struct {
	*api.Request
	callback chan response
}

type response struct {
	*api.Response
	error
}

// Connect ...
func (c *connection) Connect(address string, port int, timeout time.Duration) error {
	c.Status = api.Status_unknown

	// Save the connection info in case we need to re-connect
	host := fmt.Sprintf("%v:%v", address, port)
	u := url.URL{Scheme: "ws", Host: host, Path: "/sc2api"}
	c.urlStr = u.String()

	conn, _, err := websocket.DefaultDialer.Dial(c.urlStr, nil)
	if err != nil {
		return err
	}
	c.conn = conn

	c.requests = make(chan request)
	callbacks := make(chan chan<- response)

	c.conn.SetCloseHandler(func(code int, text string) error {
		//control.Error(ClientError_ConnectionClosed)
		close(c.requests)
		return nil
	})

	// Send worker
	go func() {
		for r := range c.requests {
			data, err := proto.Marshal(r.Request)
			if err != nil {
				r.callback <- response{nil, err}
				continue
			}
			err = c.conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				r.callback <- response{nil, err}
				continue
			}
			callbacks <- r.callback
		}
		close(callbacks)
	}()

	// Receive worker
	go func() {
		for cb := range callbacks {
			_, data, err := c.conn.ReadMessage()
			if err != nil {
				cb <- response{nil, err}
				continue
			}

			r := &api.Response{}
			err = proto.Unmarshal(data, r)
			if err != nil {
				cb <- response{nil, err}
				continue
			}

			cb <- response{r, c.onResponse(r)}
		}
	}()

	_, err = c.ping(api.RequestPing{})()
	return err
}

func (c *connection) onResponse(r *api.Response) error {
	if r.Status != api.Status_nil {
		c.Status = r.Status
	}
	// for _, e := range r.Error {
	// 	// TODO: error callback
	// }
	if len(r.Error) > 0 {
		return fmt.Errorf("%v", r.Error)
	}
	return nil
}

func (c *connection) request(req *api.Request) func() (*api.Response, error) {
	out := make(chan response, 1)
	c.requests <- request{req, out}
	return func() (*api.Response, error) {
		r := <-out
		return r.Response, r.error
	}
}
