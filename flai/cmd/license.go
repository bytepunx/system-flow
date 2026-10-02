package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/license"
)

// newLicenseCmd prints the license built into the binary (S-0231): the terms
// flai, flaiover, and the system-flow repository are distributed under, so that
// a copy of the license travels with every copy of the software, as it requires.
func newLicenseCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "license",
		Short: "Print the license flai is distributed under",
		Long: `Print the Bytepunx Shield License: the terms flai, flaiover, and the
system-flow repository are distributed under. The text is built into the binary,
so it is the license of the very flai that prints it. It is LICENSE.md at the
root of the repository, and the dashboard shows the same text under Host.`,
		Example: `  flai license
  flai license --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if a.jsonOut {
				return a.printJSON(map[string]string{"name": license.Name(), "text": license.Text()})
			}
			_, err := fmt.Fprint(a.out, license.Text())
			return err
		},
	}
}
