package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

//https://udf.su/ws-1006-error-handling
//https://stackoverflow.com/questions/37696527/go-gorilla-websockets-on-ping-pong-fail-user-disconnct-call-function

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var Devices []*Device
var CPanels []*CPanel

func main() {

	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}

		deviceID := r.URL.Query().Get("device")

		if "" == deviceID {
			log.Println("Новая панель управления")
			cp := &CPanel{Conn: &Connection{conn: conn}}
			go cp.Conn.read(false)
			CPanels = append(CPanels, cp)
		} else {
			log.Println("Устройство подключено")
			device := &Device{Conn: &Connection{conn: conn}}
			go device.Conn.read(true)
			Devices = append(Devices, device)
		}

	})

	log.Println("Запуск сервера")
	http.ListenAndServe(":8654", nil)
}
