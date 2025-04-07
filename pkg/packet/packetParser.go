package packet

import (
	"encoding/binary"
	"fmt"
	"net"
)

type PacketInfo struct {
	SrcIP   string `json:"src_ip"`
	DstIP   string `json:"dst_ip"`
	SrcPort uint16 `json:"src_port,omitempty"`
	DstPort uint16 `json:"dst_port,omitempty"`
	Proto   string `json:"protocol"`
}

func ParsePacket(data []byte) (PacketInfo, error) {
	if len(data) < 20 {
		return PacketInfo{}, fmt.Errorf("too short packet")
	}

	version := data[0] >> 4
	if version != 4 {
		return PacketInfo{}, fmt.Errorf("only IPv4 supported")
	}

	srcIP := net.IP(data[12:16]).String()
	dstIP := net.IP(data[16:20]).String()
	proto := data[9]

	var info PacketInfo
	info.SrcIP = srcIP
	info.DstIP = dstIP

	switch proto {
	case 6: // TCP
		if len(data) < 34 {
			return PacketInfo{}, fmt.Errorf("malformed TCP packet")
		}
		info.SrcPort = binary.BigEndian.Uint16(data[20:22])
		info.DstPort = binary.BigEndian.Uint16(data[22:24])
		info.Proto = "TCP"
	case 17: // UDP
		if len(data) < 28 {
			return PacketInfo{}, fmt.Errorf("malformed UDP packet")
		}
		info.SrcPort = binary.BigEndian.Uint16(data[20:22])
		info.DstPort = binary.BigEndian.Uint16(data[22:24])
		info.Proto = "UDP"
	default:
		return PacketInfo{}, fmt.Errorf("unsupported protocol: %d", proto)
	}

	return info, nil
}
