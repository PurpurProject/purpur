package server

import (
	"net"

	"io/ioutil"

	"encoding/binary"
	"encoding/json"

	"github.com/CalmBit/elytra/connutil"
	"github.com/CalmBit/elytra/packetutil"
	"github.com/CalmBit/purpur/packetdata"
	logging "github.com/op/go-logging"
)

var ClientConnectionMap map[string]*ClientConnection

var Log *logging.Logger

var CurrentStatus *ServerStatus

const minecraftVersionName = "1.11.2"
const minecraftProtocolVersion = 316

func HandleConnection(connection *ClientConnection) {

	Log.Debug("Connection handler initiated")

	for !connection.isClosed {

		packet, packetSize, packetID, err := readPacketHeader(connection)

		if err != nil {
			DestroyClientConnection(connection)
			Log.Error("Connection Terminated On Read: " + err.Error())
			return
		}

		Log.Debug("Packet Size: ", packetSize)
		Log.Debug("Packet ID: ", packetID)
		Log.Debugf("Packet Contents: %v\n", packet)

		reader := packetutil.CreatePacketReader(packet)

		switch connection.state {
		case HANDSHAKE:
			{
				switch packetID {
				case 0x00:
					{
						pkt, err := packetdata.SBHandshakePacketCreate(packetSize, reader)
						if err != nil {
							DestroyClientConnection(connection)
							Log.Error(err.Error())
						}

						connection.KeepAlive()
						connection.state = int(pkt.NextState)
						break
					}
				case 0xFE:
					{
						Log.Warning("Legacy Ping Request received - this may mean the main status isn't working!")
						DestroyClientConnection(connection)
						return
					}
				}
				break
			}
		case STATUS:
			{
				switch packetID {
				case 0x00:
					{
						connection.KeepAlive()

						writer := packetutil.CreatePacketWriter(0x00)

						marshaledStatus, err := json.Marshal(*CurrentStatus)
						if err != nil {
							Log.Error(err.Error())
							DestroyClientConnection(connection)
							return
						}
						writer.WriteString(string(marshaledStatus))

						SendData(connection, writer)
						break
					}
				case 0x01:
					{
						writer := packetutil.CreatePacketWriter(0x01)
						mirror, _ := reader.ReadLong()
						writer.WriteLong(mirror)
						SendData(connection, writer)
						DestroyClientConnection(connection)
					}
				}
			}
		}
	}
}

func SendData(connection *ClientConnection, writer *packetutil.PacketWriter) {
	connection.conn.Write(writer.GetPacket())
}

func getPacketData(conn net.Conn) ([]byte, error) {
	return ioutil.ReadAll(conn)
}

func readPacketHeader(conn *ClientConnection) ([]byte, int32, int32, error) {

	packetSize, err := connutil.ParseVarIntFromConnection(conn.conn)

	if err != nil {
		return nil, 0, 0, err
	}

	if packetSize == 254 && conn.state == HANDSHAKE {
		preBufferSize := 29
		preBuffer := make([]byte, preBufferSize)
		conn.conn.Read(preBuffer)
		postBufferSize := int(binary.BigEndian.Uint16(preBuffer[25:]))
		postBuffer := make([]byte, postBufferSize)
		size := preBufferSize + postBufferSize
		conn.conn.Read(postBuffer)
		return append(preBuffer, postBuffer...), int32(size), 0xFE, nil
	}

	packetID, err := connutil.ParseVarIntFromConnection(conn.conn)

	if err != nil {
		return nil, 0, 0, err
	}

	if packetSize-1 == 0 {
		return nil, packetSize, packetID, nil
	}
	packet := make([]byte, packetSize-1)
	conn.conn.Read(packet)

	return packet, packetSize - 1, packetID, nil
}
