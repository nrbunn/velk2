package structs

import (
	"net"
)

const ConnNew = "CONN_NEW"
const ConnLinkDead = "CONN_LINK_DEAD"
const ReadBuffer = 512

type Connection struct {
	conn       net.Conn
	Status     string
	readBuffer int
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		conn:       conn,
		Status:     ConnNew,
		readBuffer: ReadBuffer,
	}
}

func (c *Connection) Write(s string) error {
	if c.Status == ConnLinkDead {
		return nil
	}
	_, err := c.conn.Write([]byte(s))
	if err != nil {
		c.Status = ConnLinkDead
		c.Close()
		return err
	}
	return nil
}

func (c *Connection) Read() (string, error) {
	for {
		buf := make([]byte, c.readBuffer)
		nr, err := c.conn.Read(buf)
		if err != nil {
			c.Status = ConnLinkDead
			c.Close()
			return "", err
		}

		data := buf[0 : nr-1]
		return string(data), nil
	}
}

func (c *Connection) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}
