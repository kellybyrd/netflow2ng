package formatter

import (
	"github.com/netsampler/goflow2/v2/format"
	gf2proto "github.com/netsampler/goflow2/v2/producer/proto"
	"github.com/sirupsen/logrus"
	"github.com/synfinatic/netflow2ng/proto"
	googleproto "google.golang.org/protobuf/proto"
)

const (
	PROTO_REMAPPED_IN_BYTES  = 200
	PROTO_REMAPPED_IN_PKTS   = 201
	PROTO_REMAPPED_OUT_BYTES = 202
	PROTO_REMAPPED_OUT_PKTS  = 203
)

// goflow2 allows some remapping of fields. We use this trick to save IN_BYTES/PKTS
// and OUT_BYTES/PKTS into custom fields in the protobuf message so that are not lost.
var MappingYamlStr = `
formatter:
  key:
    - sampler_address
  protobuf: # manual protobuf fields addition
    - name: in_bytes
      index: 200 # keep in sych with NFv9_REMAPPED_IN_BYTES
      type: varint
    - name: in_packets
      index: 201 # keep in sych with NFv9_REMAPPED_IN_PKTS
      type: varint
    - name: out_bytes
      index: 202 # keep in sych with NFv9_REMAPPED_OUT_BYTES
      type: varint
    - name: out_packets
      index: 203 # keep in sych with NFv9_REMAPPED_OUT_PKTS
      type: varint
# Decoder mappings
netflowv9:
  mapping:
    - field: 1
      destination: in_bytes
    - field: 2
      destination: in_packets
    - field: 23
      destination: out_bytes
    - field: 24
      destination: out_packets
`

var log *logrus.Logger //nolint:unused

func SetLogger(l *logrus.Logger) {
	log = l
}

func init() {
	format.RegisterFormatDriver("ntopjson", &NtopngJson{})
	format.RegisterFormatDriver("ntoptlv", &NtopngTlv{})
}

func castToExtendedFlowMsg(data interface{}) (*proto.ExtendedFlowMessage, error) {

	ppm, ok := data.(*gf2proto.ProtoProducerMessage)
	if !ok {
		log.Fatal("could not cast Format data to ProtoProducerMessage")
	}

	// Marshal to binary
	bin, err := googleproto.Marshal(ppm)
	if err != nil {
		log.Fatal("could not marshal ProtoProducerMessage to binary", err)
	}
	// Unmarshal into your custom struct
	efm := &proto.ExtendedFlowMessage{}
	if err := googleproto.Unmarshal(bin, efm); err != nil {
		log.Fatal("could not unmarshal binary to ExtendedFlowMsg", err)
	}

	// Need to assign the BaseFlow field explicitly, as it is not Unmarshalled automatically.
	efm.BaseFlow = &ppm.FlowMessage

	return efm, nil
}
