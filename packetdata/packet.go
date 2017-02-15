package packetdata

type UncompressedPacket interface {
	Length() int32
	PacketID() int32
}
