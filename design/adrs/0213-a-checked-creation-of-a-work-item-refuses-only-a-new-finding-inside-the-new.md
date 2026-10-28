---
id: ADR-0213
title: "A checked creation of a work item refuses only a new finding inside the new item's story, as a checked edit does"
status: proposed
date: 2026-10-28
supersedes: []
superseded_by: []
refines: [ADR-0197]
topics: [cli]
---

# ADR-0213 A checked creation of a work item refuses only a new finding inside the new item's story, as a checked edit does

## Context

[ADR-0197](0197-a-checked-edit-of-a-work-item-refuses-only-a-new-finding-inside-its-story-s.md) scoped the checked edit of a work item to its story and left the other checked writers diffing two unscoped runs of `flai check` (its point 6), since no instance asked for more. [I-0191](../issues/I-0191-a-work-item-write-in-a-project-with-no-gitignore-makes-flai-cache-items-lock-and-then-refuses-itself-on-layout-gitignore.md) is that instance.

`itemnew.Create` is the checked creation: `flai story new`, `flai task new`, and `flai epic new` with `--body-stdin`, `--after`, or `--autocommit`, the MCP tool `item_new` with a body or `after`, `flai issue story`, and the dashboard's new form. In a project with no `.gitignore` and no `.flai-cache`, its first run of the check finds nothing. Writing the item takes the items lock of [ADR-0184](0184-every-work-item-write-takes-one-items-lock-and-merges-a-copy-read-before.md), which makes `.flai-cache/items.lock`, and the second run finds `layout.gitignore` on the repository's root. The diff calls it new, and the creation refuses itself. Run again, it stands, since the first run then knows the finding.

## Decision

A checked creation of a work item refuses only a new finding inside the new item's story, as a checked edit does.

1. **The second run is scoped to the new item's story.** `itemnew.Create` scopes the check it runs with the item in place to the new item's story, the story itself or a task's parent, with `check.ScopeToStories`, giving as the changed paths the item's file and its parent's, which it rewrites to list the item. `check.ScopeToItem` does this for both writers: `itemedit.Apply`'s scope moved there.
2. **What the scope marks outside is a warning.** A new finding outside the story, such as `layout.gitignore` on the repository's root, is returned among the creation's warnings, and the item stands.
3. **What the scope leaves out refuses nothing.** A new finding the scoped check drops, such as one on another open story's narrative, is neither a refusal nor a warning.
4. **What is inside still refuses.** A new finding inside the story's scope, such as a markdown finding on the new item or a `task.after` naming another story's task, refuses the creation and leaves nothing behind. A new `wip.overlap`, `touches.folder`, or `wip.after-hub` stays a warning, whichever story it names.
5. **An epic has no story.** A new epic refuses a new finding only on its own file or its parent's, and returns every other new finding as a warning.

The items lock stays where ADR-0184 puts it, and flai writes no `.gitignore` of its own; the main checkout's `flai check` still reports `layout.gitignore`.

## Consequences

- The first work item written in a project with no `.gitignore` is kept, with `layout.gitignore` among its warnings; I-0191's cause no longer occurs.
- A creation can no longer be refused for a finding it did not bring, from another writer or from the repository itself.
- A document save and `flai adr new` still diff two unscoped runs; no instance asks for more.

## Alternatives considered

- **Take the items lock outside the project, or only where `.flai-cache` is ignored.** It would move the lock ADR-0184 settled, and every other writer that makes `.flai-cache` would bring the same finding.
- **Write a `.gitignore` on the first write.** flai would change a tracked file the operator did not ask it to change, outside the item's commit.
- **Run the first check after the lock is taken.** It hides only this finding, not one another writer brings between the runs, which ADR-0197 already answers by scope.
