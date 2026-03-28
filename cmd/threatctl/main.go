package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fabianoflorentino/threatctl/internal/pcap"
	"github.com/fabianoflorentino/threatctl/pkg/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "threatctl",
	Short: "ThreatCTL — Low-level Network Forensics Toolkit",
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze [pcap file]",
	Short: "Analyze a PCAP file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		fmt.Printf("threatctl v%s — analyzing %s\n", version.Version, path)
		if err := pcap.ParsePCAP(path); err != nil {
			return err
		}
		fmt.Println("analysis complete (stub)")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(summaryCmd)
}

var summaryCmd = &cobra.Command{
	Use:   "summary [pcap file]",
	Short: "Print a JSON summary of a PCAP file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		topN, _ := cmd.Flags().GetInt("top")
		s, err := pcap.SummarizePCAP(path, topN)
		if err != nil {
			return err
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(s)
	},
}

func init() {
	summaryCmd.Flags().IntP("top", "t", 5, "Top N items for IPs/ports")
}
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
