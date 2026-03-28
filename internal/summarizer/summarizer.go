package summarizer

import "github.com/fabianoflorentino/threatctl/internal/pcap"

// Summarizer defines the core port for summarizing packet captures.
type Summarizer interface {
	Summarize(path string, topN int) (*pcap.Summary, error)
}

// PCAPSummarizer is an adapter that implements Summarizer using internal/pcap.
type PCAPSummarizer struct{}

func NewPCAPSummarizer() *PCAPSummarizer { return &PCAPSummarizer{} }

func (p *PCAPSummarizer) Summarize(path string, topN int) (*pcap.Summary, error) {
	return pcap.SummarizePCAP(path, topN)
}
