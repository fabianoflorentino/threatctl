package cmd

import (
	"encoding/json"
	"os"

	"github.com/fabianoflorentino/threatctl/internal/pcap"
	"github.com/spf13/cobra"
)

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
	rootCmd.AddCommand(summaryCmd)
	summaryCmd.Flags().IntP("top", "t", 5, "Top N items for IPs/ports")
}
