// Package cmd wires the flai command tree.
package cmd

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/logx"
	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// app carries state shared by all commands.
type app struct {
	configPath string // --config
	jsonOut    bool   // --json
	yes        bool   // --yes
	verbose    bool   // --verbose
	out        io.Writer
	errOut     io.Writer
	runner     execx.Runner
	log        *slog.Logger // structured events on errOut, see design/conventions/logging.md

	stdinIsTerminal *bool                              // tests override terminal detection
	confirm         func(title string) (bool, error)   // tests answer confirmations
	serveStarter    func() (serve.Status, bool, error) // tests do not start a process
	cwd             string                             // tests override the working directory
	clock           func() time.Time                   // tests override the clock
}

// Execute runs the CLI and returns the process exit code.
func Execute() int {
	return run(os.Args[1:], os.Stdout, os.Stderr)
}

func run(args []string, out, errOut io.Writer) int {
	root := newRootCmd(out, errOut)
	root.SetArgs(args)
	if err := root.Execute(); err != nil {
		var ee *exitError
		if errors.As(err, &ee) {
			if ee.msg != "" {
				rootApp(root).fail(err)
			}
			return ee.code
		}
		a := rootApp(root)
		a.fail(err)
		return 1
	}
	return 0
}

// fail logs the final error once, at the boundary, as a fatal event.
func (a *app) fail(err error) {
	if a.log == nil {
		a.initLogger()
	}
	logx.Fatal(a.log, "command failed", "component", "cmd", "err", err.Error())
}

// initLogger builds the logger from --verbose, LOG_LEVEL, LOG_FORMAT, and
// whether stderr is a terminal.
func (a *app) initLogger() {
	opt := logx.FromEnv()
	opt.Verbose = a.verbose
	if f, ok := a.errOut.(*os.File); ok {
		opt.IsTerminal = term.IsTerminal(int(f.Fd()))
	}
	a.log = logx.New(a.errOut, opt)
}

var apps = map[*cobra.Command]*app{}

func rootApp(root *cobra.Command) *app { return apps[root] }

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
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			a.initLogger()
		},
	}
	apps[root] = a
	root.SetOut(out)
	root.SetErr(errOut)
	pf := root.PersistentFlags()
	pf.StringVar(&a.configPath, "config", "", "path to config file (default ~/.flai/config.json, or $FLAI_CONFIG)")
	pf.BoolVar(&a.jsonOut, "json", false, "print structured JSON output")
	pf.BoolVarP(&a.yes, "yes", "y", false, "answer yes to confirmations")
	pf.BoolVarP(&a.verbose, "verbose", "v", false, "debug-level log events on stderr (LOG_LEVEL, LOG_FORMAT also apply)")

	root.AddCommand(
		newVersionCmd(a), newSelfUpgradeCmd(a), newConfigCmd(a), newNewCmd(a), newImportCmd(a), newUpgradeCmd(a), newTemplateCmd(a),
		newItemCmd(a, "epic"), newItemCmd(a, "story"), newItemCmd(a, "task"), newShowCmd(a),
		newMoveCmd(a), newBlockCmd(a), newUnblockCmd(a), newTouchesCmd(a), newEditCmd(a), newBoardCmd(a), newOrderCmd(a), newPushCmd(a),
		newStreamCmd(a), newArchiveCmd(a), newMigrateCmd(a), newCheckCmd(a), newStatsCmd(a), newPrimeCmd(a), newIssueCmd(a), newThreadCmd(a), newMCPCmd(a), newReleaseCmd(a), newAcceptCmd(a), newDashboardCmd(a), newDocCmd(a),
	)
	// Registered apart from the list above so that stories adding commands at
	// the same time do not meet on one line.
	root.AddCommand(newAdrCmd(a))
	root.AddCommand(newServeCmd(a))
	root.AddCommand(newHostAPICmd(a))
	return root
}

// loadConfig resolves the config path and loads (or creates) the file.
func (a *app) loadConfig() (config.Config, string, error) {
	path := config.ResolvePath(a.configPath)
	cfg, created, err := config.Load(path)
	if err != nil {
		return config.Config{}, path, err
	}
	if created {
		a.logger().Info("config created with defaults", "component", "config", "path", path)
	} else {
		a.logger().Debug("config loaded", "component", "config", "path", path)
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

// logger returns the structured logger, building it if a command ran
// without the root's PersistentPreRun (tests calling helpers directly).
func (a *app) logger() *slog.Logger {
	if a.log == nil {
		a.initLogger()
	}
	return a.log
}
