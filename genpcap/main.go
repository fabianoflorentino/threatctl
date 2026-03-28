package main

import (
	"log"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

func main() {
	if err := os.MkdirAll("samples", 0o755); err != nil {
		log.Fatalf("mkdir: %v", err)
	}
	f, err := os.Create("internal/samples/example.pcap")
	if err != nil {
		log.Fatalf("create pcap: %v", err)
	}
	defer f.Close()

	w := pcapgo.NewWriter(f)
	w.WriteFileHeader(65536, layers.LinkTypeEthernet)

	// Build a DNS packet (UDP)
	eth := &layers.Ethernet{SrcMAC: []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x01}, DstMAC: []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x02}, EthernetType: layers.EthernetTypeIPv4}
	ip := &layers.IPv4{Version: 4, IHL: 5, TTL: 64, SrcIP: []byte{192, 168, 0, 2}, DstIP: []byte{8, 8, 8, 8}, Protocol: layers.IPProtocolUDP}
	udp := &layers.UDP{SrcPort: 12345, DstPort: 53}
	dns := &layers.DNS{ID: 0x1a2b, QR: false, RD: true, Questions: []layers.DNSQuestion{{Name: []byte("example.com"), Type: layers.DNSTypeA, Class: layers.DNSClassIN}}}
	udp.SetNetworkLayerForChecksum(ip)

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{ComputeChecksums: true, FixLengths: true}
	if err := gopacket.SerializeLayers(buf, opts, eth, ip, udp, dns); err != nil {
		log.Fatalf("serialize dns: %v", err)
	}
	w.WritePacket(gopacket.CaptureInfo{Timestamp: time.Now(), CaptureLength: len(buf.Bytes()), Length: len(buf.Bytes())}, buf.Bytes())

	// Build a simple TCP HTTP GET packet (no payload reassembly)
	eth2 := *eth
	ip2 := *ip
	ip2.DstIP = []byte{93, 184, 216, 34}
	tcp := &layers.TCP{SrcPort: 55555, DstPort: 80, Seq: 110, SYN: true, Window: 14600}
	tcp.SetNetworkLayerForChecksum(&ip2)
	buf2 := gopacket.NewSerializeBuffer()
	if err := gopacket.SerializeLayers(buf2, opts, &eth2, &ip2, tcp); err != nil {
		log.Fatalf("serialize tcp: %v", err)
	}
	w.WritePacket(gopacket.CaptureInfo{Timestamp: time.Now(), CaptureLength: len(buf2.Bytes()), Length: len(buf2.Bytes())}, buf2.Bytes())

	log.Println("wrote internal/samples/example.pcap")
}
