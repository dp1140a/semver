/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dp1140a/semver/pkg/types"
	"github.com/dp1140a/semver/pkg/util"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "semver",
	Short: "A semantic versioning tool",
	Long: `Run by itself semver will return the current version string. For example if the current version is 1.2.3:
   $ semver --> 1.2.3
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		runVersion("string")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

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

	versionStr := strings.TrimSpace(string(CUR_VER)) // <-- key fix
	v := types.NewVersionFromString(versionStr)

	switch strings.ToLower(format) {
	case "string":
		fmt.Println(v.String())
	case "json":
		fmt.Println(v.Json())
	case "pretty":
		fmt.Println(v.PrettyPrint())
	default:
		fmt.Printf("%v is an unknown format. Optons are [string | json | pretty]\n\n", format)
	}
}

func init() {}
