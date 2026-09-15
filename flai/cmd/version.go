package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
)

func newVersionCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version, commit, and build date",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			info := buildinfo.Get()
			if a.jsonOut {
				return a.printJSON(info)
			}
			_, err := fmt.Fprintf(a.out, "flai %s\ncommit: %s\nbuilt:  %s\ngo:     %s\n", info.Version, info.Commit, info.Date, info.Go)
			return err
		},
	}
}
