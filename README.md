# ThreatCTL

ThreatCTL is a lightweight, low-level network forensics toolkit focused on offline PCAP analysis. It provides command-line tools to analyze packet captures, extract indicators, and generate concise summaries for incident response and research.

Status: experimental (MVP)

Features

- Analyze PCAP files for suspicious activity and artifacts
- Extract indicators (IPs, domains, file hashes) from captures
- Produce human-readable summaries for quick triage
- Small, dependency-light Go codebase suitable for local analysis

Prerequisites

- Go 1.18+ installed

Build
Clone the repository and build the main binary and helper tools:

```bash
git clone https://github.com/your-org/threatctl.git
cd threatctl
go build -o bin/threatctl ./
go build -o bin/genpcap ./genpcap
```

Usage
Basic analyze command (example):

```bash
./bin/threatctl analyze samples/example.pcap
```

Other commands

- `analyze` — run detailed analysis on a PCAP file
- `summary` — produce a short summary report

Run tests

```bash
go test ./...
```

Repository layout

- `cmd/` — CLI entry points (`analyze`, `summary`, root command)
- `internal/` — core analysis, pcap parsing, and summarizer logic
- `genpcap/` — helper tool for generating sample PCAPs
- `pkg/version` — version info used by the CLI
- `samples/` — example PCAPs

Contributing
See CONTRIBUTING.md for guidelines. Please open issues for feature requests or bugs.

License
This project is licensed under the terms in the LICENSE file.

Contact
Create issues in this repository for questions or support.
