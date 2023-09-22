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

// majorCmd represents the major command
var majorCmd = &cobra.Command{
	Use:   "major",
	Short: "Will bump the current Major version",
	Long: `If our current version is 0.1.0:
   $ semver bump minor --> 1.0.0

Bumping will reset all lower order versions to 0 and remove build or pre-release values.  For example if our current version is 1.2.3-alpha:
   
   $ semver bump major --> 2.0.0
`,
	Run: func(cmd *cobra.Command, args []string) {
		runBumpMajor()
	},
}

func init() {
	BumpCmd.AddCommand(majorCmd)
}

func runBumpMajor() {
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
	fmt.Println("Bumping Major")
	v.IncrementMajor()
	fmt.Printf("New Version: %v\n", v.String())
	err = util.WriteVersionFile(v.String())
	if err != nil {
		fmt.Printf("Error writing VERSION file: %v. Exiting", err)
		os.Exit(-1)
	}
	os.Exit(0)
}
