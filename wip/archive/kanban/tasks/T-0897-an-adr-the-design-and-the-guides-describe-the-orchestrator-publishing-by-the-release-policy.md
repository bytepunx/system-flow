---
id: T-0897
type: task
nature: feature
title: An ADR, the design, and the guides describe the orchestrator publishing by the release policy
status: done
parent: S-0222
owner: alex
created: 2026-10-05T04:47:15Z
updated: 2026-10-06T12:38:02Z
transitions:
  - to: ready
    at: 2026-10-06T12:33:52Z
    by: agent-S-0222
  - to: in-progress
    at: 2026-10-06T12:33:53Z
    by: agent-S-0222
  - to: done
    at: 2026-10-06T12:38:02Z
    by: agent-S-0222
stream: S-0222
tags: [flai]
touches: [design/adrs, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md]
after: [T-0891, T-0895]
usage:
  source: log
  seconds: 249
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 88
      output: 24594
      cache_read: 3977885
      cache_write: 198957
      cost: 2.6599
---
# T-0897 An ADR, the design, and the guides describe the orchestrator publishing by the release policy

## Work

ADR-0067 says agents publish only when the operator asks. Add an ADR that refines it with `flai adr new`: the operator's `publish` permission is the asking. The ADR covers:

- the orchestrator publishes only through `release_publish`, and only under its release policy
- `whole_epics` holds a batch back under every policy
- a refusal is logged and goes to the operator on a thread, and is never worked around

If S-0218's ADR already decides this, link it from the design instead and write none. Say which in the narrative's `## Decisions`.

Describe publishing in `design/system/strategic-agents.md`'s orchestrator section: when it evaluates, what each policy publishes, what `release_publish` refuses, what is logged, and the thread on a refusal. Add `release_publish` to the MCP tools in `design/system/flai-cli.md` and `docs/users/flai.md`. Say in `docs/operators/index.md`, under the push host action, that the orchestrator's `publish` permission needs the `push` action on too.

This task waits for T-0891 and T-0895, whose tool and prompt it describes.

## Done when

- the ADR is accepted and listed in `design/adrs/README.md`, or the design links S-0218's
- `strategic-agents.md`, `flai-cli.md`, `flai.md`, and the operators guide describe the tool, the policies, and the refusal as built
- `scripts/lint-md.sh` and `flai check --strict` report nothing on the changed files

## Notes
