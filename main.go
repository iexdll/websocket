package main

import (
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
var Status string

func main() {

	Status = "yellow"
	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Новый клиент")
		c := &Client{ws: ws}
		go c.readMsg()
		c.writeMsg(Status)

		Clients = append(Clients, c)
	})

	log.Println("Запуск сервера")
	http.ListenAndServe(":8654", nil)
}

func (c *Client) writeMsg(msg string) {
	w, err := c.ws.NextWriter(websocket.TextMessage)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write([]byte(msg))
	w.Close()
}

func (c *Client) readMsg() {
	for {
		_, r, err := c.ws.NextReader()
		if err != nil {
			log.Println(err)
			return
		}
		buf, err := ioutil.ReadAll(r)
		if err != nil {
			log.Println(err)
			return
		}
		msg := string(buf)
		if "click" == msg {
			if Status == "red" {
				Status = "green"
			} else {
				Status = "red"
			}
			for _, cl := range Clients {
				cl.writeMsg(Status)
			}
		}
	}
}
