package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	ws *websocket.Conn
}

var Clients []*Client

func main() {

	Device = &DeviceState{}

	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Новый клиент")
		c := &Client{ws: ws}
		go c.readMsg()
		c.writeMsg()

		Clients = append(Clients, c)
	})

	log.Println("Запуск сервера")
	http.ListenAndServe(":8654", nil)
}

func (c *Client) writeMsg() {

	w, err := c.ws.NextWriter(websocket.TextMessage)
	if err != nil {
		log.Println(1, err)
		return
	}

	deviceJson, err := json.Marshal(Device)

	w.Write(deviceJson)
	w.Close()

}

type ClientMessage struct {
	Action string `json:"action"`
	Value  string `json:"value"`
}

type DeviceState struct {
	Lamp1 bool   `json:"lamp1"`
	Lamp2 bool   `json:"lamp2"`
	Temp  string `json:"temp"`
}

var Device *DeviceState

func (c *Client) readMsg() {
	for {
		_, r, err := c.ws.NextReader()
		if err != nil {
			log.Println(err)
			return
		}

		message := &ClientMessage{}
		b, err := ioutil.ReadAll(r)
		_ = json.Unmarshal(b, &message)

		if message.Action == "click" {
			if message.Value == "lamp1" {
				Device.Lamp1 = !Device.Lamp1
			} else if message.Value == "lamp2" {
				Device.Lamp2 = !Device.Lamp2
			}
		} else if message.Action == "temp" {
			Device.Temp = message.Value
		}

		for _, cl := range Clients {
			cl.writeMsg()
		}

	}
}
