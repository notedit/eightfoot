// Date: 2012-06-23
// Author: notedit<notedit@gmail.com>

package rpc

import (
	"bufio"
	"errors"
	"io"
	"net"
	"sync"
	"time"
    
    "github.com/ugorji/go-msgpack"
)

// timeout
const DefaultTimeout = time.Duration(10000) * time.Millisecond

// connection pool number
const DefaultConnectionPool = 10

type Client struct {
	addr     net.Addr
	mutex    sync.Mutex
	Timeout  time.Duration
	freeconn []*conn
}


type conn struct {
	rwc		io.ReadWriteCloser
	dec 	*msgpack.Decoder
	enc 	*msgpack.Encoder
	encBuf	*bufio.Writer
	c  		*Client
}

func NewClient(cn io.ReadWriteCloser,c *Client) *conn{
	encBuf := bufio.NewWriter(cn)
	con := &conn{cn,msgpack.NewDecoder(cn,nil),msgpack.NewEncoder(cn),encBuf,c}
	return con
}


func (cn *conn) WriteRequest(req *clientRequest, body interface{}) (err error) {
	if err = cn.enc.Encode(req);err != nil {
		return
	}
	if err = cn.enc.Encode(body);err != nil {
		return
	}
	return cn.encBuf.Flush()
}

func (cn *conn) ReadResponse(res *clientResponse, reply interface{}) (err error) {
	if err = cn.ReadResponseHeader(res); err != nil {
		return
	}
	if res.Operation == 3 {
		cn.ReadResponseBody(&struct{}{})
		return errors.New(res.Error)
	}
	err = cn.ReadResponseBody(reply)
	return
}

func (cn *conn) ReadResponseHeader(res *clientResponse) (err error) {
	return cn.dec.Decode(res)
}

func (cn *conn) ReadResponseBody(reply interface{}) (err error) {
	return cn.dec.Decode(reply)
}

type clientRequest struct {
	Operation uint8
	Method    string
}

type clientResponse struct {
	Operation uint8
	Error     string
}

func New(server string) *Client {
	addr, err := net.ResolveTCPAddr("tcp", server)
	if err != nil {
		panic(err)
	}
	return &Client{addr: addr, freeconn: make([]*conn, 0, DefaultConnectionPool)}
}

func (c *Client) dial() (net.Conn, error) {
	cn, err := net.Dial(c.addr.Network(), c.addr.String())
	if err != nil {
		return nil, err
	}
	return cn, nil
}

func (c *Client) getConn() (*conn, error) {
	cn, ok := c.getFreeConn()
	if ok {
		return cn, nil
	}

	nc, err := c.dial()
	if err != nil {
		return nil, err
	}

	return NewClient(nc,c), nil
}


func (c *Client) getFreeConn() (*conn, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if len(c.freeconn) == 0 {
		return nil, false
	}
	cn := c.freeconn[len(c.freeconn)-1]
	c.freeconn = c.freeconn[:len(c.freeconn)-1]
	return cn, true
}

func (c *Client) release(cn *conn) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if len(c.freeconn) >= DefaultConnectionPool {
		cn.rwc.Close()
		return
	}
	c.freeconn = append(c.freeconn, cn)
}

func (c *Client) call(req *clientRequest, args interface{}, reply interface{}) (err error) {
	cn, err := c.getConn()
	if err != nil {
		return
	}
	defer c.release(cn)
	if err = cn.WriteRequest(req, args); err != nil {
		return err
	}
	res := &clientResponse{}
	if err = cn.ReadResponse(res, reply); err != nil {
		return
	}
	return
}

func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	req := new(clientRequest)
	req.Method = serviceMethod
	req.Operation = uint8(1)
	err := c.call(req, args, reply)
	if err != nil {
		return err
	}
	return nil
}
