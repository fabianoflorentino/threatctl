package pcap

import (
	"errors"
	"os"
)

// ParsePCAP is a stub that validates the file exists and is readable.
// Real parsing will use gopacket and robust validation/sandboxing.
func ParsePCAP(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return errors.New("provided path is a directory, expected a pcap file")
	}
	// TODO: implement parsing using gopacket
	return nil
}
