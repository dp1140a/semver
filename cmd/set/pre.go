/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package set

import (
	"fmt"
	"os"

	"github.com/dp1140a/semver/types"
	"github.com/dp1140a/semver/util"
	"github.com/spf13/cobra"
)

// preCmd represents the pre command
var preCmd = &cobra.Command{
	Use:   "pre [(optional) pre-release value]",
	Args:  cobra.MaximumNArgs(1),
	Short: "Set Version Pre Release information",
	Long: `Will set the pre-release on a version.  For example if the current version is 1.2.3:

   $ semver set pre alpha-123 --> 1.2.3-alpha-123

If no pre string argument is given it will set the pre-release accordingly:
  If no pre-release value will set to alpha.  For example if the current version is 1.2.3
  
  $ semver set pre --> 1.2.3-alpha

  If pre-release value is alpha will set to beta.  For example if the current version is 1.2.3-alpha
  
  $ semver set pre --> 1.2.3-beta

  If pre-release value is beta will set to rc1.0.  For example if the current version is 1.2.3-beta
  
  $ semver set pre --> 1.2.3-rc1.0

NOTE: Setting the pre-release value WILL delete the current build value since pre-release is a higher precedence.
  
`,
	Run: func(cmd *cobra.Command, args []string) {
		pre := ""
		if len(args) == 1 {
			pre = args[0]
		}

		setPre(pre)
	},
}

func init() {
	SetCmd.AddCommand(preCmd)
}

func setPre(pre string) {
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
	fmt.Println("Setting Pre-Release")

	if pre == "" {
		switch v.PreRelease {
		case "": //if empty set to alpha
			pre = "alpha"
			break
		case "alpha": //if alpha set to beta
			pre = "beta"
			break
		case "beta": //if beta set to rc-1.0
			pre = "rc-1.0"
			break
		default: //if anything else notify and leave it alone
			fmt.Printf("Current Pre-Release value is %v. No value specified.  Exiting", v.PreRelease)
			os.Exit(0)
		}
	}

	v.SetPre(pre)
	v.SetBuild("")
	fmt.Printf("New Version: %v\n", v.String())
	err = util.WriteVersionFile(v.String())
	if err != nil {
		fmt.Printf("Error writing VERSION file: %v. Exiting", err)
		os.Exit(-1)
	}
	os.Exit(0)
}
