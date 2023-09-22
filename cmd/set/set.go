/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package set

import (
	"fmt"
	"os"
	"strings"

	"github.com/dp1140a/semver/cmd"
	"github.com/dp1140a/semver/types"
	"github.com/dp1140a/semver/util"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var SetCmd = &cobra.Command{
	Use:   "set [new-version]",
	Args:  cobra.ExactArgs(1),
	Short: "sets the version, build, or pre-release values",
	Long: `By itself (with no subcommand) the set command will set the version to the passed in argument.  For example if our current version is 1.2.3:
   $semver version 4.5.6 --> 4.5.6
   $semver versiion 1.0.0-beta+exp.sha.5114f85 --> 1.0.0-beta+exp.sha.5114f85
`,
	Run: func(cmd *cobra.Command, args []string) {
		ver := ""
		if len(args) == 1 {
			ver = args[0]
		}

		ver = strings.TrimSuffix(ver, "\n")
		strings.TrimSpace(ver)
		setVersion(ver)
	},
}

func init() {
	cmd.RootCmd.AddCommand(SetCmd)
}

func setVersion(ver string) {
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
	fmt.Println("Setting Version")
	v = types.NewVersionFromString(ver)
	fmt.Printf("New Version: %v\n", v.String())
	err = util.WriteVersionFile(v.String())
	if err != nil {
		fmt.Printf("Error writing VERSION file: %v. Exiting", err)
		os.Exit(-1)
	}
	os.Exit(0)
}
