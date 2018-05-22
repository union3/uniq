package client

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

var byteSpace = []byte(" ")
var byteNewLine = []byte("\n")

type Command struct {
	Name   []byte
	Params [][]byte
	Body   []byte
}

func (c *Command) String() string {
	if len(c.Params) > 0 {
		return fmt.Sprintf("%s %s", c.Name, string(bytes.Join(c.Params, byteSpace)))
	}
	return string(c.Name)
}

func (c *Command) WriteTo(w io.Writer) (int64, error) {
	var total int64
	var buf [4]byte
	n, err := w.Write(c.Name)
	total += int64(n)
	if err != nil {
		return total, err
	}
	for _, param := range c.Params {
		n, err := w.Write(byteSpace)
		total += int64(n)
		if err != nil {
			return total, err
		}
		n, err = w.Write(param)
		total += int64(n)
		if err != nil {
			return total, err
		}
	}
	n, err = w.Write(byteNewLine)
	total += int64(n)
	if err != nil {
		return total, err
	}
	if c.Body != nil {
		bufs := buf[:]
		binary.BigEndian.PutUint32(bufs, uint32(len(c.Body)))
		n, err := w.Write(bufs)
		total += int64(n)
		if err != nil {
			return total, err
		}
		n, err = w.Write(c.Body)
		total += int64(n)
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

func Ping() *Command {
	return &Command{[]byte("PING"), nil, nil}
}

func Identify(js map[string]interface{}) (*Command, error) {
	body, err := json.Marshal(js)
	if err != nil {
		return nil, err
	}
	return &Command{[]byte("IDENTIFY"), nil, body}, nil
}

func Auth(secret string) (*Command, error) {
	return &Command{[]byte("AUTH"), nil, []byte(secret)}, nil
}
