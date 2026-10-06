package issues

import (
	"fmt"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// What strategic agents spent on an issue (S-0227): an analyzer's activity
// is charged to the issues it named, in the shape ADR-0083 gives an item's
// strategic usage. An issue has no agents' figures of its own, so its usage
// carries nothing but strategic entries.

// ChargeStrategic adds what the strategic agent kind spent on the issue with
// id, the agents' figures of u as apportioned to it, to that issue's
// strategic usage, and saves it. The issue is read just before it is saved,
// and its updated is left as it is, as an item's is. It returns the issue
// charged, or nil when u spent nothing.
func ChargeStrategic(r *workitem.Repo, id, kind string, u *usage.Usage) (*Issue, error) {
	if u.Empty() {
		return nil, nil
	}
	if !workitem.IsActivityKind(kind) {
		return nil, fmt.Errorf("%s is not a strategic agent's kind: %s", kind, strings.Join(workitem.ActivityKinds, ", "))
	}
	is, err := Get(r, id)
	if err != nil {
		return nil, fmt.Errorf("charge the %s's usage to issue %s: %w", kind, id, err)
	}
	if is.Usage == nil {
		is.Usage = &usage.Usage{}
	}
	is.Usage.AddStrategic(kind, u)
	if err := is.Save(); err != nil {
		return nil, fmt.Errorf("charge the %s's usage to issue %s: %w", kind, is.ID, err)
	}
	return is, nil
}

// usageBlock is an issue's usage as front matter, its strategic entries
// written as an item's are; "" when it has none.
func usageBlock(u *usage.Usage) string {
	if u == nil || len(u.Strategic) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("usage:\n")
	workitem.StrategicBlock(&b, u.Strategic)
	return b.String()
}

// usageErrors are what is wrong with an issue's usage: any agents' figures,
// which an issue does not have, and any strategic entry an item's usage
// would be refused for.
func usageErrors(u *usage.Usage) []string {
	if u == nil {
		return nil
	}
	var errs []string
	if u.Source != "" || u.Seconds != 0 || u.Estimated || len(u.Models) > 0 {
		errs = append(errs, "usage carries agents' figures (source, seconds, estimated, models): an issue's usage is its strategic entries only")
	}
	return append(errs, workitem.StrategicErrors(u.Strategic, "an issue")...)
}
