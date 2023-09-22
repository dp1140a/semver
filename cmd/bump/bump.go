/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package bump

import (
	"github.com/dp1140a/semver/cmd"
	"github.com/spf13/cobra"
)

// bumpCmd represents the bump command
var (
	DryRun bool

	BumpCmd = &cobra.Command{
		Use:   "bump",
		Short: "Will bump the current version",
		Long: `If no subcommand is specified this command will bump the Patch version.  For example if our current version is 0.1.0:

   $ semver bump --> 0.1.1 
   Is the same as
   $ semver bump patch --> 0.1.1

Bumping will reset all lower order versions to 0 and remove build or pre-release values.  For instance if our current version is 1.2.3-alpha:
   
   $ semver bump major --> 2.0.0
   $ semver bump minor --> 1.3.0
   $ semver bump patch --> 1.2.4

`,
		Run: func(cmd *cobra.Command, args []string) {
			runBump()
		},
	}
)

var CUR_VER []byte

func init() {
	cmd.RootCmd.AddCommand(BumpCmd)
	BumpCmd.PersistentFlags().BoolVarP(&DryRun, "dry", "d", false, "Shows what the next version will be.  Will not write VERSION file")

}

func runBump() {
	runBumpPatch()
}
