---
id: T-3473
type: task
nature: feature
title: The workflow file's keys raise flai.minimum on publish, and flai upgrade keeps a project's file without them
status: in-progress
parent: S-0806
owner: alex
created: 2026-11-16T23:50:39Z
updated: 2026-11-17T08:32:23Z
transitions:
  - to: ready
    at: 2026-11-17T08:32:22Z
    by: agent-S-0806
  - to: in-progress
    at: 2026-11-17T08:32:23Z
    by: agent-S-0806
stream: S-0806
tags: [flai]
touches: [flai/internal/workitem/front-matter-fields.txt, flai/internal/workitem/fields_test.go, flai/internal/upgrade/upgrade_test.go]
after: [T-3469]
---
# T-3473 The workflow file's keys raise flai.minimum on publish, and flai upgrade keeps a project's file without them

## Work

An older flai's strict parse refuses the new keys as an unknown field, in words that name no release. Make the refusal one that names the flai that reads them: list the workflow file's keys in `flai/internal/workitem/front-matter-fields.txt`, as lines for the file, a state, a transition, and `verify_at`, taken from the code, and keep them equal to the code in `TestFieldsFileIsTheCode`. The first publish from the release that carries this then raises `flai.minimum` (`release.RaiseMinimum`), and an older flai stops on the manifest with `manifest.TooOldError`, naming the release and the upgrade, before it reads the file.

When the shape T-3469 gives has no struct whose fields the test can list, say so in the narrative's `## Decisions` and ask on a thread on the story before choosing another way.

Add a case to `upgrade_test.go`: a project whose own `design/workflows/board.md` lacks the keys keeps it unchanged after `flai upgrade` from a template whose file has them, as `upgrade.kept` does today (ADR-0307).

It waits for T-3469, whose structs the fields file lists.

## Done when

- `front-matter-fields.txt` lists the workflow file's keys, and `TestFieldsFileIsTheCode` passes.
- The new `upgrade_test.go` case passes, and `flai test` on the changed paths passes.

## Notes
