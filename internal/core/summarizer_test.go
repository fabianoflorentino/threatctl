package core

import "testing"

func TestPCAPSummarizer_Summarize(t *testing.T) {
	s := NewPCAPSummarizer()
	sum, err := s.Summarize("../../samples/example.pcap", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum.TotalPackets == 0 {
		t.Fatalf("expected packets > 0")
	}
	if sum.FlowCount == 0 {
		t.Fatalf("expected flow count > 0")
	}
	if len(sum.SrcIPs) == 0 || len(sum.DstIPs) == 0 {
		t.Fatalf("expected IP maps to be populated")
	}
}
