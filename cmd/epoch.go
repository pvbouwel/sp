/*
Copyright © 2025 Peter Van Bouwel <https://github.com/pvbouwel>
*/
package cmd

import (
	"os"

	"github.com/pvbouwel/sp/epoch"
	"github.com/spf13/cobra"
)

// epochCmd represents the epoch command
var epochCmd = &cobra.Command{
	Use:   "epoch",
	Short: "Replace epoch occurrences",
	Long:  `Replace all epoch occurences in the input`,
	Run: func(cmd *cobra.Command, args []string) {
		stdoutWriter = epoch.NewEpoch(os.Stdout)
		stderrWriter = epoch.NewEpoch(os.Stderr)
	},
}

func init() {
	rootCmd.AddCommand(epochCmd)
}
