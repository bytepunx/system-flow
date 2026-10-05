---
id: TH-0141
title: "S-0246 plan: two tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0246-flai-guard-refuses-a-sub-agent-s-shell-command-whose-heredoc-text-reads-like-a-flai-write.md
  item: S-0246
status: resolved
participants: [planner-S-0246, alex]
created: 2026-10-05T05:45:23Z
updated: 2026-10-05T05:49:42Z
---

# TH-0141 S-0246 plan: two tasks in two layers

On wip/kanban/stories/S-0246-flai-guard-refuses-a-sub-agent-s-shell-command-whose-heredoc-text-reads-like-a-flai-write.md.

## Entries

### 2026-10-05T05:45:23Z planner-S-0246
The plan for S-0246 is two tasks in two layers.

- **Layer 1: T-0920.** Make `flai guard` read a heredoc's text as input, not as commands. The fix is in `commands` in `flai/internal/guard/guard.go`, with tests in `guard_test.go` and one hook-input test in `flai/cmd/guard_test.go`.
- **Layer 2: T-0926**, after T-0920. Update `flai-cli.md`, `agent-context.md`, `flai-reference.md`, and `flai.md` to say how the guard reads a heredoc. Then close I-0058 with `flai issue close`. It waits because its words and the close reason describe T-0920's code.

**The fix I propose**, since I-0058 recommends none. The guard is not a shell. It splits a line at newlines, so a heredoc's lines are read as commands. In S-0209's instance an apostrophe also opened a quote that never closed. The proposal:

- Recognise `<<WORD` and `<<-WORD`, with the delimiter bare, quoted, or escaped, and take each body out of the line as text, up to its closing delimiter.
- Keep checking a body where it really runs:
  - A body fed to a shell (`bash <<EOF`) is checked as a script, as `bash -c` already is.
  - In a body whose delimiter is not quoted, `$(...)` and backticks are checked.

**Assumptions:**

- No ADR. This narrows the guard's reading to what a shell does, a refinement of ADR-0060 like S-0175's, and `flai-cli.md` is where it is designed.
- A program other than a shell that reads a heredoc, such as `python3 -` that calls `subprocess` on flai, is not checked. That matches "careless, not hiding" in ADR-0060.
- The sub-agent's and the planner's checks share `commands`, so both get the fix. No template or settings file changes.
- The story's touches include I-0058's file and `design/issues/summary.md`, because `flai issue close` writes both. That overlaps S-0258's claim on `design/issues`, and `design/system/flai-cli.md` overlaps it already.

The cost of delay inputs are asked on TH-0138.

### 2026-10-05T05:49:42Z alex
Resolved.
