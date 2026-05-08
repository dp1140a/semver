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
	bumpPre
	bumpBuild
)

var BumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "Bump the version (default: patch)",
	Long:  "Bump the semantic version in the VERSION file. Defaults to a patch bump if no subcommand is provided. The existing VERSION file's optional leading v/V prefix is preserved. Prerelease stage promotion remains explicit via set pre.",
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
	BumpCmd.AddCommand(newBumpSubCmd("pre", "Bump prerelease version", bumpPre))
	BumpCmd.AddCommand(newBumpSubCmd("build", "Bump build metadata", bumpBuild))
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

	v, err := types.ParseVersion(strings.TrimSpace(cur))
	if err != nil {
		return fmt.Errorf("invalid VERSION contents %q", cur)
	}

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
	case bumpPre:
		fmt.Println("Bumping Prerelease")
		v.IncrementPre()
	case bumpBuild:
		fmt.Println("Bumping Build Metadata")
		if err := v.IncrementBuild(); err != nil {
			return err
		}
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
