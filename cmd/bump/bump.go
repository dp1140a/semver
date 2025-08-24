package bump

import (
	"fmt"
	"os"
	"strings"

	"github.com/dp1140a/semver/cmd"
	"github.com/dp1140a/semver/pkg/cli"
	"github.com/dp1140a/semver/pkg/types"
	"github.com/spf13/cobra"
)

type bumpKind int

const (
	bumpPatch bumpKind = iota
	bumpMinor
	bumpMajor
)

var BumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "Bump the version (default: patch)",
	Long:  "Bump the semantic version in the VERSION file. Defaults to a patch bump if no subcommand is provided.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// default to patch when no subcommand is specified
		return runBump(cmd, bumpPatch)
	},
}

func init() {
	cmd.RootCmd.AddCommand(BumpCmd)

	// Per-call flag (no package-level globals)
	BumpCmd.PersistentFlags().BoolP(
		"dry", "d", false,
		"Show what the next version would be; do not write VERSION",
	)

	// Subcommands using the same runner
	BumpCmd.AddCommand(newBumpSubCmd("patch", "Bump patch version", bumpPatch))
	BumpCmd.AddCommand(newBumpSubCmd("minor", "Bump minor version", bumpMinor))
	BumpCmd.AddCommand(newBumpSubCmd("major", "Bump major version", bumpMajor))
}

func newBumpSubCmd(name, desc string, kind bumpKind) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: desc,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBump(cmd, kind)
		},
	}
}

func runBump(cmd *cobra.Command, kind bumpKind) error {
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

	v := types.NewVersionFromString(strings.TrimSpace(cur))

	switch kind {
	case bumpPatch:
		fmt.Println("Bumping Patch")
		v.IncrementPatch()
	case bumpMinor:
		fmt.Println("Bumping Minor")
		v.IncrementMinor()
	case bumpMajor:
		fmt.Println("Bumping Major")
		v.IncrementMajor()
	default:
		return fmt.Errorf("unknown bump kind: %v", kind)
	}

	next := v.String()
	if dry {
		cli.RenderDry(next)
		return nil
	}

	fmt.Printf("New Version: %s\n", next)
	return cli.WriteVersion(next)
}
