---
id: T-3653
type: task
nature: remediation
title: StoryLeft reads the kinds of the story's project's workflow, so a dev server stays up in a state the workflow adds
status: backlog
parent: S-0851
owner: alex
created: 2026-11-28T00:03:45Z
updated: 2026-11-28T00:03:45Z
transitions: []
stream: S-0851
tags: [flai, workflow, components]
touches: [flai/internal/serve/components.go, flai/internal/serve/components_test.go, design/system/component-library.md]
---
# T-3653 StoryLeft reads the kinds of the story's project's workflow, so a dev server stays up in a state the workflow adds

## Work

I-0406's second read is `StoryLeft` in `flai/internal/serve/components.go`. It sets `wf := baselineFlow`, and its comment says the story is all a look reads of the project. I-0402 is the defect that read causes. S-0756's end-to-end test adds `testing`, a state of the active kind, between in-progress and review. `StoryLeft` answers true for a story in `testing`, because `isUnderway` on the baseline does not know that state. So flai serve stops the story's dev server while the story is still active work.

The look does have the repository at hand. `ComponentsOptions.Story` defaults to `readStory`, which opens the project with `workitem.Open(root)` and drops it after `repo.Get`.

- Give `StoryLeft` the workflow to read the kinds by: `StoryLeft(wf workflow.Workflow, it *workitem.Item, since time.Time)`.
- Have the look get the story's project's workflow with the story. Either `Story` returns `repo.Flow()` beside the item, or `ComponentsOptions` gains a `Flow func(root string) workflow.Workflow` that defaults to the project's. Choose the one that changes `cmd/serve_actions.go` least: it builds `ComponentsOptions` without `Story`. Fall back to `baselineFlow` only where no project can be opened, as `Options.flow` does.
- Fix the test helper `stories` in `components_test.go` to match. Add cases to `TestStoryLeft` on a workflow that adds an active `testing` state and one that renames review. A story in `testing` has not left. A story that moved from the renamed review to done has left.
- Under § Leaving review in `design/system/component-library.md`, say that the states are read by kind from the story's project's workflow, with in progress and in review being the baseline's names for them.

It waits for nothing. It shares no path with the curate scheduler task, so the two run together.

## Done when

- `grep -n baselineFlow flai/internal/serve/components.go` finds only a fallback for a project that cannot be opened, if any.
- `TestStoryLeft` has a case for a story in a state of the active kind that the baseline lacks, which answers false, and it passes. The baseline cases are unchanged.
- `design/system/component-library.md` § Leaving review names the kinds the stop reads.
- `flai test` on the three paths passes.

## Notes
