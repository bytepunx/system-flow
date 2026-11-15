---
id: ADR-0307
title: "flai reads the workflow file of the checkout it opens, refuses to open a repository whose file does not parse or has an error finding of the workflow rules, save for flai check, which reports them, re-reads it at each MCP call, and flai upgrade adds the baseline file only to a project that lacks it and keeps a project's file always"
status: accepted
date: 2026-11-15
supersedes: []
superseded_by: []
refines: [ADR-0301]
---

# ADR-0307 flai reads the workflow file of the checkout it opens, refuses to open a repository whose file does not parse or has an error finding of the workflow rules, save for flai check, which reports them, re-reads it at each MCP call, and flai upgrade adds the baseline file only to a project that lacks it and keeps a project's file always

## Context

ADR-0301 makes the board a file, `design/workflows/board.md`, which a project may hold in place of the baseline flai carries. §7 says the template ships the baseline, that `flai upgrade` adds it to a project that lacks it and never changes it afterwards, and that a file replaces the baseline whole. §9 leaves to S-0755 a project's own file, the template, upgrade, and the guides. Four things ADR-0301 does not settle:

- which checkout's file a command reads when it runs in a story's worktree, where `wip/` is the main checkout's;
- what a command does with a file that does not parse, or that parses into a workflow flai cannot run;
- how a process that keeps a repository open, the MCP server, sees a file that changes while it runs;
- how `flai upgrade` keeps a file the template ships, since today it replaces an unchanged shipped file and calls an edited one a conflict.

## Decision

flai reads the workflow file of the checkout it opens a repository in, and refuses to open one whose file it cannot run, naming the finding, everywhere but `flai check`, which reports it; the MCP server reads the file again when it changes; and `flai upgrade` adds the baseline file to a project that lacks it and never writes a project's file after.

### 1. The file read is the opened checkout's

`workitem.Open` reads `<layout.design>/workflows/board.md` under the repository's root, the checkout flai was started in, as `flai check` reads it today. In a story's worktree that is the story branch's file, so a story that changes the board runs and tests its own board before it is accepted; the main checkout, `flai serve`, and the dashboard read main's. A project with no file runs the baseline.

### 2. A file flai cannot run refuses the open

A file that does not parse, or whose workflow has a finding `flai check` reports as an error under a `workflow.` code (`workflow.file`, `workflow.rule`, `workflow.timer`, `workflow.state`, `workflow.dead-end`, `workflow.no-pull`, `workflow.no-accept`), makes `workitem.Open` fail with that finding: its code, its path and line, and its message. Nothing falls back to the baseline. Every command, the MCP server, and `flai serve` stop on it. A warning, `workflow.unreachable`, does not stop anything. `workflow.item`, an item naming a state the file lacks, is a finding on the item, not on the file, and stops nothing.

The rules that judge a workflow move from `flai check` into the workflow package, so that the load and the check judge it alike. `flai check` alone opens the repository without refusing: it runs the baseline for its other rules, reports the file's findings as today, and leaves out the per-item workflow judgments while the file is refused, since they would judge the items against a board the project does not run.

### 3. The MCP server reads it again when it changes

The MCP server keeps its repository open. At each tool call it reads the file again when its size or modification time changed since it was last read, and refuses the call with the finding while it cannot be run. `flai serve` opens the repository afresh at each look and logs a refused open once per distinct finding.

### 4. Upgrade keeps the project's file

`flai upgrade` adds `design/workflows/board.md` to a project that lacks it and never writes it once the project has one, whatever its content and whatever the lock records. When the template's file differs from the hash the lock recorded for it, upgrade says that the template's baseline changed and that the operator carries the change over in a story; the lock then records the new hash, so it says so once per change.

### 5. No advisory for a project without the file

`flai check` gives no advisory to a project with no workflow file. Such a project runs the baseline, which is a complete board, and `flai upgrade` adds the file.

## Consequences

- A story's agent can add a state on its branch and see the board, the moves, and the check follow it in its worktree before the operator accepts it.
- A broken file stops all work on the project but `flai check`, which says what to fix. That is the price of never running a board other than the one the file says.
- The MCP server pays one `stat` per call.
- A project's file drifts from the template's baseline until a story carries the change over; upgrade names each change once.

## Alternatives considered

- **Read main's file in a worktree, as `wip/` is main's.** A story that changes the board could not run it before acceptance, and S-0756's end-to-end story could not be tested on its branch.
- **Fall back to the baseline on a broken file.** The board would then run rules the operator removed, silently.
- **The baseline marker, so upgrade merges the file as a convention.** ADR-0301 §1 gives the file no marker, since flai reads nothing from its body, and a merge would replace the project's front matter, which is the whole workflow.
