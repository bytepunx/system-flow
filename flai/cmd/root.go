// Package cmd wires the flai command tree.
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// app carries state shared by all commands.
type app struct {
	configPath string // --config
	jsonOut    bool   // --json
	yes        bool   // --yes
	out        io.Writer
	errOut     io.Writer
	runner     execx.Runner

	stdinIsTerminal *bool            // tests override terminal detection
	cwd             string           // tests override the working directory
	clock           func() time.Time // tests override the clock
}

// Execute runs the CLI and returns the process exit code.
func Execute() int {
	return run(os.Args[1:], os.Stdout, os.Stderr)
}

func run(args []string, out, errOut io.Writer) int {
	root := newRootCmd(out, errOut)
	root.SetArgs(args)
	if err := root.Execute(); err != nil {
		fmt.Fprintf(errOut, "flai: %v\n", err)
		return 1
	}
	return 0
}

func newRootCmd(out, errOut io.Writer) *cobra.Command {
	return newRootCmdWith(&app{out: out, errOut: errOut})
}

func newRootCmdWith(a *app) *cobra.Command {
	if a.runner == nil {
		a.runner = execx.System{}
	}
	out, errOut := a.out, a.errOut
	root := &cobra.Command{
		Use:   "flai",
		Short: "Manage monorepos that follow the system-flow standard",
		Long: `flai creates and imports projects that follow the system-flow standard,
manages work items and agent narratives, validates the repository, prints
flow metrics, and runs the flaiover dashboard.

Configuration is read from ~/.flai/config.json (override with --config or
FLAI_CONFIG). Every command that prints data accepts --json.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.SetOut(out)
	root.SetErr(errOut)
	pf := root.PersistentFlags()
	pf.StringVar(&a.configPath, "config", "", "path to config file (default ~/.flai/config.json, or $FLAI_CONFIG)")
	pf.BoolVar(&a.jsonOut, "json", false, "print structured JSON output")
	pf.BoolVarP(&a.yes, "yes", "y", false, "answer yes to confirmations")

	root.AddCommand(
		newVersionCmd(a), newConfigCmd(a), newNewCmd(a), newTemplateCmd(a),
		newItemCmd(a, "epic"), newItemCmd(a, "story"), newItemCmd(a, "task"), newShowCmd(a),
		newMoveCmd(a), newBlockCmd(a), newUnblockCmd(a), newBoardCmd(a),
		newStreamCmd(a), newArchiveCmd(a),
	)
	return root
}

// loadConfig resolves the config path and loads (or creates) the file.
func (a *app) loadConfig() (config.Config, string, error) {
	path := config.ResolvePath(a.configPath)
	cfg, created, err := config.Load(path)
	if err != nil {
		return config.Config{}, path, err
	}
	if created && !a.jsonOut {
		fmt.Fprintf(a.errOut, "flai: created %s with defaults\n", path)
	}
	return cfg, path, nil
}

// printJSON writes v as indented JSON to stdout.
func (a *app) printJSON(v any) error {
	enc := json.NewEncoder(a.out)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// relPath shows a path relative to root when possible.
func relPath(root, p string) string {
	if rel, err := filepath.Rel(root, p); err == nil {
		return rel
	}
	return p
}
