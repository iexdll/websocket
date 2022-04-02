package main

import (
	"encoding/json"
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

		buf, err := ioutil.ReadAll(r)
		if err != nil {
			log.Println(5, err)
			return
		}

		if isDevice {
			if len(Devices) > 0 {
				d := &Device{}
				_ = json.Unmarshal(buf, d)
				Devices[0].Lamp1 = d.Lamp1
				Devices[0].Lamp2 = d.Lamp2
				Devices[0].Temp = d.Temp
			}
			for _, item := range CPanels {
				item.Conn.write(buf)
			}
		} else {
			for _, item := range Devices {
				item.Conn.write(buf)
			}
		}

	}
}
