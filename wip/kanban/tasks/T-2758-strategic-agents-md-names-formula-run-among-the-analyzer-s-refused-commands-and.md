---
id: T-2758
type: task
nature: improvement
title: strategic-agents.md names formula run among the analyzer's refused commands, and flai-cli.md lists flai checks and flai serve checks
status: in-progress
parent: S-0626
owner: alex
created: 2026-11-09T17:13:21Z
updated: 2026-11-09T17:27:18Z
transitions:
  - to: ready
    at: 2026-11-09T17:27:17Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T17:27:18Z
    by: system-flow
stream: S-0626
tags: [design, docs, flai]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md]
---
# T-2758 strategic-agents.md names formula run among the analyzer's refused commands, and flai-cli.md lists flai checks and flai serve checks

## Work

- `design/system/strategic-agents.md` § The analyzer's guard, the row at line 930: the guard refuses the analyzer `formula run` too. `cliAuthors` in `flai/internal/guard/guard.go:501-510` holds `"formula": {"run"}`, and the analyzer's rules at lines 521-533 refuse what it holds. Add `formula run` to the row.
- `design/system/flai-cli.md` § Commands: add rows for `flai checks run|status|cancel|tail` (`flai/cmd/checks.go:21-25`) and `flai serve checks set|timeout|clear|show` (`flai/cmd/serve_actions.go:1010-1114`), each saying what it does in one line, as the generated `docs/users/flai-reference.md` gives it.

Bump each document's `updated`. S-0619, S-0620, and S-0629 also touch `flai-cli.md`, and S-0620 touches `strategic-agents.md`. A trial merge judges the overlap.

It waits for no task: no other task changes these files.

## Done when

- The analyzer's guard row names `formula run`.
- § Commands has a row for `flai checks` and for `flai serve checks`, each with its subcommands.
- `flai check --strict` and the markdown lint pass on both files.

## Notes
