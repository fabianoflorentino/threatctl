package pcap

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

// Summary holds basic metadata counts extracted from a PCAP file.
type Summary struct {
	TotalPackets   int            `json:"total_packets"`
	ProtocolCounts map[string]int `json:"protocol_counts"`
	SrcIPs         map[string]int `json:"src_ips"`
	DstIPs         map[string]int `json:"dst_ips"`
	SrcPorts       map[uint16]int `json:"src_ports"`
	DstPorts       map[uint16]int `json:"dst_ports"`
	Earliest       time.Time      `json:"earliest_timestamp"`
	Latest         time.Time      `json:"latest_timestamp"`
	FlowCount      int            `json:"flow_count"`
	TopSrcIPs      []KeyCount     `json:"top_src_ips,omitempty"`
	TopDstIPs      []KeyCount     `json:"top_dst_ips,omitempty"`
	TopSrcPorts    []KeyCount     `json:"top_src_ports,omitempty"`
	TopDstPorts    []KeyCount     `json:"top_dst_ports,omitempty"`
}

// KeyCount represents a key and its count for top-N lists.
type KeyCount struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

// ParsePCAP is a minimal validator that checks existence and readability.
func ParsePCAP(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := pcapgo.NewReader(f); err != nil {
		return err
	}
	return nil
}

// SummarizePCAP reads the pcap file and returns a Summary with basic counts.
// topN controls how many top items are included for IPs/ports.
func SummarizePCAP(path string, topN int) (*Summary, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader, err := pcapgo.NewReader(f)
	if err != nil {
		return nil, err
	}

	src := gopacket.NewPacketSource(reader, reader.LinkType())
	summary := &Summary{
		ProtocolCounts: make(map[string]int),
		SrcIPs:         make(map[string]int),
		DstIPs:         make(map[string]int),
		SrcPorts:       make(map[uint16]int),
		DstPorts:       make(map[uint16]int),
	}

	flowMap := make(map[string]struct{})
	for packet := range src.Packets() {
		summary.TotalPackets++

		ci := packet.Metadata().Timestamp
		if summary.Earliest.IsZero() || ci.Before(summary.Earliest) {
			summary.Earliest = ci
		}
		if summary.Latest.IsZero() || ci.After(summary.Latest) {
			summary.Latest = ci
		}

		// Network layer
		if netLayer := packet.NetworkLayer(); netLayer != nil {
			switch nl := netLayer.(type) {
			case *layers.IPv4:
				summary.SrcIPs[nl.SrcIP.String()]++
				summary.DstIPs[nl.DstIP.String()]++
				summary.ProtocolCounts["ipv4"]++
			case *layers.IPv6:
				summary.SrcIPs[nl.SrcIP.String()]++
				summary.DstIPs[nl.DstIP.String()]++
				summary.ProtocolCounts["ipv6"]++
			default:
				summary.ProtocolCounts[fmt.Sprintf("net-%T", nl)]++
			}
		}

		// Transport layer: TCP/UDP
		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)
			summary.ProtocolCounts["tcp"]++
			summary.SrcPorts[uint16(tcp.SrcPort)]++
			summary.DstPorts[uint16(tcp.DstPort)]++
			if packet.NetworkLayer() != nil {
				srcIP := packet.NetworkLayer().NetworkFlow().Src().String()
				dstIP := packet.NetworkLayer().NetworkFlow().Dst().String()
				proto := "tcp"
				key := fmt.Sprintf("%s:%d-%s:%d-%s", srcIP, tcp.SrcPort, dstIP, tcp.DstPort, proto)
				flowMap[key] = struct{}{}
			}
		}
		if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
			udp, _ := udpLayer.(*layers.UDP)
			summary.ProtocolCounts["udp"]++
			summary.SrcPorts[uint16(udp.SrcPort)]++
			summary.DstPorts[uint16(udp.DstPort)]++
			if packet.NetworkLayer() != nil {
				srcIP := packet.NetworkLayer().NetworkFlow().Src().String()
				dstIP := packet.NetworkLayer().NetworkFlow().Dst().String()
				proto := "udp"
				key := fmt.Sprintf("%s:%d-%s:%d-%s", srcIP, udp.SrcPort, dstIP, udp.DstPort, proto)
				flowMap[key] = struct{}{}
			}
		}

		// Application protocols (simple detection)
		if packet.Layer(layers.LayerTypeDNS) != nil {
			summary.ProtocolCounts["dns"]++
		}
		// Note: HTTP/TLS/SNI detection would require deeper parsing or reassembly.
	}

	summary.FlowCount = len(flowMap)

	// compute top-N lists
	summary.TopSrcIPs = topNFromStringMap(summary.SrcIPs, topN)
	summary.TopDstIPs = topNFromStringMap(summary.DstIPs, topN)
	sp := make(map[string]int)
	dp := make(map[string]int)
	for k, v := range summary.SrcPorts {
		sp[fmt.Sprintf("%d", k)] = v
	}
	for k, v := range summary.DstPorts {
		dp[fmt.Sprintf("%d", k)] = v
	}
	summary.TopSrcPorts = topNFromStringMap(sp, topN)
	summary.TopDstPorts = topNFromStringMap(dp, topN)

	return summary, nil
}

// topNFromStringMap returns the top N entries from a map[string]int sorted by count desc.
func topNFromStringMap(m map[string]int, n int) []KeyCount {
	if n <= 0 {
		return nil
	}
	items := make([]KeyCount, 0, len(m))
	for k, v := range m {
		items = append(items, KeyCount{Key: k, Count: v})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Count > items[j].Count })
	if len(items) > n {
		items = items[:n]
	}
	return items
}
