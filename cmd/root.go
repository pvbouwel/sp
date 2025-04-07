/*
Copyright © 2025 Peter Van Bouwel <https://github.com/pvbouwel>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/pvbouwel/sp/streams"
	"github.com/spf13/cobra"
)


var appName string
var appArgs []string
const appSeparator string = "--"
var stdoutWriter io.Writer
var stderrWriter io.Writer

func isAppSepartor(s string) bool {
	return s == appSeparator
}

func proccessAppArguments() {
	separatorIdx := slices.IndexFunc(os.Args, isAppSepartor)
	if separatorIdx != -1 {
		if len(os.Args) == separatorIdx {
			fmt.Fprintf(os.Stderr, "Invalid invocation after %san application is expected", appSeparator)
			os.Exit(1)
		}
		appName = os.Args[separatorIdx + 1]
		if len(os.Args) > separatorIdx + 1{
			appArgs = os.Args[separatorIdx + 2:]
		}
		os.Args = os.Args[0:separatorIdx]
	} else {
		appArgs = make([]string, 0)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sp",
	Short: "Unix style stream processing utility",
	Long: fmt.Sprintf(`A helper for when you like unix style piping.

The default way to use it is to pipe it onto a command in order to process its stdout.
If you want to process both stdout and stderr and keep them appart you specify the sp
command followed by -- followed by the command you want to run. Potentially chained. For example:
	sp color --force %s sp epoch %s ./your_scripts/print_epochs.sh

Would color stderr and stdout differently (see color subcommand for defaults and it will also replace epochs).
`, appSeparator, appSeparator),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	proccessAppArguments()
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprint(os.Stderr, "Encountered issues processing sp initialization")
		os.Exit(1)
	}
	if appName == "" {
		if stdoutWriter == nil {
			fmt.Fprint(os.Stderr, "After sp initialization stdout writer was still nil")
			os.Exit(1)
		}
		os.Exit(streams.NewPipedApp(stdoutWriter).Run())
	} else {
		if stderrWriter == nil {
			stderrWriter = os.Stderr
		}
		if stdoutWriter == nil {
			fmt.Fprint(os.Stderr, "After sp initialization stdout writer was still nil")
			os.Exit(1)
		}
		os.Exit(streams.NewSpawnedApp(NewSyncedWriter(stdoutWriter), NewSyncedWriter(stderrWriter), appName, appArgs).Run())
	}
}

func init() {}


