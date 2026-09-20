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
)

// flai hostapi calls one method of the table flai serve offers dashboards
// (ADR-0029), for the project in the working directory, without a
// connection: what a dashboard would be answered, on a terminal. It is how
// the table is looked at and debugged, and how flaiover's tests read a
// fixture through the same code the dashboard will be served by.
func newHostAPICmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "hostapi [method] [params-json]",
		Short: "Answer one method of the dashboard's API for this project, as flai serve would",
		Long: `Without arguments, list the methods flai serve offers a dashboard. With a
method, and optionally its params as a JSON object, print the result as
JSON. An error is printed as {"error": {"code", "message"}} with exit 2.`,
		Example: `  flai hostapi
  flai hostapi board.get '{"all":true}'
  flai hostapi item.get '{"id":"S-0042"}'`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			methods := hostapi.Methods(buildinfo.Version, a.now)
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
			res, rerr := method(cmd.Context(), p, params)
			enc := json.NewEncoder(a.out)
			if rerr != nil {
				_ = enc.Encode(map[string]any{"error": rerr})
				return &exitError{code: 2}
			}
			return enc.Encode(res)
		},
	}
}
