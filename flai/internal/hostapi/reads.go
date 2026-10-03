package hostapi

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/metrics"
	"github.com/bytepunx/system-flow/flai/internal/perf"
	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/statsread"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The dashboard's reads that a command also prints (S-0159): until then
// each started flai and read its --json, and the dashboard asks for some of
// them at every change. They are answered here, in flai serve's process,
// from the code the command prints with, and with what flai's answer was:
// the command's JSON as data, the warnings it logged, and its failure read
// as a failure of the process was, so the dashboard sees no difference.

// work is one read: what the command computes, given the project, the
// runner its git goes through, and the logger its warnings go to.
type work func(r execx.Runner, repo *workitem.Repo, log *slog.Logger) (any, error)

// answer runs a read under the phase named for it and returns what the
// command would have answered.
func answer(ctx context.Context, p channel.Project, phase string, w work) (any, *channel.Error) {
	done := perf.Track(ctx, "repo.open")
	repo, err := workitem.Open(p.Root)
	done()
	if err != nil {
		return nil, failed(err)
	}
	r := execx.Timed(ctx, execx.System{})
	repo.Git = r
	var events bytes.Buffer
	log := slog.New(slog.NewJSONHandler(&events, &slog.HandlerOptions{Level: slog.LevelWarn}))
	done = perf.Track(ctx, phase)
	v, err := w(r, repo, log)
	done()
	ran := Ran{Events: logged(&events)}
	if err != nil {
		ran.Exit = 1
		ran.Events = append(ran.Events, map[string]any{"level": "FATAL", "msg": "command failed", "err": err.Error()})
		return outcome(ran, nil, nil)
	}
	if ran.Stdout, err = json.Marshal(v); err != nil {
		return nil, failed(err)
	}
	return outcome(ran, nil, nil)
}

// logged reads back the events a read logged, as ExecRunner reads a flai's.
func logged(b *bytes.Buffer) []map[string]any {
	var out []map[string]any
	sc := bufio.NewScanner(b)
	sc.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for sc.Scan() {
		var ev map[string]any
		if json.Unmarshal(sc.Bytes(), &ev) == nil {
			out = append(out, ev)
		}
	}
	return out
}

// idOnly decodes {"id": ...} and checks it is an item's ID.
func idOnly(raw json.RawMessage) (string, *channel.Error) {
	in, e := decode[struct {
		ID string `json:"id"`
	}](raw)
	if e != nil {
		return "", e
	}
	return in.ID, needID(in.ID)
}

func readMethods(now func() time.Time) map[string]channel.Method {
	return map[string]channel.Method{
		// item.show: what flai edit --show prints, an item's own words and
		// the hash an edit gives back (S-0085).
		"item.show": func(ctx context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			id, e := idOnly(raw)
			if e != nil {
				return nil, e
			}
			return answer(ctx, p, "item.show", func(_ execx.Runner, repo *workitem.Repo, _ *slog.Logger) (any, error) {
				return itemedit.Show(repo, id)
			})
		},

		// item.move.preview: what cancelling an item would take with it, as
		// flai move <id> cancelled --dry-run prints it (ADR-0028).
		"item.move.preview": func(ctx context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			id, e := idOnly(raw)
			if e != nil {
				return nil, e
			}
			by := owner(p)
			return answer(ctx, p, "cancel.preview", func(r execx.Runner, repo *workitem.Repo, _ *slog.Logger) (any, error) {
				it, err := repo.Get(id)
				if err != nil {
					return nil, err
				}
				// flai validates the move before it lists anything, and a
				// cancellation needs a reason
				return preview.Cancel(r, repo, it, by, "preview", now())
			})
		},

		// accept.preview: what would block an acceptance, as flai accept
		// --dry-run prints it.
		"accept.preview": func(ctx context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			id, e := idOnly(raw)
			if e != nil {
				return nil, e
			}
			by := owner(p)
			return answer(ctx, p, "accept.preview", func(r execx.Runner, repo *workitem.Repo, log *slog.Logger) (any, error) {
				it, err := repo.Get(id)
				if err != nil {
					return nil, err
				}
				return preview.Accept(r, repo, it, by, now(), log)
			})
		},

		// stream.diff: a story branch against the main branch, as flai
		// stream diff prints it (S-0041).
		"stream.diff": func(ctx context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			id, e := idOnly(raw)
			if e != nil {
				return nil, e
			}
			return answer(ctx, p, "stream.diff", func(r execx.Runner, repo *workitem.Repo, _ *slog.Logger) (any, error) {
				it, err := repo.Get(id)
				if err != nil {
					return nil, err
				}
				return storygit.StoryDiff(r, repo, it.ID)
			})
		},

		// stats.get: flow metrics, as flai stats prints them.
		"stats.get": func(ctx context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			in, e := decode[struct {
				Since  string `json:"since"`
				Type   string `json:"type"`
				By     string `json:"by"`
				Bucket string `json:"bucket"`
			}](raw)
			if e != nil {
				return nil, e
			}
			if in.Since != "" && !sinceValue.MatchString(in.Since) {
				return nil, bad("since is a number and d, w, or h")
			}
			if !itemType[in.Type] {
				return nil, bad("type must be epic, story, or task")
			}
			if in.By != "" && in.By != "nature" && in.By != "type" && in.By != "parent" {
				return nil, bad("by must be nature, type, or parent")
			}
			if in.Bucket != "" && in.Bucket != metrics.BucketHour && in.Bucket != metrics.BucketDay && in.Bucket != metrics.BucketWeek {
				return nil, bad("bucket must be hour, day, or week")
			}
			return answer(ctx, p, "stats.compute", func(r execx.Runner, repo *workitem.Repo, log *slog.Logger) (any, error) {
				since := in.Since
				if since == "" {
					since = metrics.DefaultWindow
				}
				window, err := metrics.ParseWindow(since)
				if err != nil {
					return nil, err
				}
				if err := metrics.CheckBucket(in.Bucket, window); err != nil {
					return nil, err
				}
				items, opt, err := statsread.Read(r, repo, log)
				if err != nil {
					return nil, err
				}
				opt.Now, opt.Since, opt.Type, opt.By, opt.Bucket = now(), window, in.Type, in.By, in.Bucket
				return metrics.Compute(items, opt), nil
			})
		},

		// publish.preview: everything release.Pending would release, without
		// changing anything, so the board can show it before the operator
		// asks for it (S-0087).
		"publish.preview": func(ctx context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, e
			}
			return answer(ctx, p, "release.pending", func(r execx.Runner, repo *workitem.Repo, _ *slog.Logger) (any, error) {
				return preview.Publish(r, repo)
			})
		},
	}
}
