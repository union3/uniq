package client

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type IdentifyResponse struct {
	AuthRequired bool `json:"auth_required"`
}

type AuthResponse struct {
	Identity string `json:"identity"`
}

type msgResponse struct {
	msg     *Message
	cmd     *Command
	success bool
	backoff bool
}

type Conn struct {
	addr string
	conn *net.TCPConn
	r    io.Reader
	w    io.Writer
}

func NewConn(addr string) (conn *Conn) {
	conn = &Conn{
		addr: addr,
	}
	return conn
}

func (c *Conn) Connect() (*IdentifyResponse, error) {
	dialer := &net.Dialer{
		LocalAddr: "127.0.0.1",
		Timeout:   tine.Second * 3,
	}
	conn, err := dialer.Dial("tcp", c.addr)
	if err != nil {
		return nil, err
	}
	c.conn = conn.(*net.TCPConn)
	c.r = conn
	c.w = conn
	_, err = c.Write(MagicV2)
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("[%s] failed to write magic - %s", c.addr, err)
	}
	resp, err := c.identify()
	if err != nil {
		return nil, err
	}
	if resp != nil && resp.AuthRequired {
		if c.config.AuthSecret == "" {
			return nil, errors.New("Auth Required")
		}
		err := c.auth(c.config.AuthSecret)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}

func (c *Conn) identify() (*IdentifyResponse, error) {
	ci := make(map[string]interface{})
	ci["client_id"] = c.config.ClientID
	ci["hostname"] = c.config.Hostname
	ci["user_agent"] = c.config.UserAgent
	ci["short_id"] = c.config.ClientID // deprecated
	ci["long_id"] = c.config.Hostname  // deprecated
	ci["tls_v1"] = c.config.TlsV1
	ci["deflate"] = c.config.Deflate
	ci["deflate_level"] = c.config.DeflateLevel
	ci["snappy"] = c.config.Snappy
	ci["feature_negotiation"] = true
	if c.config.HeartbeatInterval == -1 {
		ci["heartbeat_interval"] = -1
	} else {
		ci["heartbeat_interval"] = int64(c.config.HeartbeatInterval / time.Millisecond)
	}
	ci["sample_rate"] = c.config.SampleRate
	ci["output_buffer_size"] = c.config.OutputBufferSize
	if c.config.OutputBufferTimeout == -1 {
		ci["output_buffer_timeout"] = -1
	} else {
		ci["output_buffer_timeout"] = int64(c.config.OutputBufferTimeout / time.Millisecond)
	}
	ci["msg_timeout"] = int64(c.config.MsgTimeout / time.Millisecond)
	cmd, err := Identify(ci)
	if err != nil {
		return nil, ErrIdentify{err.Error()}
	}
	err = c.WriteCommand(cmd)
	if err != nil {
		return nil, ErrIdentify{err.Error()}
	}
	frameType, data, err := ReadUnpackedResponse(c)
	if err != nil {
		return nil, ErrIdentify{err.Error()}
	}
	if frameType == FrameTypeError {
		return nil, ErrIdentify{string(data)}
	}
	// check to see if the server was able to respond w/ capabilities
	// i.e. it was a JSON response
	if data[0] != '{' {
		return nil, nil
	}
	resp := &IdentifyResponse{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, ErrIdentify{err.Error()}
	}
	c.log(LogLevelDebug, "IDENTIFY response: %+v", resp)
	// now that connection is bootstrapped, enable read buffering
	// (and write buffering if it's not already capable of Flush())
	c.r = bufio.NewReader(c.r)
	if _, ok := c.w.(flusher); !ok {
		c.w = bufio.NewWriter(c.w)
	}
	return resp, nil
}

func (c *Conn) auth(secret string) error {
	cmd, err := Auth(secret)
	if err != nil {
		return err
	}
	err = c.WriteCommand(cmd)
	if err != nil {
		return err
	}
	frameType, data, err := ReadUnpackedResponse(c)
	if err != nil {
		return err
	}
	if frameType == FrameTypeError {
		return errors.New("Error authenticating " + string(data))
	}
	resp := &AuthResponse{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return err
	}
	c.log(LogLevelInfo, "Auth accepted. Identity: %q %s Permissions: %d",
		resp.Identity, resp.IdentityUrl, resp.PermissionCount)
	return nil
}

type flusher interface {
	Flush() error
}

func (c *Conn) Flush() error {
	if f, ok := c.w.(flusher); ok {
		return f.Flush()
	}
	return nil
}

func (c *Conn) WriteCommand(cmd *Command) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	_, err := cmd.WriteTo(c)
	if err != nil {
		return err
	}
	err = c.Flush()
	if err != nil {
		return err
	}
	return nil
}
