---
id: T-0926
type: task
nature: improvement
title: flai guard's design and docs say how it reads a heredoc, and I-0058 is closed
status: backlog
parent: S-0246
owner: alex
created: 2026-10-05T05:45:02Z
updated: 2026-10-05T05:45:02Z
transitions: []
stream: S-0246
tags: [flai, docs]
touches: [design/system/flai-cli.md, design/system/agent-context.md, docs/users/flai-reference.md, docs/users/flai.md, design/issues/I-0058-flai-guard-refuses-a-sub-agent-s-shell-command-whose-heredoc-text-reads-like-a-flai-write.md, design/issues/summary.md]
after: [T-0920]
---
# T-0926 flai guard's design and docs say how it reads a heredoc, and I-0058 is closed

## Work

Describe what T-0920 built wherever the guard's reading of a command line is described, then close the issue. This task waits for T-0920, because the words must match what the code does and the close reason names the fix.

- `design/system/flai-cli.md`, the `flai guard` row. After the sentence on how `internal/guard` splits a shell line, say three things:
  - A heredoc's body is taken out as text, for `<<` and `<<-` with a bare, quoted, or escaped delimiter.
  - A body fed to a shell is checked as a script.
  - A body whose delimiter is not quoted has its command substitutions checked.
- `docs/users/flai-reference.md`, the `flai guard` entry, and `docs/users/flai.md`, the paragraph on what a sub-agent may do. In each, add one plain sentence: a heredoc's text is not read as commands unless a shell runs it.
- `design/system/agent-context.md`, the paragraph on ADR-0060's refinement. Add one sentence: the guard took heredoc text for commands (I-0058, S-0246) and now reads it as a shell does, which refines ADR-0060 rather than reversing it.
- Run `flai issue close I-0058 --reason "..."` in the story's worktree. The reason names T-0920's fix and its tests. The command rewrites the issue and `design/issues/summary.md`.

## Done when

- The three documents say how the guard reads a heredoc, in words that match T-0920's code.
- I-0058 is closed, with a reason that names the fix, and `summary.md` no longer lists it as open.
- The markdown lint and `flai check --strict` pass.

## Notes
