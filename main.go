package main

import (
	"net"
	"os"

	"time"

	"github.com/PurpurProject/purpur/server"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("purpur")

var format = logging.MustStringFormatter(
	"%{color}[%{time:2006-01-02 15:04:05.000}] [%{level}] [%{shortfunc}]%{color:reset} %{message}",
)

var purpurVersion = "0.0.1"

func main() {

	/* STARTUP */

	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend1Format := logging.NewBackendFormatter(backend1, format)

	backend1LF := logging.AddModuleLevel(backend1Format)

	var loggingLevel logging.Level

	loggingLevel = logging.DEBUG

	backend1LF.SetLevel(loggingLevel, "")

	logging.SetBackend(backend1LF)

	server.Log = log
	server.CurrentStatus = server.CreateStatusObject()

	log.Info("Purpur v.", purpurVersion, " starting up...")

	srv, err := net.Listen("tcp", ":25565")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	for {
		conn, err := srv.Accept()
		if err != nil {
			log.Error(err.Error())
			continue
		}
		conn.SetDeadline(time.Now().Add(time.Second * 5))
		go server.HandleConnection(server.SpawnClientConnection(conn, server.HANDSHAKE))
	}
}
