---
id: T-2073
type: task
nature: feature
title: A role file names its charge rule, item, named, or issues, the baseline role files carry theirs, and flai check reports an unknown rule
status: done
parent: S-0425
owner: alex
created: 2026-10-23T15:34:05Z
updated: 2026-11-02T19:17:57Z
transitions:
  - to: ready
    at: 2026-11-02T19:02:09Z
    by: agent-S-0425
  - to: in-progress
    at: 2026-11-02T19:02:10Z
    by: agent-S-0425
  - to: done
    at: 2026-11-02T19:17:57Z
    by: agent-S-0425
stream: S-0425
tags: [cli, template]
touches: [flai/internal/roles/roles.go, flai/internal/roles/roles_test.go]
after: [T-1790]
usage:
  source: log
  seconds: 947
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 2835
      cache_read: 508632
      cache_write: 13409
      cost: 0.2383
---
# T-2073 A role file names its charge rule, item, named, or issues, the baseline role files carry theirs, and flai check reports an unknown rule

## Work

S-0414 already gave the role file its `charge` field: `flai/internal/roles` parses it with `item`, `named`, `issues`, or `none`, defaulting to `none` (`Charges`, `roles.go`), `flai check` reports a `charge` that is none of them (`roles.file`, tested in `roles_test.go`'s "bad charge"), and the baseline role files here, in the template, and in `flai/internal/roles/baseline` carry `item` on the planner, `named` on the orchestrator, `issues` on the analyzer, and `none` on the curator. What is left, as the ADR of T-1790 names it:

- `flai check` reports `roles.file` on `charge` for `charge: item` on a role whose `item` is `none`: such a run is started for no item, so the rule could charge nothing.

Waits for T-1790, which fixes the rule. Runs beside T-1791 and T-1793 in its layer: they share no path.

## Done when

`flai test` on `flai/internal/roles` passes, with a test for `charge: item` on a role that takes no item, and `flai check --strict` reports nothing on the role files.

## Notes

Narrowed on pull, 2026-11-02: the field, the unknown-rule check, and the baseline files were done by S-0414 (T-1722), so `flai/internal/check`, the six role files, and `template/CHANGELOG.md` left its touches.
