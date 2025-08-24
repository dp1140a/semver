package set

import (
	"fmt"
	"os"

	"github.com/dp1140a/semver/pkg/cli"
	"github.com/dp1140a/semver/pkg/types"
	"github.com/spf13/cobra"
)

var preCmd = &cobra.Command{
	Use:   "pre",
	Short: "Set or clear the prerelease identifier",
	Long:  "Set the prerelease (e.g., rc.1). Use --clear to remove it.",
	RunE: func(cmd *cobra.Command, args []string) error {
		dry, _ := cmd.Flags().GetBool("dry")
		val, _ := cmd.Flags().GetString("value")
		clr, _ := cmd.Flags().GetBool("clear")

		if (val == "" && !clr) || (val != "" && clr) {
			return fmt.Errorf("exactly one of --value or --clear must be provided")
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

		fmt.Printf("Current Version: %s\n", cur)
		fmt.Println("Setting Prerelease")

		v := types.NewVersionFromString(cur)
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
