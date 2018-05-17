package client

type Consumer struct {
	id          int64
	topic       string
	channel     string
	connections map[string]*Conn
}
