package packets

import (
	"github.com/google/gopacket"
	"time"
)

type LayerStruct struct {
	LayerContents []byte
	LayerPayload  []byte
	LayerType     gopacket.LayerType
	LayerFlow     gopacket.Flow
}

type AppLayerStruct struct {
	LayerContents []byte
	LayerPayload  []byte
	LayerFlow     []byte
	LayerType     gopacket.LayerType
}

type Metadata struct {
	CaptureLength  int
	InterfaceIndex int
	Length         int
	Timestamp      time.Time
	CaptureInfo    gopacket.CaptureInfo
	AncillaryData  []interface{}
	Truncated      bool
}

type PacketStruct struct {
	Layers           []gopacket.Layer
	Data             []byte
	Dump             string
	Metadata         Metadata
	ApplicationLayer AppLayerStruct
	TransportLayer   LayerStruct
	NetworkLayer     LayerStruct
	LinkLayer        LayerStruct
}

