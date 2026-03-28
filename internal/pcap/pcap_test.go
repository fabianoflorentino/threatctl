package pcap

import (
	"net"
	"os"
	"testing"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

func TestParsePCAP_Nonexistent(t *testing.T) {
	if err := ParsePCAP("nonexistent.pcap"); err == nil {
		t.Fatalf("expected error for nonexistent file")
	}
}

func TestSummarizePCAP_Sample(t *testing.T) {
	sum, err := SummarizePCAP("../samples/example.pcap", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum.TotalPackets != 2 {
		t.Fatalf("expected 2 packets, got %d", sum.TotalPackets)
	}
	if sum.FlowCount == 0 {
		t.Fatalf("expected flow count > 0")
	}
}

func TestParsePCAP_Valid(t *testing.T) {
	if err := ParsePCAP("../samples/example.pcap"); err != nil {
		t.Fatalf("expected no error for valid pcap: %v", err)
	}
}

func TestSummarizePCAP_TopNZero(t *testing.T) {
	sum, err := SummarizePCAP("../samples/example.pcap", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sum.TopSrcIPs) != 0 || len(sum.TopDstIPs) != 0 {
		t.Fatalf("expected no top lists when topN=0")
	}
}

func TestSummarizePCAP_IPv6UDP(t *testing.T) {
	// create a temporary pcap with an IPv6 UDP packet and ensure summarize reads it
	f, err := os.CreateTemp("", "sample-*.pcap")
	if err != nil {
		t.Fatalf("tmp file: %v", err)
	}
	name := f.Name()
	defer os.Remove(name)

	w := pcapgo.NewWriter(f)
	if err := w.WriteFileHeader(65536, layers.LinkTypeEthernet); err != nil {
		t.Fatalf("write header: %v", err)
	}

	eth := &layers.Ethernet{SrcMAC: []byte{0x02, 0, 0, 0, 0, 1}, DstMAC: []byte{0x02, 0, 0, 0, 0, 2}, EthernetType: layers.EthernetTypeIPv6}
	ip6 := &layers.IPv6{Version: 6, HopLimit: 64, SrcIP: net.ParseIP("fe80::1"), DstIP: net.ParseIP("ff02::2"), NextHeader: layers.IPProtocolUDP}
	udp := &layers.UDP{SrcPort: 40000, DstPort: 53}
	udp.SetNetworkLayerForChecksum(ip6)

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{ComputeChecksums: true, FixLengths: true}
	if err := gopacket.SerializeLayers(buf, opts, eth, ip6, udp); err != nil {
		t.Fatalf("serialize: %v", err)
	}
	ci := gopacket.CaptureInfo{Timestamp: time.Now(), CaptureLength: len(buf.Bytes()), Length: len(buf.Bytes())}
	if err := w.WritePacket(ci, buf.Bytes()); err != nil {
		t.Fatalf("write pkt: %v", err)
	}
	f.Close()

	sum, err := SummarizePCAP(name, 2)
	if err != nil {
		t.Fatalf("summarize: %v", err)
	}
	if sum.ProtocolCounts["udp"] == 0 {
		t.Fatalf("expected udp count > 0")
	}
	if sum.TotalPackets == 0 {
		t.Fatalf("expected total packets > 0")
	}
}
