package manifest

import (
	"fmt"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
)

// The values of orchestration.policy, the same as flai order --by's.
const (
	// OrderCOD orders by cost of delay, the largest first.
	OrderCOD = "cod"
	// OrderWSJF orders by cost of delay divided by forecast duration, the
	// largest first.
	OrderWSJF = "wsjf"
	// OrderThroughput orders by forecast duration, the shortest first.
	OrderThroughput = "throughput"
	// OrderFIFO leaves the operator's order alone. It is the default.
	OrderFIFO = "fifo"
)

// OrderPolicies are the values of orchestration.policy and flai order --by,
// in the order they are listed to the operator.
var OrderPolicies = []string{OrderCOD, OrderWSJF, OrderThroughput, OrderFIFO}

// Release is the project's say about when accepted stories not yet released
// are due a release.
type Release struct {
	// Policy is one of ReleasePolicies; empty means ReleaseJudgement.
	Policy string `yaml:"policy,omitempty" json:"policy,omitempty"`
	// Value is, for ReleaseThreshold, the unreleased cost of delay per week,
	// in the project's currency, at which a release is due; unset means none.
	Value *float64 `yaml:"value,omitempty" json:"value,omitempty"`
	// Count is, for ReleaseThreshold, the number of accepted stories not yet
	// released at which a release is due; unset means none.
	Count *int `yaml:"count,omitempty" json:"count,omitempty"`
	// Epic is, for ReleaseTheme, the epic whose stories, every one accepted,
	// make a release due.
	Epic string `yaml:"epic,omitempty" json:"epic,omitempty"`
	// Tag is, for ReleaseTheme, the tag whose stories, every one accepted,
	// make a release due.
	Tag string `yaml:"tag,omitempty" json:"tag,omitempty"`
	// WholeEpics, under every policy, holds a release back while a story
	// accepted and not yet released belongs to an epic in neither review nor
	// done (S-0222). Off when unset.
	WholeEpics bool `yaml:"whole_epics,omitempty" json:"whole_epics,omitempty"`
}

// The values of orchestration.release.policy.
const (
	// ReleaseJudgement leaves the release to the operator: it is never met by
	// itself. It is the default.
	ReleaseJudgement = "judgement"
	// ReleaseThreshold is met by the unreleased cost of delay per week against
	// Value, or the accepted stories not yet released against Count.
	ReleaseThreshold = "threshold"
	// ReleaseTheme is met when every story of an epic or a tag is accepted.
	ReleaseTheme = "theme"
)

// ReleasePolicies are the values of orchestration.release.policy.
var ReleasePolicies = []string{ReleaseJudgement, ReleaseThreshold, ReleaseTheme}

var epicID = regexp.MustCompile(`^E-\d{3,}$`)

// PolicyOrDefault is orchestration.policy, or OrderFIFO when it is empty. It
// does not say whether the policy is one of OrderPolicies; Errors does.
func (o Orchestration) PolicyOrDefault() string {
	if s := strings.TrimSpace(o.Policy); s != "" {
		return s
	}
	return OrderFIFO
}

// PolicyOrDefault is orchestration.release.policy, or ReleaseJudgement when
// it is empty. It does not say whether the policy is one of ReleasePolicies;
// Errors does.
func (r Release) PolicyOrDefault() string {
	if s := strings.TrimSpace(r.Policy); s != "" {
		return s
	}
	return ReleaseJudgement
}

// Errors are what is wrong with the orchestration settings, one sentence
// each; none when they are valid.
func (o Orchestration) Errors() []string {
	var errs []string
	if p := o.PolicyOrDefault(); !slices.Contains(OrderPolicies, p) {
		errs = append(errs, fmt.Sprintf("orchestration.policy %q is not an order policy; write cod (cost of delay), wsjf (cost of delay by duration), throughput (shortest first), or fifo (the operator's order), or remove it for fifo", o.Policy))
	}
	errs = append(errs, o.Release.errors()...)
	errs = append(errs, o.Permissions.errors("orchestration.permissions")...)
	if _, err := o.Permissions.ConfirmAfterDuration(); err != nil {
		errs = append(errs, err.Error())
	}
	if _, err := o.CoalesceDuration(); err != nil {
		errs = append(errs, err.Error())
	}
	if _, err := o.RunTimeoutDuration(); err != nil {
		errs = append(errs, err.Error())
	}
	if _, err := o.OrchestrationSchedule(); err != nil {
		errs = append(errs, err.Error())
	}
	return append(errs, o.Agent.problems("orchestration.agent")...)
}

