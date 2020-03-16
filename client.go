package douyu

import (
	"encoding/binary"
	"sync"
	"time"

	"github.com/go-sdk/logx"
	ws "github.com/gorilla/websocket"
)

const (
	danmuURL  = "wss://danmuproxy.douyu.com:8501/"
	maxLength = 10240
)

type Client struct {
	rid     string
	conn    *ws.Conn
	finish  chan bool
	message chan map[string]string
	wMu     sync.Mutex
}

func NewClient() *Client {
	c := &Client{
		finish:  make(chan bool, 1),
		message: make(chan map[string]string),
		wMu:     sync.Mutex{},
	}
	return c
}

func (c *Client) SetRoomId(rid string) *Client {
	c.rid = rid
	return c
}

func (c *Client) Start() (chan bool, error) {
	conn, _, err := ws.DefaultDialer.Dial(danmuURL, nil)
	if err != nil {
		return nil, err
	}

	c.conn = conn

	go c.listen()
	go c.keepAlive()

	c.write(Encode("type", "loginreq", "roomid", c.rid))
	c.write(Encode("type", "joingroup", "rid", c.rid, "gid", "-9999"))

	return c.finish, nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.write(Encode("type", "logout"))
		_ = c.conn.Close()
	}
}

func (c *Client) read() []byte {
	if c.conn == nil {
		return nil
	}
	_, bs, err := c.conn.ReadMessage()
	if err != nil {
		c.Close()
		logx.Errorf("websocket: read error: %v", err)
		return nil
	}
	return bs
}

func (c *Client) write(bs []byte) {
	if c.conn == nil {
		return
	}
	c.wMu.Lock()
	defer c.wMu.Unlock()
	err := c.conn.WriteMessage(ws.BinaryMessage, bs)
	if err != nil {
		logx.Errorf("websocket: write error: %v", err)
	}
}

func (c *Client) listen() {
	defer close(c.finish)
	for {
		go c.handle(c.read())
	}
}

func (c *Client) handle(bs []byte) {
	defer func() { recover() }()
	if len(bs) <= 4 {
		return
	}
	length := int(binary.LittleEndian.Uint32(bs[:4]))
	if length > maxLength {
		return
	}
	go func() { c.message <- Decode(bs[:length+4]) }()
	go c.handle(bs[length+4:])
}

func (c *Client) GetMessage() chan map[string]string {
	return c.message
}

func (c *Client) keepAlive() {
	time.Sleep(5 * time.Second)
	for {
		c.write(Encode("type", "mrkl"))
		time.Sleep(30 * time.Second)
	}
}
