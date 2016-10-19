package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	Version   string
	GitCommit string
	BuildDate string
)

// versionCmd represents the plan command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "display version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\nGit Commit: %s\nBuild Date: %s\n",
			color.GreenString(Version), color.GreenString(GitCommit), color.GreenString(BuildDate))
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