func (r Release) errors() []string {
	var errs []string
	switch r.PolicyOrDefault() {
	case ReleaseJudgement:
	case ReleaseThreshold:
		if r.Value == nil && r.Count == nil {
			errs = append(errs, "orchestration.release is a threshold with neither value nor count; write value (the unreleased cost of delay per week), count (the accepted stories not yet released), or both")
		}
		if v := r.Value; v != nil && (math.IsNaN(*v) || math.IsInf(*v, 0) || *v < 0) {
			errs = append(errs, fmt.Sprintf("orchestration.release.value %v is not an amount of zero or more; write the unreleased cost of delay per week, in the project's currency, at which a release is due", *v))
		}
		if c := r.Count; c != nil && *c < 0 {
			errs = append(errs, fmt.Sprintf("orchestration.release.count %d is not a number of zero or more; write the number of accepted stories not yet released at which a release is due", *c))
		}
	case ReleaseTheme:
		epic, tag := strings.TrimSpace(r.Epic), strings.TrimSpace(r.Tag)
		switch {
		case epic == "" && tag == "":
			errs = append(errs, "orchestration.release is a theme with neither epic nor tag; write epic (an epic ID such as E-0001) or tag (a tag), one of them")
		case epic != "" && tag != "":
			errs = append(errs, "orchestration.release is a theme with both epic and tag; keep one of them")
		case epic != "" && !epicID.MatchString(epic):
			errs = append(errs, fmt.Sprintf("orchestration.release.epic %q is not an epic ID; write one such as E-0001", r.Epic))
		}
	default:
		errs = append(errs, fmt.Sprintf("orchestration.release.policy %q is not a release policy; write judgement (the operator decides), threshold (a value or a count of unreleased work), or theme (an epic or a tag accepted), or remove it for judgement", r.Policy))
	}
	return errs
}

// Acceptance is the project's say about how a story in review is accepted
// (S-0392). Each key empty means its default, so a manifest without the
// block accepts as flai always has.
type Acceptance struct {
	// Mode is one of AcceptanceModes; empty means AcceptLocal.
	Mode string `yaml:"mode,omitempty" json:"mode,omitempty"`
	// Forge is the forge, one of Forges, that pull-request acceptance goes
	// through; empty means ForgeGitHub. Under AcceptLocal it is unused.
	Forge string `yaml:"forge,omitempty" json:"forge,omitempty"`
	// MergeMethod is how the forge merges an accepted pull request, one of
	// MergeMethods; empty means MergeMerge.
	MergeMethod string `yaml:"merge_method,omitempty" json:"merge_method,omitempty"`
}

// The values of acceptance.mode.
const (
	// AcceptLocal merges a story's branch on this machine, as flai accept
	// does. It is the default.
	AcceptLocal = "local"
	// AcceptPullRequest accepts a story through a pull request on the
	// project's forge.
	AcceptPullRequest = "pull-request"
)

// AcceptanceModes are the values of acceptance.mode, in the order they are
// listed to the operator.
var AcceptanceModes = []string{AcceptLocal, AcceptPullRequest}

// ForgeGitHub is GitHub, the default forge.
const ForgeGitHub = "github"

// Forges are the values of acceptance.forge.
var Forges = []string{ForgeGitHub}

// The values of acceptance.merge_method, as the forge merges a pull request.
const (
	// MergeMerge merges with a merge commit. It is the default.
	MergeMerge = "merge"
	// MergeSquash squashes the branch into one commit.
	MergeSquash = "squash"
	// MergeRebase rebases the branch's commits onto the base.
	MergeRebase = "rebase"
)

// MergeMethods are the values of acceptance.merge_method, in the order they
// are listed to the operator.
var MergeMethods = []string{MergeMerge, MergeSquash, MergeRebase}

// ModeOrDefault is acceptance.mode, or AcceptLocal when it is empty. It does
// not say whether the mode is one of AcceptanceModes; Errors does.
func (a Acceptance) ModeOrDefault() string {
	if s := strings.TrimSpace(a.Mode); s != "" {
		return s
	}
	return AcceptLocal
}

// PullRequest reports whether stories are accepted through pull requests.
func (a Acceptance) PullRequest() bool {
	return a.ModeOrDefault() == AcceptPullRequest
}

// ForgeOrDefault is acceptance.forge, or ForgeGitHub when it is empty, under
// either mode. It does not say whether the forge is one of Forges; Errors
// does.
func (a Acceptance) ForgeOrDefault() string {
	if s := strings.TrimSpace(a.Forge); s != "" {
		return s
	}
	return ForgeGitHub
}

// MergeMethodOrDefault is acceptance.merge_method, or MergeMerge when it is
// empty. It does not say whether the method is one of MergeMethods; Errors
// does.
func (a Acceptance) MergeMethodOrDefault() string {
	if s := strings.TrimSpace(a.MergeMethod); s != "" {
		return s
	}
	return MergeMerge
}

// Errors are what is wrong with the acceptance settings, one sentence each;
// none when they are valid. A forge or a merge method set under AcceptLocal
// is valid, and unused.
func (a Acceptance) Errors() []string {
	var errs []string
	if m := a.ModeOrDefault(); !slices.Contains(AcceptanceModes, m) {
		errs = append(errs, fmt.Sprintf("acceptance.mode %q is not an acceptance mode; write local (merge on this machine) or pull-request (merge through a pull request on the forge), or remove it for local", a.Mode))
	}
	if f := a.ForgeOrDefault(); !slices.Contains(Forges, f) {
		errs = append(errs, fmt.Sprintf("acceptance.forge %q is not a forge flai reaches; write %s, or remove it for github", a.Forge, strings.Join(Forges, ", ")))
	}
	if m := a.MergeMethodOrDefault(); !slices.Contains(MergeMethods, m) {
		errs = append(errs, fmt.Sprintf("acceptance.merge_method %q is not a merge method; write merge (a merge commit), squash (one commit), or rebase (the commits onto the base), or remove it for merge", a.MergeMethod))
	}
	return errs
}

