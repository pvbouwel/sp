/*
Copyright © 2025 Peter Van Bouwel <https://github.com/pvbouwel>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func getShellRc() string {
	shellVar := os.Getenv("SHELL")
	shellBin := filepath.Base(shellVar)
	if shellBin == "." {
		return "/replace/with/path/to/your/rc/file"
	}
	return fmt.Sprintf("${HOME}/.%src", shellBin)
}

// colorCmd represents the color command
var aliasCmd = &cobra.Command{
	Use:   "aliases",
	Short: "Show shell aliases",
	Long:  `Show example alias commands to install in your shell-rc file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("# ===START OUTPUT sp aliases===")
		fmt.Println("# This output can be sourced from your shell-rc file:")
		fmt.Println("#   sp aliases > \"$HOME/.sp-aliases\"")
		fmt.Printf("#   echo 'source \"$HOME/.sp-aliases\"' >> %s\n\n", getShellRc())

		fmt.Println(`# Rainbow colours
alias sp-rainbow="sp color --color-type rotating --rotating-type random --stride-length 15-25"

# Colour JSON depending on values of the field called levelname and have alternating colours if subsequent lines match
alias sp-json-traffic-levelname='sp color --ignore-case --color-type JSON --json-key levelname --colors INFO.0.255.0,INFO.0.155.0,WARNING.255.128.0,WARNING.155.128.0,ERROR.255.0.0,ERROR.155.0.0'
alias sp-json-traffic-level='sp color --ignore-case --color-type JSON --json-key levelname --colors INFO.0.255.0,INFO.0.155.0,WARNING.255.128.0,WARNING.155.128.0,ERROR.255.0.0,ERROR.155.0.0'

# Allow colouring of stoud and stderr differently
alias sp-stdouterr="sp color"

# Replace epoch occurrences with human readable time
alias sp-epoch="sp epoch"`)

		fmt.Println("# ===END OUTPUT sp aliases===")

		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
}
