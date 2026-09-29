package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/perf"
)

// flai hostapi calls one method of the table flai serve offers dashboards
// (ADR-0029), for the project in the working directory, without a
// connection: what a dashboard would be answered, on a terminal. It is how
// the table is looked at and debugged, and how flaiover's tests read a
// fixture through the same code the dashboard will be served by.
func newHostAPICmd(a *app) *cobra.Command {
	var timing bool
	c := &cobra.Command{
		Use:   "hostapi [method] [params-json]",
		Short: "Answer one method of the dashboard's API for this project, as flai serve would",
		Long: `Without arguments, list the methods flai serve offers a dashboard. With a
method, and optionally its params as a JSON object, print the result as
JSON. An error is printed as {"error": {"code", "message"}} with exit 2.

With --timing, the request is timed as flai serve times it (S-0152), with
no connection at all, and its "request answered" event, with the phases
the time went to, is logged on stderr at info whatever it took.`,
		Example: `  flai hostapi
  flai hostapi board.get '{"all":true}'
  flai hostapi item.get '{"id":"S-0042"}'
  flai hostapi --timing inbox.designer 2>&1 >/dev/null`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			methods := hostapi.MethodsFor(buildinfo.Version, a.now, a.host())
			if len(args) == 0 {
				names := make([]string, 0, len(methods))
				for name := range methods {
					names = append(names, name)
				}
				sort.Strings(names)
				if a.jsonOut {
					return a.printJSON(names)
				}
				fmt.Fprintln(a.out, strings.Join(names, "\n"))
				return nil
			}
			repo, err := a.project()
			if err != nil {
				return err
			}
			method, ok := methods[args[0]]
			if !ok {
				return fmt.Errorf("flai offers no method %s; flai hostapi lists them", args[0])
			}
			params := json.RawMessage(`{}`)
			if len(args) == 2 {
				if !json.Valid([]byte(args[1])) {
					return fmt.Errorf("params must be a JSON object")
				}
				params = json.RawMessage(args[1])
			}
			p := channel.Project{Key: repo.Manifest.Key, Name: repo.Manifest.Name, Root: repo.MainRoot}
			ctx, rec := cmd.Context(), (*perf.Recorder)(nil)
			if timing {
				ctx, rec = perf.Start(ctx)
			}
			res, rerr := method(ctx, p, params)
			answer, err := json.Marshal(map[string]any{"error": rerr})
			if rerr == nil {
				done := perf.Track(ctx, "encode")
				answer, err = json.Marshal(res)
				done()
			}
			if err != nil {
				return err
			}
			if rec != nil {
				answered := perf.Answered{Transport: "none", Method: args[0], Project: p.Key, Bytes: len(answer)}
				if rerr != nil {
					answered.Err = rerr.Message
				}
				rec.Log(a.logger(), answered, 0)
			}
			fmt.Fprintln(a.out, string(answer))
			if rerr != nil {
				return &exitError{code: 2}
			}
			return nil
		},
	}
	c.Flags().BoolVar(&timing, "timing", false, "log the request's duration and phases on stderr, as flai serve times it")
	return c
}
