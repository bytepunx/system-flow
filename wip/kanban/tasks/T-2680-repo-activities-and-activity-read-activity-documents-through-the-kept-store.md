---
id: T-2680
type: task
nature: improvement
title: Repo.Activities and Activity read activity documents through the kept store
status: backlog
parent: S-0602
owner: alex
created: 2026-11-07T23:50:29Z
updated: 2026-11-07T23:50:29Z
transitions: []
stream: S-0602
tags: [cli]
touches: [flai/internal/workitem/activity.go, flai/internal/workitem/activity_test.go]
after: [T-2660]
---
# T-2680 Repo.Activities and Activity read activity documents through the kept store

## Work

Move `Repo.Activity` and `Repo.Activities` in `flai/internal/workitem/activity.go` onto a package-level `Kept[*Activity]`, with `ReadActivity` as its parse and a deep copy of an activity and its entries as its clone. Keep the order of `ActivityKinds`, a missing document left out, and the error for a document whose kind is not its file's name. `AppendActivity` and `AppendRefusal` save through `atomicfile` under `lockActivity`, so the next read sees them by stat. `stats.get` reads these documents.

It waits for T-2660, whose store it uses.

## Done when

- `activity_test.go` shows that an entry appended, a document added, and one removed are each seen on the next `Activities`, and that a second read with nothing changed parses no file.
- It shows that two reads at once return equal activities.
- `flai test flai/internal/workitem/activity.go flai/internal/workitem/activity_test.go` passes.

## Notes
