---
id: T-3624
type: task
nature: improvement
title: The orchestrator's role file and the strategic-agents convention say when it keeps a removal and how it logs it
status: backlog
parent: S-0852
owner: alex
created: 2026-11-27T23:58:09Z
updated: 2026-11-27T23:58:09Z
transitions: []
stream: S-0852
tags: [flai, template]
touches: [design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/internal/roles/baseline/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md]
after: [T-3621, T-3622]
---
# T-3624 The orchestrator's role file and the strategic-agents convention say when it keeps a removal and how it logs it

## Work

In I-0415 the orchestrator could only open TH-1321 and recommend that the operator run `flai accept --keep-removal`. Tell it what it may do now, in the baseline role file `flai/internal/roles/baseline/orchestrator.md`, its copies `design/workflows/roles/orchestrator.md` and `template/root/design/workflows/roles/orchestrator.md`, and step 1 and step 2 of "Accepting a story" in `design/conventions/strategic-agents.md` and its template copy. Add a `setting:` entry for the new permission where the role file lists `accept_reviews`.

While the permission is on, the orchestrator reads each path the preview lists as `undoes_main`. It reads main's lines the branch removes beside the lines the branch adds there, and the story's narrative. It keeps a path only when the narrative's `## Decisions` names it as a rewrite and the branch's lines restate main's. Then it gives `keep_removals` to `accept`, or to `queue_add` while the queue runs, with a reason naming that entry. It logs each keep with `activity_log`. Otherwise, and while the permission is off, it leaves the path to the operator, as now. Change the baseline above the marker, since the baseline changes for every project.

It waits for T-3621 and T-3622, whose tool fields it names.

## Done when

- The three role files and the two convention copies say it, the copies match, and the role file's `setting:` lists the permission.
- `flai check --strict` passes, including the `roles` checks, and `flai test` passes on the changed paths.

## Notes
