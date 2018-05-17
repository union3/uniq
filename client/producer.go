package client

type Producer struct {
	id          int64
	topic       string
	channel     string
	connections map[string]*Conn
}
