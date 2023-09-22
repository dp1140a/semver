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

// patchCmd represents the patch command
var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Will bump the current Patch version",
	Long: `If our current version is 0.1.0:

   $ semver bump patch --> 0.1.1

Bumping will reset all lower order versions to 0 and remove build or pre-release values.  For example if our current version is 1.2.3-alpha:

   $ semver bump patch --> 1.2.4
`,
	Run: func(cmd *cobra.Command, args []string) {
		runBumpPatch()
	},
}

func init() {
	BumpCmd.AddCommand(patchCmd)
}

func runBumpPatch() {
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
	fmt.Println("Bumping Patch")
	v.IncrementPatch()
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
