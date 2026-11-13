---
id: T-2868
type: task
nature: feature
title: flai kluster deployments and its records list each story namespace with its column, head, last deploy and checks run, who asked, and the count against the limit
status: done
parent: S-0663
owner: alex
created: 2026-11-10T08:24:39Z
updated: 2026-11-13T14:02:19Z
transitions:
  - to: ready
    at: 2026-11-13T13:42:41Z
    by: agent-S-0663
  - to: in-progress
    at: 2026-11-13T13:42:41Z
    by: agent-S-0663
  - to: done
    at: 2026-11-13T14:02:19Z
    by: agent-S-0663
stream: S-0663
tags: [cli]
touches: [flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, docs/users/flai-reference.md, docs/operators/settings.md, flai/internal/serve/kluster_requests.go]
usage:
  source: log
  seconds: 1178
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 126
      output: 694
      cache_read: 8839209
      cache_write: 198148
      cost: 3.3816
---
# T-2868 flai kluster deployments and its records list each story namespace with its column, head, last deploy and checks run, who asked, and the count against the limit

## Work

Make the list of deployments say what the page and the command both show. Start from what S-0660 and S-0661 built: the deploy a story's agent asks for while its story is in progress, the limit on story namespaces and its refusal, and the checks run against a namespace with who asked. Read the ADR S-0659 wrote and `design/system/kluster.md` § Commands (S-0577, S-0578) first.

- In `flai/internal/kluster/deploy.go`, `Deployments` gives each story namespace's row its story's column (in progress, review, or any other), whether its record's commit is the story branch's head, its last deploy and its last checks run, each with its outcome, its log, its commit, and who asked (the story's agent, the operator, or `flai serve`), as S-0660's and S-0661's records hold them.
- It also gives the number of story namespaces that live, the limit, and the requests refused for the limit, each with its story, who asked, and when. Where S-0660 refuses a request over the limit, in `flai/internal/serve/kluster.go` or beside it, keep each refusal under `.flai-cache/kluster/` so the list can read it. Keep the newest few, not every one.
- In `flai/cmd/kluster.go`, `flai kluster deployments` adds `COLUMN`, `HEAD`, `CHECKS`, and `BY` to its table, and a line under it with the count against the limit and the last refusal. `--json` gives the same as fields beside the records, keeping every field it gives today, so the dashboard and older readers still read it.
- Regenerate `docs/users/flai-reference.md` with `make flai-reference` when the command's help changes.

It waits for no task: it is the first layer. The host read, T-2869, reads what it adds.

## Done when

- `flai kluster deployments` and `--json` list a story in progress's namespace with its column, head, last deploy and checks run with their outcomes, logs, and who asked, and the count against the limit with the last refusal, in tests in `flai/internal/kluster/deploy_test.go` and `flai/cmd/kluster_test.go`.
- A refusal over the limit is kept and listed, in a test in `flai/internal/serve/kluster_test.go`.
- `flai test` passes on the changed paths, `TestReferenceIsCurrent` among them.

## Notes
