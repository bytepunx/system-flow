---
id: T-3469
type: task
nature: feature
title: The workflow package parses verify on a transition, runs on a state, and verify_at, and lists the suites a file names
status: in-progress
parent: S-0806
owner: alex
created: 2026-11-16T23:50:26Z
updated: 2026-11-17T08:17:47Z
transitions:
  - to: ready
    at: 2026-11-17T08:17:41Z
    by: agent-S-0806
  - to: in-progress
    at: 2026-11-17T08:17:47Z
    by: agent-S-0806
stream: S-0806
tags: [flai]
touches: [flai/internal/workflow/workflow.go, flai/internal/workflow/workflow_test.go, flai/internal/workflow/baseline.go, flai/internal/workflow/suites.go, flai/internal/workflow/suites_test.go]
---
# T-3469 The workflow package parses verify on a transition, runs on a state, and verify_at, and lists the suites a file names

## Work

Add the three keys of ADR-0312 to the workflow file's YAML shapes in `flai/internal/workflow/workflow.go`, keeping the strict decoding (`yaml.Strict()`) that refuses any other key:

- `verify:` on a transition: a list of suites or `all`, with a severity, `refuse` unless the file says `warn`, as ADR-0312 §2 says. Follow the shape the accepted ADR gives; when it leaves the severity's place open, take a mapping `{suites: <list or all>, severity: warn}` beside the bare list or `all`, and record the choice in the narrative's `## Decisions`.
- `runs:` on a state: `suites:`, a list or `all`, and `kluster:`, any of `deploy` and `checks` (§4).
- `verify_at:` at the top level, with `batch`, `publish`, and `ci`, each a list of suites or `all` (§6).

Refuse at parse time, as `workflow.file`, a suite name that is not kebab-case, a `kluster:` entry other than `deploy` and `checks`, and a `verify_at` key other than the three. Make `Clone` in `baseline.go` deep-copy the new lists. Add `suites.go`: the suites a file names, each with where it is named (a transition, a state, or a point), and whether the file names `all` anywhere, for the check and for the stories after this one.

This task waits for no other task: the rest of the story builds on these shapes.

## Done when

- `workflow_test.go` parses each key in every form it takes, and refuses an unknown key inside each and an unknown top-level key, as before.
- `suites_test.go` lists the suites of a file that names some, `all`, and none.
- `TestBaselineIsACopy` still passes, and `flai test` on the changed paths passes.

## Notes
