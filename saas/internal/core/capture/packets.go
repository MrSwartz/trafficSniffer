package packets

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
	e "saas/internal/debug/err"
	"strconv"
)

// buffered read
func Capture(iface Interface, filter string) []PacketStruct {
	buf := make([]PacketStruct, 10000)
	
	handle, err := pcap.OpenLive(iface.Device, iface.Snaplen, iface.Promisc, iface.Timeout)
	if err != nil {
		log.Println(err)
	}
	defer handle.Close()
	if err := handle.SetBPFFilter(filter); err != nil {
		log.Println(err)
	}
	source := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range source.Packets() {
		fmt.Println(retPS(packet))
		buf = append(buf, retPS(packet))
	}
	return buf
}

func retPS(p gopacket.Packet) (ps PacketStruct) {
	if p.Metadata() != nil {
		var m Metadata
		m.CaptureInfo = p.Metadata().CaptureInfo
		m.AncillaryData = p.Metadata().AncillaryData
		m.CaptureLength = p.Metadata().CaptureLength
		m.InterfaceIndex = p.Metadata().InterfaceIndex
		m.Length = p.Metadata().Length
		m.Timestamp = p.Metadata().Timestamp
		m.Truncated = p.Metadata().Truncated
		ps.Metadata = m
	}
	if p.NetworkLayer() != nil {
		var ls LayerStruct
		ls.LayerPayload = p.NetworkLayer().LayerPayload()
		ls.LayerType = p.NetworkLayer().LayerType()
		ls.LayerContents = p.NetworkLayer().LayerContents()
		ls.LayerFlow = p.NetworkLayer().NetworkFlow()
		ps.NetworkLayer = ls
	}
	if p.TransportLayer() != nil {
		var tl LayerStruct
		tl.LayerPayload = p.TransportLayer().LayerPayload()
		tl.LayerType = p.TransportLayer().LayerType()
		tl.LayerContents = p.TransportLayer().LayerContents()
		tl.LayerFlow = p.TransportLayer().TransportFlow()
		ps.TransportLayer = tl
	}
	if p.LinkLayer() != nil {
		var ll LayerStruct
		ll.LayerPayload = p.LinkLayer().LayerPayload()
		ll.LayerType = p.LinkLayer().LayerType()
		ll.LayerContents = p.LinkLayer().LayerContents()
		ll.LayerFlow = p.LinkLayer().LinkFlow()
		ps.NetworkLayer = ll
	}
	if p.ApplicationLayer() != nil {
		var al AppLayerStruct
		al.LayerPayload = p.ApplicationLayer().LayerPayload()
		al.LayerType = p.ApplicationLayer().LayerType()
		al.LayerContents = p.ApplicationLayer().LayerContents()
		al.LayerFlow = p.ApplicationLayer().Payload()
	}
	ps.Data = p.Data()
	ps.Dump = p.Dump()
	ps.Layers = p.Layers()
	return
}

func (i Interface)Send(srcIP, dstIP, srcPort, dstPort, srcMAC, dstMAC string) {
	handle, err := pcap.OpenLive(i.Device, i.Snaplen, i.Promisc, i.Timeout)
	e.FatalErr("", err)

	rBytes := []byte{10, 20, 30}
	err = handle.WritePacketData(rBytes)
	e.FatalErr("Error writing bytes to network device. ", err)

	buffer := gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(buffer, i.Options, &layers.Ethernet{}, &layers.IPv4{}, &layers.TCP{}, gopacket.Payload(rBytes))

	out := buffer.Bytes()
	err = handle.WritePacketData(out)
	e.FatalErr("Error sending packet to network service", err)

	sm, err := net.ParseMAC(srcMAC)
	e.CheckErr(err)

	dm, err := net.ParseMAC(dstMAC)
	e.CheckErr(err)

	ethl := &layers.Ethernet{SrcMAC: sm, DstMAC: dm}

	ipl := &layers.IPv4{SrcIP: net.ParseIP(srcIP), DstIP: net.ParseIP(dstIP)}

	sp, err := strconv.ParseInt(srcPort, 10, 64)
	e.CheckErr(err)

	dp, err := strconv.ParseInt(dstPort, 10, 64)
	e.CheckErr(err)

	tcpl := &layers.TCP{SrcPort: layers.TCPPort(sp), DstPort: layers.TCPPort(dp)}

	buffer = gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(buffer, i.Options, ethl, ipl, tcpl, gopacket.Payload(rBytes))
	out = buffer.Bytes()

	handle.Close()
}
