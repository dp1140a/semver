/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package set

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dp1140a/semver/types"
	"github.com/dp1140a/semver/util"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [(optional) build value]",
	Args:  cobra.MaximumNArgs(1),
	Short: "Set Version Build information",
	Long: `Will set the build on a version.  For example if the current version is 1.2.3:

   $ semver set build mybuild-123 --> 1.2.3+mybuild-123

If no build string argument is given it will set the build to the short version of the current git HEAD hash.
This is equivalent to setting the build to the output of:

   $ git rev-parse --short HEAD

For example if the current version is 1.2.3:

  $ semver set build --> 1.2.3+b113571 (if that was the current hash)
`,
	Run: func(cmd *cobra.Command, args []string) {
		build := ""
		if len(args) == 1 {
			build = args[0]
		}

		setBuild(build)
	},
}

func init() {
	SetCmd.AddCommand(buildCmd)
}

func setBuild(build string) {
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
	fmt.Println("Setting Build")
	if build == "" {
		cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
		bytes, err := cmd.Output()
		if err != nil {
			fmt.Printf("Error getting git build info: %v", err)
			os.Exit(-1)
		}
		build = string(bytes)
	}

	v.SetBuild(build)
	fmt.Printf("New Version: %v\n", v.String())
	err = util.WriteVersionFile(v.String())
	if err != nil {
		fmt.Printf("Error writing VERSION file: %v. Exiting", err)
		os.Exit(-1)
	}
	os.Exit(0)
}
