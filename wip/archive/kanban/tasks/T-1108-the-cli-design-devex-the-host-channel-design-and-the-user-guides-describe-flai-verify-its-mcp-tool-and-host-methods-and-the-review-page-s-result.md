---
id: T-1108
type: task
nature: improvement
title: The CLI design, devex, the host channel design, and the user guides describe flai verify, its MCP tool and host methods, and the review page's result
status: done
parent: S-0270
owner: alex
created: 2026-10-06T22:53:26Z
updated: 2026-10-07T04:45:04Z
transitions:
  - to: ready
    at: 2026-10-07T04:33:23Z
    by: agent-S-0270
  - to: in-progress
    at: 2026-10-07T04:33:24Z
    by: agent-S-0270
  - to: done
    at: 2026-10-07T04:45:04Z
    by: agent-S-0270
stream: S-0270
tags: [flai, docs]
touches: [design/system/flai-cli.md, design/system/devex.md, design/system/dashboard-host-channel.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md, design/system/flaiover-dashboard.md, design/system/project-manifest.md, design/system/agent-context.md, docs/operators/settings.md, design/system/conventions.md, design/system/strategic-agents.md, design/system/workflow.md, docs/operators/index.md]
after: [T-1075, T-1076, T-1082, T-1089, T-1096]
usage:
  source: log
  seconds: 700
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 175
      output: 73552
      cache_read: 12732549
      cache_write: 369556
      cost: 6.3244
---
# T-1108 The CLI design, devex, the host channel design, and the user guides describe flai verify, its MCP tool and host methods, and the review page's result

## Work

Criterion 5, and the documentation for criteria 2 to 4.

- `design/system/flai-cli.md`: add `flai verify <story> [--json] [--record-issues]` to the command table, the `verify` MCP tool, and the prompt change from T-1096.
- `design/system/devex.md`: in the test tiers and close-out rows, say that `flai verify` runs the tiers the diff selects and that the close-out calls it.
- `design/system/dashboard-host-channel.md`: add `verify.run`, gated by the `checks` host action, and `verify.status`.
- `docs/users/flai.md`: rewrite the close-out paragraphs near "closes out with `scripts/close-out.sh`" and the verifier row of the sub-agents table, and describe `flai verify` and its output.
- `docs/users/flai-reference.md`: regenerate with `make flai-reference`, not by hand.
- `docs/users/flaiover.md`: describe the verification result on the review page.
- Waits for every other task: it describes what they built.

## Done when

- [ ] Each document describes `flai verify`, its MCP tool and host methods, the close-out's call, and the review page's result as they were built.
- [ ] `docs/users/flai-reference.md` is the output of `make flai-reference`.
- [ ] `flai check --strict` and the markdown lint pass.

## Notes

Drafted by the planner.
