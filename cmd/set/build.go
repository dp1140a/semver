package set

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dp1140a/semver/pkg/cli"
	"github.com/dp1140a/semver/pkg/types"
	"github.com/dp1140a/semver/pkg/util"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build [value]",
	Short: "Set or clear the build metadata",
	Long:  "Set the build metadata (e.g., build.42). Pass a value directly, use --git to derive from git, or use --clear to remove it. With no value and no flags, the build metadata defaults to the short git HEAD hash. The existing VERSION file's optional leading v/V prefix is preserved.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// flags (read per-call; no globals)
		dry, _ := cmd.Flags().GetBool("dry")
		flagVal, _ := cmd.Flags().GetString("value")
		useGit, _ := cmd.Flags().GetBool("git")
		clear, _ := cmd.Flags().GetBool("clear")
		val := flagVal
		if len(args) == 1 {
			if flagVal != "" {
				return fmt.Errorf("provide the build metadata as either an argument or --value, not both")
			}
			val = args[0]
		}
		if val == "" && !useGit && !clear {
			useGit = true
		}

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
			return fmt.Errorf("provide build metadata, use --git, or use --clear")
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

		v, err := types.ParseVersion(strings.TrimSpace(cur))
		if err != nil {
			return fmt.Errorf("invalid VERSION contents %q", cur)
		}

		switch {
		case clear:
			v.SetBuild("")
		case useGit:
			out, err := exec.Command("git", "rev-parse", "--short", "HEAD").CombinedOutput()
			if err != nil {
				return fmt.Errorf("error getting git build info: %w", err)
			}
			build := strings.TrimSpace(string(out))
			if !util.ValidBuildString(build) {
				return fmt.Errorf("invalid git-derived build metadata %q", build)
			}
			v.SetBuild(build)
		default:
			if !util.ValidBuildString(val) {
				return fmt.Errorf("invalid build metadata %q", val)
			}
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