// MergeQueue is the project's say about the merge queue flai serve runs
// behind the merge-queue host action (S-0399, ADR-0173). A manifest without
// the block batches DefaultBatchMax stories at most, and waits
// DefaultMergeQueueCoalesce for a batch to fill (S-0552).
type MergeQueue struct {
	// BatchMax is the most enqueued stories one batch holds, a whole number
	// of one or more; unset means DefaultBatchMax.
	BatchMax *int `yaml:"batch_max,omitempty" json:"batch_max,omitempty"`
	// Coalesce is how long after the first of a batch's stories was
	// enqueued the queue waits for more before it builds a batch of fewer
	// than BatchMax, a Go duration such as 3m; empty means
	// DefaultMergeQueueCoalesce, and 0 builds it at once.
	Coalesce string `yaml:"coalesce,omitempty" json:"coalesce,omitempty"`
}

// DefaultBatchMax is the most stories a batch of the merge queue holds when
// merge_queue.batch_max is not set (ADR-0173).
const DefaultBatchMax = 4

// DefaultMergeQueueCoalesce is how long the merge queue waits for a batch to
// fill when merge_queue.coalesce is not set: 3 minutes.
const DefaultMergeQueueCoalesce = 3 * time.Minute

// CoalesceDuration is merge_queue.coalesce as a duration:
// DefaultMergeQueueCoalesce when it is empty, and zero when the queue waits
// for none.
func (q MergeQueue) CoalesceDuration() (time.Duration, error) {
	s := strings.TrimSpace(q.Coalesce)
	if s == "" {
		return DefaultMergeQueueCoalesce, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("merge_queue.coalesce %q is not a duration; write how long the merge queue waits for a batch to fill, such as 3m or 90s, 0 to build each batch at once, or remove it for 3m", q.Coalesce)
	}
	if d < 0 {
		return 0, fmt.Errorf("merge_queue.coalesce %q is negative; write how long the merge queue waits for a batch to fill, such as 3m or 90s, 0 to build each batch at once, or remove it for 3m", q.Coalesce)
	}
	return d, nil
}

// MergeQueueCoalesce is merge_queue.coalesce as a duration, zero when the
// queue waits for no batch to fill, or DefaultMergeQueueCoalesce when it is
// not set or not valid, which MergeQueue.Errors refuses at load.
func (m Manifest) MergeQueueCoalesce() time.Duration {
	d, err := m.MergeQueue.CoalesceDuration()
	if err != nil {
		return DefaultMergeQueueCoalesce
	}
	return d
}

// BatchMax is merge_queue.batch_max, or DefaultBatchMax when it is not set.
// It does not say whether the value is one or more; MergeQueue.Errors does.
func (m Manifest) BatchMax() int {
	if b := m.MergeQueue.BatchMax; b != nil {
		return *b
	}
	return DefaultBatchMax
}

// batchMaxProblem says why merge_queue.batch_max, as data writes it, is not
// a whole number, which a load into MergeQueue would otherwise fail on with
// the decoder's words or, for a fraction such as 2.5, cut to a whole one; ""
// when it is one, is not set, or data does not parse, which the load says.
func batchMaxProblem(data []byte) string {
	var raw struct {
		MergeQueue struct {
			BatchMax any `yaml:"batch_max"`
		} `yaml:"merge_queue"`
	}
	if yaml.Unmarshal(data, &raw) != nil {
		return ""
	}
	written := ""
	switch v := raw.MergeQueue.BatchMax.(type) {
	case nil, int, int64, uint64:
		return ""
	case float64:
		if v == math.Trunc(v) && !math.IsInf(v, 0) {
			return ""
		}
		written = strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		written = strconv.Quote(v)
	default:
		written = fmt.Sprint(v)
	}
	return fmt.Sprintf("merge_queue.batch_max %s is not a whole number; write the most enqueued stories one batch of the merge queue holds, such as %d, or remove it for %d", written, DefaultBatchMax, DefaultBatchMax)
}

// Errors are what is wrong with the merge queue's settings, one sentence
// each; none when they are valid.
func (q MergeQueue) Errors() []string {
	var errs []string
	if b := q.BatchMax; b != nil && *b < 1 {
		errs = append(errs, fmt.Sprintf("merge_queue.batch_max %d is not a number of one or more; write the most enqueued stories one batch of the merge queue holds, such as %d, or remove it for %d", *b, DefaultBatchMax, DefaultBatchMax))
	}
	if _, err := q.CoalesceDuration(); err != nil {
		errs = append(errs, err.Error())
	}
	return errs
}
