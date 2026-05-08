package set

import (
	"fmt"
	"os"
	"strings"

	"github.com/dp1140a/semver/pkg/cli"
	"github.com/dp1140a/semver/pkg/types"
	"github.com/dp1140a/semver/pkg/util"
	"github.com/spf13/cobra"
)

var preCmd = &cobra.Command{
	Use:   "pre [value]",
	Short: "Set or clear the prerelease identifier",
	Long:  "Set the prerelease (e.g., rc.1, alpha.2, beta.1). Use --clear to remove it. The existing VERSION file's optional leading v/V prefix is preserved. Stage changes like alpha.2 -> beta.1 are explicit set operations rather than bump behavior.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dry, _ := cmd.Flags().GetBool("dry")
		flagVal, _ := cmd.Flags().GetString("value")
		clr, _ := cmd.Flags().GetBool("clear")
		val := flagVal
		if len(args) == 1 {
			if flagVal != "" {
				return fmt.Errorf("provide the prerelease as either an argument or --value, not both")
			}
			val = args[0]
		}

		if (val == "" && !clr) || (val != "" && clr) {
			return fmt.Errorf("provide a prerelease value or --clear")
		}

		cwd, _ := os.Getwd()
		cur, err := cli.ReadVersion()
		if err != nil {
			return err
		}
		if cur == "" {
			cli.PrintNoVersionMsg(cwd)
			return nil
		}
		if val != "" {
			val = strings.TrimSpace(val)
			if !util.ValidPrereleaseString(val) {
				return fmt.Errorf("invalid prerelease value %q", val)
			}
		}

		fmt.Printf("Current Version: %s\n", cur)
		fmt.Println("Setting Prerelease")

		v, err := types.ParseVersion(cur)
		if err != nil {
			return fmt.Errorf("invalid VERSION contents %q", cur)
		}
		if clr {
			v.SetPre("")
		} else {
			v.SetPre(val)
		}

		next := v.String()
		if dry {
			cli.RenderDry(next)
			return nil
		}
		fmt.Printf("New Version: %s\n", next)
		return cli.WriteVersion(next)
	},
}

func init() {
	SetCmd.AddCommand(preCmd)
	preCmd.Flags().String("value", "", "Prerelease value to set (e.g., rc.1)")
	preCmd.Flags().Bool("clear", false, "Clear the prerelease")
}
