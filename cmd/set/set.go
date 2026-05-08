package set

import (
	"fmt"
	"os"
	"strings"

	"github.com/dp1140a/semver/cmd"
	"github.com/dp1140a/semver/pkg/cli"
	"github.com/dp1140a/semver/pkg/types"
	"github.com/spf13/cobra"
)

var SetCmd = &cobra.Command{
	Use:   "set <version>",
	Short: "Set the full semantic version",
	Long:  "Set the semantic version in the VERSION file (e.g., 1.2.3, 1.2.3-rc.1+build.5, or v1.2.3-rc.1). If the provided version omits a leading v/V, the existing VERSION file style is preserved.",
	Args:  cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		verArg := strings.TrimSpace(args[0])
		return runSetVersion(c, verArg)
	},
}

func init() {
	cmd.RootCmd.AddCommand(SetCmd)
	SetCmd.PersistentFlags().BoolP(
		"dry", "d", false,
		"Show what the new version would be; do not write VERSION",
	)
}

func runSetVersion(cmd *cobra.Command, verArg string) error {
	dry, _ := cmd.Flags().GetBool("dry")

	cwd, _ := os.Getwd()
	cur, err := cli.ReadVersion()
	if err != nil {
		return err
	}
	if cur == "" {
		cli.PrintNoVersionMsg(cwd)
		return nil
	}

	fmt.Printf("Current Version: %s\n", cur)
	fmt.Println("Setting Version")

	curVersion, err := types.ParseVersion(cur)
	if err != nil {
		return fmt.Errorf("invalid VERSION contents %q", cur)
	}

	v, err := types.ParseVersion(verArg)
	if err != nil {
		return fmt.Errorf("invalid semantic version %q", verArg)
	}
	if v.Prefix() == "" {
		v.SetPrefix(curVersion.Prefix())
	}
	next := v.String()

	if dry {
		cli.RenderDry(next)
		return nil
	}

	fmt.Printf("New Version: %s\n", next)
	return cli.WriteVersion(next)
}
