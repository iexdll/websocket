package main

import (
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
)

type Connection struct {
	conn *websocket.Conn
}

func (c *Connection) write(message []byte) {

	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		log.Println(1, err)
		return
	}

	_, err = w.Write(message)
	if err != nil {
		log.Println(2, err)
		return
	}

	err = w.Close()
	if err != nil {
		log.Println(3, err)
		return
	}

}

func (c *Connection) read(isDevice bool) {
	for {

		_, r, err := c.conn.NextReader()
		if err != nil {
			log.Println(4, err)
			return
		}

		message, err := ioutil.ReadAll(r)
		if err != nil {
			log.Println(5, err)
			return
		}

		if isDevice {
			for _, item := range CPanels {
				item.Conn.write(message)
			}
		} else {
			for _, item := range Devices {
				item.Conn.write(message)
			}
		}

	}
}
