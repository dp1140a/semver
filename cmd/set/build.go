package set

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dp1140a/semver/pkg/cli"
	"github.com/dp1140a/semver/pkg/types"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Set or clear the build metadata",
	Long:  "Set the build metadata (e.g., build.42). Use --git to derive from git, or --clear to remove it.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// flags (read per-call; no globals)
		dry, _ := cmd.Flags().GetBool("dry")
		val, _ := cmd.Flags().GetString("value")
		useGit, _ := cmd.Flags().GetBool("git")
		clear, _ := cmd.Flags().GetBool("clear")

		// exactly one of --value, --git, --clear
		count := 0
		if val != "" {
			count++
		}
		if useGit {
			count++
		}
		if clear {
			count++
		}
		if count != 1 {
			return fmt.Errorf("exactly one of --value, --git, or --clear must be provided")
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
		fmt.Println("Setting Build Metadata")

		v := types.NewVersionFromString(strings.TrimSpace(cur))

		switch {
		case clear:
			v.SetBuild("")
		case useGit:
			out, err := exec.Command("git", "rev-parse", "--short", "HEAD").CombinedOutput()
			if err != nil {
				return fmt.Errorf("error getting git build info: %w", err)
			}
			v.SetBuild(strings.TrimSpace(string(out)))
		default:
			v.SetBuild(val)
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
	SetCmd.AddCommand(buildCmd)
	buildCmd.Flags().String("value", "", "Build metadata to set (e.g., build.42)")
	buildCmd.Flags().Bool("git", false, "Use git rev-parse --short HEAD for build metadata")
	buildCmd.Flags().Bool("clear", false, "Clear the build metadata")
}
