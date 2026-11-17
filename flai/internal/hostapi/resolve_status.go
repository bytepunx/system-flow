package hostapi

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/closeout"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// resolveStatusAnswer is resolve.status's answer: the story, and its last
// close-out or rebase and resolve job, left out when it has had none.
type resolveStatusAnswer struct {
	ID  string           `json:"id"`
	Job *resolveJobState `json:"job,omitempty"`
}

// resolveJobState is a story's last job as resolve.status answers it: the
// job's fields, its step, tier, steps, and log as closeout.Status reads
// them, and whether it is running.
type resolveJobState struct {
	closeout.JobStatus
	Running bool `json:"running"`
}

// resolveStatus answers resolve.status in flai serve's process (S-0804):
// story's last close-out or rebase and resolve job under the project's
// .flai-cache/verify/jobs, running or ended, as closeout.Status reads it
// with its progress from its log, which the review page follows while a
// Rebase and resolve runs and shows once it ended; with no job, the story
// alone. It runs nothing and records nothing.
func resolveStatus(ctx context.Context, _ Host, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
	id, e := verifyStatusOf(raw)
	if e != nil {
		return nil, e
	}
	return answer(ctx, p, "resolve.status", func(_ execx.Runner, repo *workitem.Repo, _ *slog.Logger) (any, error) {
		if it, err := repo.Get(id); err == nil {
			id = it.ID
		}
		st, ok, err := closeout.Status(repo.CacheDir(), id)
		if err != nil || !ok {
			return resolveStatusAnswer{ID: id}, err
		}
		return resolveStatusAnswer{ID: id, Job: &resolveJobState{JobStatus: st, Running: st.Running()}}, nil
	})
}
