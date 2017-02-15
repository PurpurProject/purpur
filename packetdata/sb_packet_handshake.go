package packetdata

import "github.com/PurpurProject/elytra/packetutil"

type SBHandshakePacket struct {
	length          int32
	ProtocolVersion int32
	ServerAddress   string
	ServerPort      uint16
	NextState       int32
}

func SBHandshakePacketCreate(length int32, reader *packetutil.PacketReader) (*SBHandshakePacket, error) {
	p := new(SBHandshakePacket)
	var err error

	p.length = length

	p.ProtocolVersion, err = reader.ReadVarInt()
	if err != nil {
		return nil, err
	}

	p.ServerAddress, err = reader.ReadString()
	if err != nil {
		return nil, err
	}

	p.ServerPort, err = reader.ReadUnsignedShort()
	if err != nil {
		return nil, err
	}

	p.NextState, err = reader.ReadVarInt()
	if err != nil {
		return p, err
	}

	return p, nil
}

func (p SBHandshakePacket) Length() int32 {
	return p.length
}

func (p SBHandshakePacket) PacketID() int32 {
	return 0x00
}
