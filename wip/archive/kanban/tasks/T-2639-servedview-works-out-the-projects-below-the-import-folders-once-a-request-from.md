---
id: T-2639
type: task
nature: improvement
title: servedView works out the projects below the import folders once a request, from the running flai serve's status when it runs, and a test with three served projects under two import folders shows the same answer
status: done
parent: S-0605
owner: alex
created: 2026-11-07T23:48:41Z
updated: 2026-11-08T00:34:39Z
transitions:
  - to: ready
    at: 2026-11-08T00:14:53Z
    by: agent-S-0605
  - to: in-progress
    at: 2026-11-08T00:14:58Z
    by: agent-S-0605
  - to: done
    at: 2026-11-08T00:34:39Z
    by: agent-S-0605
stream: S-0605
tags: [cli]
touches: [flai/cmd/serve_project.go, flai/cmd/serve_actions.go, flai/internal/serve/below.go, flai/internal/serve/below_test.go, flai/cmd/serve_project_test.go]
usage:
  source: log
  seconds: 1181
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 58
      output: 14890
      cache_read: 7998401
      cache_write: 192956
      cost: 3.1575
---
# T-2639 servedView works out the projects below the import folders once a request, from the running flai serve's status when it runs, and a test with three served projects under two import folders shows the same answer

## Work

- `importRootOf` (`flai/cmd/serve_project.go`) runs `serve.FindBelow` over every import folder, with a manifest load per project found, only to ask whether one root is among them. Replace that with a helper in `flai/internal/serve/below.go`, such as `ImportRootOf(roots []string, root string) string`, that answers the same from the path alone: for each import folder, in order, that holds `root`, whether `repositories` would reach `root` from it (no hidden or `skipScanning` component on the way, no symbolic link followed, within `candidateDepth`, no repository above `root` at or below the folder, `root` itself a repository) and `root` has a `system-flow.yaml`. A handful of stats per project, no walk, no manifest load. `servedBelow` and every other caller use it.
- `servedView` (`flai/cmd/serve_actions.go`) does not show the candidates, so it lists the served projects without `serve.FindCandidates`' walk: split `listServedProjects` so that `flai serve project list` still gets candidates and `servedView` does not ask for them. It loads the configuration once a call.
- So while flai serve runs, `servedView` walks no folder: the unserved list comes from its status as now. When it does not run, `belowImportRoots`' one `FindBelow` is the only walk.
- Leave what `settings.get`'s `host.projects` answers as it is: `running`, `served` with each `below`, and `unserved`.
- Tests: in `flai/internal/serve/below_test.go`, the helper agrees with `repositories` plus the manifest check on a tree with two import folders, a hidden folder, a `node_modules` folder, a repository nested in a repository, one past `candidateDepth`, a repository with no manifest, and a symbolic link. In `flai/cmd/serve_project_test.go`, two import folders, three served projects below them (two below one, one below the other), and a registered project below neither: each project's `below` and the `unserved` list are what they were before the change.
- Waits for no task: it is the first layer.

## Done when

- `servedView` walks no folder for `below`, and none at all while flai serve runs.
- The new tests pass, and the existing `servedView`, `flai serve project`, and `serve` package tests pass unchanged.

## Notes

Rewritten by agent-S-0605 from the planner's draft: the status cannot answer `below` for a registered project below an import folder, since `serve.Place` leaves registered projects out of `ImportProjects` and `Unserved`, so the path check replaces the walk instead.
