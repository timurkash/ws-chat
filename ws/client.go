package ws

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type (
	RegisterFunc  func(*Client)
	BroadcastFunc func([]byte)

	Client struct {
		Conn      *websocket.Conn
		Send      chan []byte
		Register  RegisterFunc
		Broadcast BroadcastFunc
	}
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (c *Client) ReadPump() {
	defer func() {
		c.Register(c)
		if err := c.Conn.Close(); err != nil {
			log.Println(err)
			//log.Fatalln(err)
		}
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		//log.Fatalln(err)
	}
	c.Conn.SetPongHandler(func(string) error {
		if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			log.Println(err)
			//log.Fatalln(err)
		}
		return nil
	})
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.ReplaceAll(message, newline, space))
		c.Broadcast(message)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := c.Conn.Close(); err != nil {
			log.Println(err)
			//log.Fatalln(err)
		}
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Println(err)
				//log.Fatalln(err)
			}
			if !ok {
				if err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Println(err)
					//log.Fatalln(err)
				}
				log.Println("websocket.CloseMessage")
				return
			}
			writeCloser, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println("conn.NextWriter error")
				return
			}
			if _, err := writeCloser.Write(message); err != nil {
				log.Println("writeCloser.Write error")
				return
			}
			{
				// Add queued chat messages to the current websocket message.
				for i := 0; i < len(c.Send); i++ {
					if _, err := writeCloser.Write(newline); err != nil {
						log.Println(err)
						//log.Fatalln(err)
					}
					if _, err := writeCloser.Write(<-c.Send); err != nil {
						log.Println(err)
						//log.Fatalln(err)
					}
				}
			}
			if err := writeCloser.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Println(err)
				//log.Fatalln(err)
			}
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
