---
id: S-0436
type: story
nature: feature
title: flai reads and saves a text file that is not markdown through file.get and file.save, and a token save writes the stylesheet with it
status: in-progress
parent: E-0026
owner: alex
created: 2026-10-16T21:31:07Z
updated: 2026-10-16T22:33:17Z
transitions:
  - to: ready
    at: 2026-10-16T21:52:18Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-16T21:52:24Z
    by: agent-S-0436
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
  seconds: 2713
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-10-16
      ceremony: 2
      hand_edits: 2
      work: 42
  models:
    - model: claude-opus-5-5
      input: 322
      output: 1516
      cache_read: 14557281
      cache_write: 654543
      cost: 7.0657
  strategic:
    - kind: orchestrator
      seconds: 381
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 46
          output: 804
          cache_read: 5798956
          cache_write: 14409
          cost: 1.4575
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

- [ ] `flai file show <path> [--story] [--json]` prints a file's raw content and its SHA-256, and `flai file save <path> --hash <h> [--story] [--message]` saves the content it reads on standard input. Both reach only a file under the manifest's design or docs folder, or the design token file or stylesheet the manifest names. Both refuse any other path, the wip folder, a markdown file, a binary file, and a file over the size limit this story sets, each with its reason.
- [ ] A save follows `doc.save`'s rules: a stale hash is a conflict, exit 3, with the current content and a diff; `flai check` runs before and after, and a finding the content introduces puts the old file back and refuses the save, exit 4; the save commits that path on the main checkout's branch as `docs:` unless `dashboard.autocommit` is false; with `--story` it reads, checks, and commits in the story's worktree, refused as ADR-0145 refuses.
- [ ] Whenever flai writes the token file, by `file.save` or by a `flai design` command, it writes the stylesheet the manifest names, generated from the new tokens, checks with both in place, and commits both paths in one commit, on main or on the story's branch (ADR-0150).
- [ ] The host methods `file.get` `{path, story?}` and `file.save` `{path, content, hash, story?, message?}` run the two commands, map exit 3 to Conflict and exit 4 to Refused, and are in the dashboard's `REQUIRED_METHODS`, so the banner names an older flai. The Design page's token save at `/api/design/tokens` then saves and commits on the host, shown by a test against a project with a token file and a stylesheet.
- [ ] `design/system/flai-cli.md`, `docs/users/flai.md` with `flai-reference.md` regenerated, `docs/users/flaiover.md`, and `design/system/design-system-tool.md` say the pair is built, what it reaches, and the size limit.

## Tasks
- T-1849 flai/internal/docedit reads and saves a non-markdown file by doc.save's rules, and writes the stylesheet with the token file
- T-1850 flai file show and flai file save, with --story, and flai design writes the stylesheet whenever it writes the token file
- T-1851 The host methods file.get and file.save run flai file show and save, and the dashboard requires them
- T-1852 The CLI design, the user docs, and the design system tool's design say file.get and file.save are built, what they reach, and the size limit

## Notes

- Drafted by the planner for E-0026. It belongs to E-0020's decisions (`design/system/design-system-tool.md` § file.get and file.save), but E-0020 is done, so it is planned here, where outcome 3 depends on it.
- ADR-0147 and ADR-0150 are accepted; this story needs no new ADR. It sets the size limit ADR-0147 leaves to it.

### Planning

- **Touches.** The story declared none; every touch comes from its four tasks, file by file.
  - **From design:** `design-system-tool.md` § file.get and file.save names `flai/cmd/file.go`, `flai/cmd/root.go`, `flai/internal/docedit/`, `flai/internal/hostapi/writes.go`, `writes_test.go`, `flai-cli.md`, `flai.md`, and the reference. `docedit/` is narrowed to a new `file.go` beside `docedit.go`.
  - **From layout:** `flai/internal/hostapi/files.go` for the read, merged in `hostapi.go` as `design.go`'s table is; `flai/cmd/design.go` for criterion 3; `flaiover/src/lib/server/agent.ts` for `REQUIRED_METHODS`. No `flai/cmd/file.go` or `file.get` exists in the tree today, and `repo.ts` already calls both.
  - **Co-change, not added:** `design/system/flaiover-dashboard.md` (25%), since the dashboard's routes do not change; `docs/operators/index.md` and `settings.md` (17%, 16%), since the pair adds no setting or host action.
  - No folder touch.
- **Forecast.** Kept at flai's 37m: size 23, at 96 s a unit over 61 done large feature stories. Delivery is flai's play-out, 43rd in the pull order.
- **Cost of delay.** 50.92 USD a week, flai's share of E-0026's 600 USD penalty by forecast duration, kept. It also unblocks S-0383's token save, which answers 501 today; E-0020 has no inputs to price that.
- The forecast carries `by: system-flow`: the CLI wrote it under the host's `FLAI_AGENT`. The planner wrote it.
