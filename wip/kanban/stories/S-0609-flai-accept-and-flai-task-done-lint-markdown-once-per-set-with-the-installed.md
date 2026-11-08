---
id: S-0609
type: story
nature: improvement
title: flai accept and flai task done lint markdown once per set with the installed markdownlint, not five times through npx, saving about 4 s an acceptance
status: ready
parent: E-0033
owner: alex
created: 2026-11-07T22:41:56Z
updated: 2026-11-08T00:28:02Z
transitions:
  - to: ready
    at: 2026-11-08T00:23:35Z
    by: orchestrator
tags: [cli]
touches: [flai/cmd/accept_markdown.go, flai/cmd/accept.go, scripts/lint-md.sh, scripts/install-tools.sh, design/tech/ci.md, design/tech/node.md, docs/contributors/index.md, flai/internal/verify/fix.go, flai/internal/verify/fix_test.go, flai/cmd/branch.go, flai/cmd/accept_fix_test.go, flai/cmd/accept_uncommitted_test.go, flai/internal/preview/accept.go, design/system/server-performance.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 91
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 29
          output: 9960
          cache_read: 761818
          cache_write: 92071
          cost: 1.0882
    - kind: orchestrator
      seconds: 74
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 7
          output: 71
          cache_read: 755595
          cache_write: 24969
          cost: 0.1993
cost_of_delay:
  inputs:
    penalty_per_week: 15
    by: orchestrator
    at: 2026-11-08T00:03:47Z
  value: 15
  by: planner-S-0609
  at: 2026-11-08T00:10:42Z
forecast:
  duration: 1h
  delivery: 2026-11-08T02:44:00Z
  basis: "Its own forecast of 1h; 5th in the pull order with an in-progress limit of 5, behind S-0598, S-0599, S-0600, S-0605, S-0606, S-0601, S-0603, S-0607 and S-0604."
  by: flai
  at: 2026-11-08T00:28:02Z
finalized:
  by: orchestrator
  at: 2026-11-08T00:21:09Z
---
# S-0609 flai accept and flai task done lint markdown once per set with the installed markdownlint, not five times through npx, saving about 4 s an acceptance

## Goal

Cause 19 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `flai accept` runs `scripts/lint-md.sh` five times, one after another, each through `npx --yes markdownlint-cli2` (`flai/cmd/accept_markdown.go`), and checks again after each `--fix`. Each call takes 1.0 to 1.2 s, of which markdownlint itself is 0.27 to 0.36 s. `flai task done` runs it through `npx` twice. `accept` also starts 58 git processes, `git status --porcelain` six times, and lists every item for `EpicWalk`. The serve log has `accept.run` at p50 16.6 s, up to 42 s.

## Acceptance criteria

- [ ] `flai accept` lints the uncommitted set and the branch's set in one call each, runs the installed `markdownlint-cli2` without `npx` where it is installed, and runs no check after a `--fix` whose exit already says what remains.
- [ ] `flai accept` reads `git status --porcelain` once per state of the checkout it needs, not six times.
- [ ] An acceptance in a scratch copy of this repository takes at least 4 s less than before, refusing the same unfixed markdown as ADR-0140 says, and server-performance.md records the before and after.

## Tasks

- T-2636 scripts/lint-md.sh runs the installed markdownlint-cli2 at its pinned version, and install-tools.sh installs it
- T-2642 verify's fix answers what its fix command's exit says remains, so a caller needs no check after a --fix
- T-2649 flai accept lints the uncommitted set and the branch's set in one call each, with no check after a fix whose exit says what remains
- T-2658 flai accept reads git status --porcelain once per state of the checkout it needs, not six times
- T-2664 Measure an acceptance before and after in a scratch copy and record it in server-performance.md

## Notes

- `releaseDue`, `tellAccepted`, and `installAccepted`, the acceptance's last steps, are independent and could run at once; the saving is small but larger when an install runs.

### Planning

Touches, and where each came from:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/cmd/accept_markdown.go`, `flai/cmd/accept.go`, `scripts/lint-md.sh` | declared | The lint calls and the script they run |
| `flai/internal/verify/fix.go`, `flai/internal/verify/fix_test.go` | layout | `CheckoutFix` and `CheckoutUnfixed` run the tier; the fix must answer what remains |
| `flai/cmd/branch.go` | layout | The rebase's fix and `refuseUnfixedBranch`, the branch set's two calls |
| `flai/internal/preview/accept.go` | co-change (26%) and layout | The preview's `git status --porcelain` reads |
| `flai/cmd/accept_fix_test.go`, `flai/cmd/accept_uncommitted_test.go` | co-change and layout | The tests that count the lint runs and the status reads |
| `scripts/install-tools.sh` | layout | `markdownlint-cli2` is on no `PATH` on this host, so "installed" needs an install step |
| `design/tech/ci.md`, `design/tech/node.md`, `docs/contributors/index.md` | design | Where markdownlint and `install-tools.sh` are recorded |
| `design/system/server-performance.md` | design | Criterion 3: the before and after |

No folder touches. The suggestions `design/system/flai-cli.md` and `docs/users/flai*.md` were left out: no command's behaviour or flags change.

No `after`: the story reads no other open story's output.

Figures:

- Forecast: `flai forecast` gave 27 m (94 s per unit of size over 124 done improvement stories on claude-opus-5-5, size 17 from 3 criteria and 14 touches). Raised to 1 h because criterion 3 builds flai at two commits and times acceptances in a scratch copy, with and without markdownlint installed, which the size does not count. Delivery is `flai forecast`'s 2026-11-08T02:56Z, 9th in the pull order, moved by the same 33 m to 03:29Z.
- Cost of delay: 15 USD a week, as `flai cod` gives it from the input `penalty_per_week: 15` the orchestrator set on TH-0897, the figure TH-0888 recommends for cause 19 (about 4 s an acceptance). Kept unadjusted: no revenue rides on it, and it ranks below the epic's other causes, as their per-acceptance or per-request costs are larger.

Assumptions:

- `flai task done` gains from the installed markdownlint through T-2636 alone. Its lint calls are not cut, since no criterion asks for it.
- The `EpicWalk` listing and running `releaseDue`, `tellAccepted`, and `installAccepted` at once are left out, since no criterion asks for them.
- The merge queue's preflight (`flai/internal/mergequeue/preflight.go`) keeps its fix-then-check, since no criterion asks for it.
