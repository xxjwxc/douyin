package msg

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type websocketClientManager struct {
	addr        *string
	path        string
	conn        *websocket.Conn
	sendMsgChan chan string
	recvMsgChan chan string
	isAlive     bool
	timeout     int
}

func NewWsClientManager(addrIp, addrPort, path string, timeout int) *websocketClientManager {
	addrString := addrIp + ":" + addrPort
	var sendChan = make(chan string, 10)
	var recvChan = make(chan string, 10)
	var conn *websocket.Conn
	td := &websocketClientManager{
		addr:        &addrString,
		path:        path,
		conn:        conn,
		sendMsgChan: sendChan,
		recvMsgChan: recvChan,
		isAlive:     false,
		timeout:     timeout,
	}
	td.dail()
	return td
}

// 链接服务端
func (wsc *websocketClientManager) dail() {
	var err error
	u := url.URL{Scheme: "wss", Host: *wsc.addr, Path: wsc.path}
	log.Printf("connecting to %s", u.String())
	dialer := websocket.Dialer{TLSClientConfig: &tls.Config{RootCAs: nil, InsecureSkipVerify: true}}
	wsc.conn, _, err = dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	wsc.isAlive = true
	log.Printf("connecting to %s 链接成功！！！", u.String())
	wsc.conn.Close()
}
