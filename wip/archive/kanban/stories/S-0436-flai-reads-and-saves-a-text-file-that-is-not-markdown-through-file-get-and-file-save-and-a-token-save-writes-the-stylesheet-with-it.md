---
id: S-0436
type: story
nature: feature
title: flai reads and saves a text file that is not markdown through file.get and file.save, and a token save writes the stylesheet with it
status: done
parent: E-0026
owner: alex
created: 2026-10-16T21:31:07Z
updated: 2026-10-17T00:45:45Z
transitions:
  - to: ready
    at: 2026-10-16T21:52:18Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-16T21:52:24Z
    by: agent-S-0436
  - to: review
    at: 2026-10-16T23:07:51Z
    by: agent-S-0436
  - to: done
    at: 2026-10-17T00:45:45Z
    by: alex
tags: [cli, dashboard]
topics: [design, cli, dashboard]
touches: [flai/internal/docedit/file.go, flai/internal/docedit/file_test.go, flai/cmd/file.go, flai/cmd/file_test.go, flai/cmd/root.go, flai/cmd/design.go, flai/cmd/design_test.go, docs/users/flai-reference.md, flai/internal/hostapi/files_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/server/agent.ts, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md, design/system/design-system-tool.md, flai/internal/docedit/docedit.go, docs/operators/settings.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4821
  priced_by: harness
  turns:
    - day: 2026-10-16
      ceremony: 5
      hand_edits: 3
      work: 63
  models:
    - model: claude-opus-5-5
      input: 386
      output: 151846
      cache_read: 21554641
      cache_write: 710629
      cost: 11.6629
  strategic:
    - kind: orchestrator
      seconds: 456
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 50
          output: 816
          cache_read: 6577235
          cache_write: 18820
          cost: 1.6537
        - model: claude-sonnet-5-5
          input: 12
          output: 61
          cache_read: 206057
          cache_write: 65408
          cost: 0.2222
cost_of_delay:
  value: 50.92
  by: planner-E-0026
  at: 2026-10-16T21:49:57Z
forecast:
  duration: 37m
  delivery: 2026-10-17T19:13:00Z
  basis: "flai's 37m from 96 s a unit over 61 done large feature stories on claude-opus-5-5, times size 23; delivery is flai's play-out of the pull order, 43rd."
  by: system-flow
  at: 2026-10-16T21:42:55Z
finalized:
  by: orchestrator
  at: 2026-10-16T21:51:56Z
---
# S-0436 flai reads and saves a text file that is not markdown through file.get and file.save, and a token save writes the stylesheet with it

## Goal

The Design page saves the operator's token changes, so that a change to the design system reaches the components this epic renders. ADR-0147 decided two host methods for files that are not markdown: `file.get` and `file.save`, with `flai file show` and `flai file save`. ADR-0150 decided that a token save writes the generated stylesheet in the same commit. S-0383 built the dashboard's half and answers 501 until flai has them. S-0383's notes leave the flai side to a story of its own, which was never drafted, and S-0385's verdict leaves the stylesheet write to it. This is that story. With it the Design page's save works end to end, and the examples story of this epic has a save to build on.

## Acceptance criteria

- [x] `flai file show <path> [--story] [--json]` prints a file's raw content and its SHA-256, and `flai file save <path> --hash <h> [--story] [--message]` saves the content it reads on standard input. Both reach only a file under the manifest's design or docs folder, or the design token file or stylesheet the manifest names. Both refuse any other path, the wip folder, a markdown file, a binary file, and a file over the size limit this story sets, each with its reason.
- [x] A save follows `doc.save`'s rules: a stale hash is a conflict, exit 3, with the current content and a diff; `flai check` runs before and after, and a finding the content introduces puts the old file back and refuses the save, exit 4; the save commits that path on the main checkout's branch as `docs:` unless `dashboard.autocommit` is false; with `--story` it reads, checks, and commits in the story's worktree, refused as ADR-0145 refuses.
- [x] Whenever flai writes the token file, by `file.save` or by a `flai design` command, it writes the stylesheet the manifest names, generated from the new tokens, checks with both in place, and commits both paths in one commit, on main or on the story's branch (ADR-0150).
- [x] The host methods `file.get` `{path, story?}` and `file.save` `{path, content, hash, story?, message?}` run the two commands, map exit 3 to Conflict and exit 4 to Refused, and are in the dashboard's `REQUIRED_METHODS`, so the banner names an older flai. The Design page's token save at `/api/design/tokens` then saves and commits on the host, shown by a test against a project with a token file and a stylesheet.
- [x] `design/system/flai-cli.md`, `docs/users/flai.md` with `flai-reference.md` regenerated, `docs/users/flaiover.md`, and `design/system/design-system-tool.md` say the pair is built, what it reaches, and the size limit.

## Tasks
- T-1849 flai/internal/docedit reads and saves a non-markdown file by doc.save's rules, and writes the stylesheet with the token file
- T-1850 flai file show and flai file save, with --story, and flai design writes the stylesheet whenever it writes the token file
- T-1851 The host methods file.get and file.save run flai file show and save, and the dashboard requires them
- T-1852 The CLI design, the user docs, and the design system tool's design say file.get and file.save are built, what they reach, and the size limit

## Notes

- Drafted by the planner for E-0026. It belongs to E-0020's decisions (`design/system/design-system-tool.md` § file.get and file.save), but E-0020 is done, so it is planned here, where outcome 3 depends on it.
- ADR-0147 and ADR-0150 are accepted; this story needs no new ADR. It sets the size limit ADR-0147 leaves to it.
- The size limit is 1 MiB (`docedit.MaxFileSize`).
- Criterion 3 is ticked with this reading: `file save` commits the token file and the stylesheet in one commit. `flai design init` is the only `flai design` command that writes the token file. It writes the stylesheet with it through the same generator and, as its help has always said, commits nothing. The two files reach one commit when the caller commits them.

### Planning

- **Touches.** The story declared none; every touch comes from its four tasks, file by file.
  - **From design:** `design-system-tool.md` § file.get and file.save names `flai/cmd/file.go`, `flai/cmd/root.go`, `flai/internal/docedit/`, `flai/internal/hostapi/writes.go`, `writes_test.go`, `flai-cli.md`, `flai.md`, and the reference. `docedit/` is narrowed to a new `file.go` beside `docedit.go`.
  - **From layout:** `flai/internal/hostapi/files.go` for the read, merged in `hostapi.go` as `design.go`'s table is; `flai/cmd/design.go` for criterion 3; `flaiover/src/lib/server/agent.ts` for `REQUIRED_METHODS`. No `flai/cmd/file.go` or `file.get` exists in the tree today, and `repo.ts` already calls both.
  - **Co-change, not added:** `design/system/flaiover-dashboard.md` (25%), since the dashboard's routes do not change; `docs/operators/index.md` and `settings.md` (17%, 16%), since the pair adds no setting or host action.
  - No folder touch.
- **Forecast.** Kept at flai's 37m: size 23, at 96 s a unit over 61 done large feature stories. Delivery is flai's play-out, 43rd in the pull order.
- **Cost of delay.** 50.92 USD a week, flai's share of E-0026's 600 USD penalty by forecast duration, kept. It also unblocks S-0383's token save, which answers 501 today; E-0020 has no inputs to price that.
- The forecast carries `by: system-flow`: the CLI wrote it under the host's `FLAI_AGENT`. The planner wrote it.
