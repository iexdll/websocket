package main

type Device struct {
	Conn  *Connection `json:"-"`
	Lamp1 bool        `json:"lamp1"`
	Lamp2 bool        `json:"lamp2"`
	Temp  string      `json:"temp"`
}
