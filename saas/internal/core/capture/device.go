package packets

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"time"
)

type Device struct {
	Name string
	Flags uint32
	Description string
	Addresses []pcap.InterfaceAddress
}

type Interface struct {
	Device   string
	Snaplen  int32
	Promisc  bool
	DevFound bool
	Err      error
	Timeout  time.Duration
	Handle   *pcap.Handle
	Buffer   gopacket.SerializeBuffer
	Options  gopacket.SerializeOptions
	Count    uint64
}

func InitTemplateInterface() Interface {
	return Interface{
		Snaplen: 1500,
		Promisc:true,
		DevFound:true,
		Err:nil,
		Timeout: 3 * time.Second,
	}
}