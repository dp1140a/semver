/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package bump

import (
	"fmt"
	"os"

	"github.com/dp1140a/semver/types"
	"github.com/dp1140a/semver/util"
	"github.com/spf13/cobra"
)

// minorCmd represents the minor command
var minorCmd = &cobra.Command{
	Use:   "minor",
	Short: "Will bump the current Minor version",
	Long: `If our current version is 0.1.0:
   $ semver bump minor --> 0.2.0

Bumping will reset all lower order versions to 0 and remove build or pre-release values.  For example if our current version is 1.2.3-alpha:

   $ semver bump minor --> 1.3.0

`,
	Run: func(cmd *cobra.Command, args []string) {
		runBumpMinor()
	},
}

func init() {
	BumpCmd.AddCommand(minorCmd)
}

func runBumpMinor() {
	cwd, _ := os.Getwd()
	if !util.VersionFileExists(cwd) {
		fmt.Printf("No VERSION file found in %v.\nPlease either change directory or first run 'semver init'\n", cwd)
		os.Exit(0)
	}

	CUR_VER, err := os.ReadFile("VERSION")
	if err != nil {
		fmt.Printf("Error reading VERSION file. %v", err)
	}
	v := types.NewVersionFromString(string(CUR_VER))
	fmt.Printf("Current Version: %v\n", v.String())
	fmt.Println("Bumping Minor")
	v.IncrementMinor()
	fmt.Printf("New Version: %v\n", v.String())
	if DryRun {
		os.Exit(0)
	} else {
		err = util.WriteVersionFile(v.String())
		if err != nil {
			fmt.Printf("Error writing VERSION file: %v. Exiting", err)
			os.Exit(-1)
		}
		os.Exit(0)
	}
}
