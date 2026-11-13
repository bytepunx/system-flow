---
id: ADR-0298
title: "flai guard refuses a Write that replaces a work item file whole with content that is not that item, and points the caller to Edit or flai edit"
status: accepted
date: 2026-11-13
supersedes: []
superseded_by: []
refines: [ADR-0060, ADR-0102]
---

# ADR-0298 flai guard refuses a Write that replaces a work item file whole with content that is not that item, and points the caller to Edit or flai edit

## Context

[I-0352](../issues/I-0352-a-story-s-agent-can-overwrite-a-work-item-file-whole-with-write-front-matter.md) records what happened on 2026-11-10. S-0665's agent overwrote the main checkout's T-3079 file with a Claude Code `Write` of the content `PLACEHOLDER`, six seconds after `flai task new` had made it. Twenty seconds later it wrote the file back. In between, every `Repo.List` failed, because the file had no front matter (I-0344, I-0345).

[tooling.md](../conventions/tooling.md) lets an agent hand-edit an item's body. Nothing told a `Write` that replaces the whole file, front matter included, from an `Edit` that changes part of it. [ADR-0060](0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md) has `flai guard` refuse a sub-agent's writes to work items through `flai` and `git`. [ADR-0102](0102-while-auto-approve-is-off-flai-guard-refuses-a-story-s-sub-agent-a-write-to-a.md) runs the guard before `Edit`, `MultiEdit`, `Write`, and `NotebookEdit` in every session `flai serve` starts. Until now the guard dropped what a `Write` would write, so no rule could judge it.

## Decision

`flai guard` refuses a `Write` that replaces a work item file whole with content that is not that item, and points the caller to `Edit` or `flai edit`.

- **Which calls.** A call that replaces its file whole: Claude Code's `Write` and OpenCode's `write`. The guard now reads the content each gives. `Edit`, `MultiEdit`, `NotebookEdit`, and an OpenCode patch change part of a file and are not judged.
- **Which files.** An epic, story, or task file, `<ID>-<slug>.md` in a folder `epics`, `stories`, or `tasks` under a folder `kanban`, the archive's included. The ID's prefix must match the folder. The guard judges the file by its path, in the main checkout or a story's worktree, and reads nothing from disk.
- **Who.** Every caller the hook runs for: a story's agent, its sub-agents, and a strategic agent. The rule comes before the sub-agent and role rules.
- **What passes.** Content that `workitem.ParseItem` reads as an item whose `id` and `type` are the ones the file's name and folder give. So an agent may still rewrite an item's body whole, so long as it keeps the item's front matter.
- **Refusal.** The guard refuses with the file, the item, and what is wrong: no front matter, front matter that does not parse, or another `id` or `type`. It says to change the body with `Edit` and the fields with `flai edit`.

## Consequences

- A `Write` like I-0352's is refused before the file changes, so `Repo.List` never meets the file without front matter.
- A whole-file rewrite that keeps valid front matter still passes. The rule does not judge fields a command owns, which `flai check` reports.
- A `Write` that makes a new item file with matching front matter passes. Making items by hand stays a convention, in `tooling.md`, rather than a guard rule.
- A write through the shell, such as a redirection or `sed -i`, is not judged. [ADR-0141](0141-flai-guard-refuses-a-read-only-sub-agent-a-shell-command-that-writes-a-file.md) covers a read-only sub-agent's shell writes. The guard is not a shell.
- A `Write`'s content now passes through the guard, which parses at most one item a call.

## Alternatives considered

- Refuse every `Write` to an item file. That needs no content, but it would also refuse a rewrite of a body that keeps its front matter. The plan's thread, TH-1173, kept the narrower rule.
- Compare the content with the file on disk. That needs the guard to read the project, which it does not do for file writes. The file's name and folder already give the item's ID and type.
- Warn instead of refusing. The 20 seconds of failed lists happened in the call itself, before any warning could be read.
