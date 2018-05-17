package client

type Config struct {
	ClientID          string
	Hostname          string
	HeartbeatInterval int64
	MsgTimeout        int64
	AuthSecret        string
}
