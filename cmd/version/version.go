package version

import (
	"fmt"

	"github.com/dp1140a/semver/cmd"
	"github.com/dp1140a/semver/pkg/version/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the current version of the semver binary",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.StringPretty())
	},
}

func init() {
	cmd.RootCmd.AddCommand(versionCmd)
}
