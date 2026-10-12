package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
)

// versionReport is flai version --json: the build metadata, and the release
// stamp and what verifying it found (S-0235).
type versionReport struct {
	buildinfo.Info
	Release buildinfo.Stamp `json:"release"`
}

func newVersionCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version, commit, build date, and whether the build is a signed release",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.printVersion(versionReport{Info: buildinfo.Get(), Release: buildinfo.OwnStamp()})
		},
	}
}

func (a *app) printVersion(v versionReport) error {
	if a.jsonOut {
		return a.printJSON(v)
	}
	_, err := fmt.Fprintf(a.out, "flai %s\ncommit: %s\nbuilt:  %s\ngo:     %s\n", v.Version, v.Commit, v.Date, v.Go)
	if err == nil {
		err = printRelease(a.out, v.Release)
	}
	return err
}

// printRelease says in one line whether the build is a signed release.
func printRelease(w io.Writer, s buildinfo.Stamp) error {
	line := "signed"
	switch {
	case s.State == buildinfo.StampMissing && s.Development:
		line = "development build, no stamp"
	case s.State == buildinfo.StampMissing:
		line = "missing stamp"
	case s.State == buildinfo.StampInvalid:
		line = "invalid stamp: " + s.Reason
	}
	_, err := fmt.Fprintf(w, "release: %s\n", line)
	return err
}
