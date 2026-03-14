package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var debugFlag bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-revpr",
	Short: "A brief description of your application",
	Long:  `todo`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		SetupLogging(debugFlag)
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "d", false, "Enable debug output (JSON to stderr)")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
