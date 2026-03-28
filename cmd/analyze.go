package cmd

import (
	"fmt"

	"github.com/fabianoflorentino/threatctl/internal/pcap"
	"github.com/fabianoflorentino/threatctl/pkg/version"
	"github.com/spf13/cobra"
)

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
}
