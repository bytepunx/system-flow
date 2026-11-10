---
id: ADR-0264
title: "A checked write of a work item runs only the rules the write can affect, before and after, in place of the whole flai check"
status: accepted
date: 2026-11-10
supersedes: []
superseded_by: []
refines: [ADR-0197, ADR-0213]
topics: [cli]
---

# ADR-0264 A checked write of a work item runs only the rules the write can affect, before and after, in place of the whole flai check

## Context

A checked write of a work item is a checked creation, `itemnew.Create`, or a checked edit, `itemedit.Apply`. The creations are `flai story new`, `flai task new`, `flai epic new` with a body, `after`, or `--autocommit`, the MCP tool `item_new`, and `flai issue story`. The edits are `flai edit`, finalize, `flai criteria tick`, and flai serve's replan. A checked write ran the whole of `flai check` twice, before and after the write. It then refused what [ADR-0197](0197-a-checked-edit-of-a-work-item-refuses-only-a-new-finding-inside-its-story-s.md) and [ADR-0213](0213-a-checked-creation-of-a-work-item-refuses-only-a-new-finding-inside-the-new.md) say it refuses: a new finding inside the written item's story.

At 3,233 items and 5,068 markdown files under `wip`, each whole run listed every item, linted every file under `wip`, started `node` for the component examples, and scanned the digest. The two runs were about 97% of a write. A write took 7 to 18 s wall, and the serve log shows `item.new` at p50 7.9 s and up to 32 s. That is cause 8 in [server-performance.md](../system/server-performance.md#where-the-time-goes).

Most of those rules read no work item. A finding they bring cannot be on the written item's story. They can still report a finding another writer brought between the two runs, but ADR-0197 and ADR-0213 already refuse nothing for such a finding.

## Decision

A checked write of a work item runs only the rules the write can affect, before and after the write, in place of the whole `flai check`.

1. **The rule set.** `check.RunItemWrite` takes the written items and the written files. It runs:
   - the item rules (`oneItem`, `history`) on the written items, their parents, and their children, since a child's rules report on its parent's file; the duplicate-ID pass still runs over every item;
   - the rules over the item graph, on every item: `wip.overlap`, `story.component-tag`, `story.after` and `task.after`, `touches.folder`, `wip.after-hub`, `wip.after-reason`, `story.unaccepted`, and `epic.lags-stories`, and the board;
   - the narrative rules on the written stories' narratives only, and the checks of `wip/agents/index.md` only when the write wrote it;
   - the markdown lint of the written files under `wip`, and the conflict markers of the written files.
2. **What it skips.** It skips every rule that reads no work item: the layout (`layout.gitignore` among them), the instructions files, the activity documents, the documentation, the ADR index, the conventions, roles, formulas, issues, experiments, analysis, the digest, the design tokens, the component examples (so it starts no `node`), threads, and conversations.
3. **The before run covers what the after run covers.**
   - A creation runs the rule set before on the parent and its file, and after on the new item, its parent, and their files.
   - An edit runs both on one set: the item, its parent, and on a reparent the former and new parent. The set's files are the item's old and new file, the parents' files, a retitled story's narrative, `wip/agents/index.md`, the files a retitle relinks, and the story's own item, task, and narrative files. The after run adds every file the edit wrote. A retitle finds the files it relinks before it writes anything.
4. **What a write refuses is unchanged.** The after run is scoped with `check.ScopeToItem` and compared with the before run as before. A new finding the write brings inside its story refuses it. A new `wip.overlap`, `touches.folder`, `wip.after-hub`, or `wip.after-reason` is a warning. An error on the edited item's own file refuses an edit. The undo puts every file back.
5. **The whole check stays.** `flai check`, `flai task done`, `flai verify`, the close-out, the merge queue, and CI still run every rule. `check.Run` is unchanged.

This refines ADR-0197 and ADR-0213, whose two runs it narrows to the rules a write can affect.

## Consequences

- A checked creation or edit costs the item list, the graph rules, and the lint of a few files, not two whole checks. [server-performance.md § A checked write runs only the rules it can affect (S-0598)](../system/server-performance.md#a-checked-write-runs-only-the-rules-it-can-affect-s-0598) records the figures.
- A finding of a skipped rule is never a warning on a write now. The first item written in a project with no `.gitignore` is still kept, as ADR-0213 decided. But `layout.gitignore` is no longer among its warnings, as ADR-0213's Consequences said it would be. The main checkout's `flai check` still reports it.
- A finding another writer brings between the two runs is no longer reported when it is on a file the write neither writes nor reads. An example is a markdown finding on another story's task, or, for an epic's edit, on a story's task. ADR-0197 point 3 returned it as a warning; now it is neither a warning nor a refusal, and the whole check reports it.
- An edit still reads its story's own files. So a finding another writer brings there still refuses the edit, as ADR-0197 has it.
- A creation reads only the new item and its parent. So a finding another writer brings on a sibling task between the two runs no longer refuses the creation.
- A file that gains a link to a retitled item after the retitle found its links, and before it writes, is not relinked.
- A rule added to `flai check` later is skipped on a write unless it is added to the rule set. One that reads work items should be weighed for it.

## Alternatives considered

- **Keep the whole check and make it faster.** S-0600 and S-0601 make a whole check cheaper, through kept parses, kept lint results, and rules run at once. But the floor of a whole run is still every file under `wip` and `node`, for rules that a write cannot bring a finding on.
- **Run the check once, after the write, and refuse every finding inside the story.** A finding already in the story would refuse an unrelated edit. The before run is what makes a refusal mean "this write brought it".
- **Scope only the markdown lint.** The lint is the largest single cost, but `node`, the digest, the documentation, and the threads would still run twice a write.
