---
id: T-2306
type: task
nature: feature
title: The strategic-agents convention and its template copy, flai-cli.md, strategic-agents.md, project-manifest.md, the roles README, and the user and operator docs say the orchestrator runs for a reason and ends
status: done
parent: S-0536
owner: alex
created: 2026-10-28T21:46:34Z
updated: 2026-10-29T08:30:23Z
transitions:
  - to: ready
    at: 2026-10-29T08:02:28Z
    by: agent-S-0536
  - to: in-progress
    at: 2026-10-29T08:02:29Z
    by: agent-S-0536
  - to: done
    at: 2026-10-29T08:30:23Z
    by: agent-S-0536
stream: S-0536
tags: [docs, conventions, template, orchestrator]
touches: [design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/system/flai-cli.md, design/system/strategic-agents.md, design/system/project-manifest.md, design/system/agent-narrative.md, design/workflows/roles/README.md, template/root/design/workflows/roles/README.md, docs/users/flai.md, docs/operators/index.md, template/CHANGELOG.md, docs/users/flai-reference.md, docs/users/conventions.md, flai/cmd/role.go, docs/users/flaiover.md, docs/operators/settings.md]
after: [T-2303, T-2304, T-2305]
usage:
  source: log
  seconds: 1674
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 259
      output: 97826
      cache_read: 18946092
      cache_write: 602161
      cost: 9.1497
---
# T-2306 The strategic-agents convention and its template copy, flai-cli.md, strategic-agents.md, project-manifest.md, the roles README, and the user and operator docs say the orchestrator runs for a reason and ends

## Work

ADR-0209 §1 to §4 and §7, criterion 5. Change every document that says the orchestrator runs as a loop, holds `wait_for_events` between decisions, or does not end:

- `design/conventions/strategic-agents.md` and `template/root/design/conventions/strategic-agents.md`, § As the orchestrator: "Wait on `wait_for_events` between decisions" becomes act on each reason the start prompt names, log each action, then end; and § Logging your activity, where a run "spans activities, as the orchestrator's does". Edit the baseline in both, since the template ships it.
- `design/system/strategic-agents.md`, § The orchestrator: Starting and stopping it, How it runs, What it is told, and The run and how it ends: the reason queue, `orchestration.coalesce`, `orchestration.run_timeout`, `orchestration.schedule`, `end: true` with `reasons done`, and the reasons in the journal.
- `design/system/flai-cli.md`: `flai serve`'s orchestrator and `wait_for_events`' end.
- `design/system/project-manifest.md`: the three keys and their defaults.
- `design/system/agent-narrative.md`, § Strategic agents' activity documents, where it says the orchestrator's run spans activities.
- `design/workflows/roles/README.md` and its template copy: `continuous` is no longer a trigger (ADR-0209), and the orchestrator's triggers are events.
- `docs/users/flai.md` and `docs/operators/index.md`: what the operator sees of the orchestrator's runs and reasons, and the three settings.
- `template/CHANGELOG.md`: the changed convention, role file, and definition.

Read each section with `doc_get` before changing it. Waits for T-2303, T-2304, and T-2305, so that the documents describe what was built.

## Done when

- No document under `design/`, `docs/`, or `template/` says the orchestrator holds `wait_for_events` without ending, or that `continuous` is a trigger, other than ADR-0087 and ADR-0207, which are immutable.
- Each document above names the reason queue, the coalescing window, the end, and the timeout where it describes the orchestrator's run.
- `flai check --strict` and the markdown lint pass on the changed paths.

## Notes
