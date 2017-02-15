package server

import (
	"net"
	"strings"
	"time"
)

type ClientConnection struct {
	conn     net.Conn
	state    int
	isClosed bool
}

func SpawnClientConnection(conn net.Conn, state int) *ClientConnection {
	if ClientConnectionMap == nil {
		ClientConnectionMap = make(map[string]*ClientConnection)
	}

	remoteAddr := strings.Split(conn.RemoteAddr().String(), ":")[0]

	_, connectionExists := ClientConnectionMap[remoteAddr]

	if connectionExists {
		Log.Debug("Client connection resurrected")
		tmp := ClientConnectionMap[remoteAddr]
		tmp.conn.Close()
		tmp.conn = conn
		ClientConnectionMap[remoteAddr] = tmp

	} else {
		Log.Debug("Client connection spawned")
		connection := new(ClientConnection)
		connection.conn = conn
		connection.state = state
		ClientConnectionMap[remoteAddr] = connection
	}

	return ClientConnectionMap[remoteAddr]
}

func DestroyClientConnection(connection *ClientConnection) {
	if ClientConnectionMap == nil {
		ClientConnectionMap = make(map[string]*ClientConnection)
	}

	Log.Debug("Client connection destroyed")

	remoteAddr := strings.Split(connection.conn.RemoteAddr().String(), ":")[0]

	connection.conn.Close()

	connection.isClosed = true

	delete(ClientConnectionMap, remoteAddr)
}

func (c *ClientConnection) KeepAlive() {
	c.conn.SetDeadline(time.Now().Add(time.Second * 5))
}
