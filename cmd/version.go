/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dp1140a/semver/types"
	"github.com/dp1140a/semver/util"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints the current version",
		Long: `Prints the current version in the chosen format. For example if the current version is 1.2.3 Format options:
   
   $ semver version --> 1.2.3 (string is default)
   $ semver version -f string --> 1.2.3
   $ semver version -f json -->
	{
	  "Major": 2,
	  "Minor": 0,
	  "Patch": 0,
	  "PreRelease": "",
	  "Build": ""
	}
   $ semver version -f pretty --> {Major: 2, Minor: 0, Patch: 0, PreRelease: "", Build: ""}

Pretty differs form json in that pretty is a pretty print of the underlying Version struct and is technically not valid json.
`,
		Run: func(cmd *cobra.Command, args []string) {
			format := cmd.Flag("format")
			runVersion(format.Value.String())
		},
	}
)

func init() {
	RootCmd.AddCommand(versionCmd)
	versionCmd.Flags().StringP("format", "f", "string", "Print Format [string | json | pretty]")
}

func runVersion(format string) {
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
	switch strings.ToLower(format) {
	case "string":
		fmt.Println(v.String())
		break
	case "json":
		fmt.Println(v.Json())
		break
	case "pretty":
		fmt.Println(v.PrettyPrint())
		break
	default:
		fmt.Printf("%v is an unknown format. Optons are [string | json | pretty]\n\n", format)
	}
}
