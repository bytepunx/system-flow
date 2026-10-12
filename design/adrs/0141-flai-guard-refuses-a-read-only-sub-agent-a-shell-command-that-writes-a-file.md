---
id: ADR-0141
title: flai guard refuses a read-only sub-agent a shell command that writes a file outside the temporary folder
status: proposed
date: 2026-10-12
supersedes: []
superseded_by: []
refines: [ADR-0059, ADR-0060]
topics: [cli, conventions, template]
---

# ADR-0141 flai guard refuses a read-only sub-agent a shell command that writes a file outside the temporary folder

## Context

[ADR-0059](0059-a-story-s-agent-hands-search-test-runs-and-verification-to-an-explorer-and-a.md) makes the explorer and the verifier read-only. Their definitions leave out `Edit` and `Write`, and [delegation.md](../conventions/delegation.md) says never to hand either an edit. The verifier keeps `Bash`, to run `flai verify S-nnnn --last`, `git diff`, and `go doc`.

[ADR-0060](0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md) has `flai guard` refuse a sub-agent's `git` and `flai` writes. It guards no file, and says the guard is not a shell. [ADR-0102](0102-while-auto-approve-is-off-flai-guard-refuses-a-story-s-sub-agent-a-write-to-a.md) added one file rule: a sub-agent's write under a `.claude/` folder.

[I-0052](../issues/I-0052-flai-guard-does-not-refuse-a-sub-agent-s-shell-commands-that-write-files-sed-i-perl-i-tee-redirection.md) records the gap. S-0189's verifier found that a sub-agent with `Bash` but no `Edit` or `Write` could still write a file through the shell: `sed -i`, `tee`, or a redirection. No sub-agent has been seen doing it. The read-only guarantee rests on the model following its prompt, which ADR-0060 rejected as a guarantee.

Task sub-agents are different. They edit files by design, and some do it through the shell, as [I-0058](../issues/I-0058-flai-guard-refuses-a-sub-agent-s-shell-command-whose-heredoc-text-reads-like-a-flai-write.md)'s instances show.

## Decision

`flai guard` refuses a read-only sub-agent a shell command that writes a file outside the system temporary folder.

- **Who.** A read-only sub-agent is one of two kinds:
  - one whose project definition, `.claude/agents/<agent_type>.md`, has a `tools` line naming neither `Edit` nor `Write`, as the template's explorer and verifier do
  - Claude Code's built-in `Explore` and `Plan`

  A task sub-agent, `general-purpose`, and any sub-agent with `Edit` or `Write` still write through the shell. A definition the guard cannot read, or one without a `tools` line, which grants every tool, is not read-only: the guard fails open, as ADR-0060 says. The calls of the story's agent and of the strategic agents are not changed.
- **What is a write.** The guard judges the simple commands it already splits, with heredoc bodies left out (S-0246). A command writes when it has one of these, and its target does not pass:
  - a redirection to a file, `>`, `>>`, `>|`, `&>`, `&>>`, or `N>`, attached or apart (`>f`, `2> f`)
  - `tee <file>`, `sed -i` or `--in-place`, `perl -i`
  - `cp`, `mv`, `rm`, `touch`, `mkdir`, `rmdir`, `truncate`, `ln`, `install`, `patch`, or `dd of=`
  - a formatter's write flag: `gofmt -w`, `goimports -w`, `go fmt`, `prettier --write` or `-w`, `eslint --fix`, or `golangci-lint` with `--fix`
- **What passes.** A target under the system temporary folder, `/tmp` or `$TMPDIR` as Go's `os.TempDir` gives it. `/dev/null` and other `/dev/` paths. A descriptor duplication, such as `2>&1` or `>&2`. A relative target is refused, since it resolves inside the checkout. The guard does not expand a variable, so `> "$TMPDIR/x"` is refused; the sub-agent writes to a `/tmp/...` path written out.
- **Refusal.** The guard exits 2. Its message names the command and the file, and tells the sub-agent to put the correction in its final message, for the story's agent to make.

## Consequences

- The explorer and the verifier cannot write a file through the common shell forms. A file they write any other way still shows in `git status` before the story can go to review.
- The guard stays not a shell. It does not see a write made by an interpreter or a script: `python -c`, a script under `scripts/`, `make`, or `go generate`.
- `git` and `flai` commands keep the rules they have under ADR-0060 and [ADR-0127](0127-flai-guard-lets-a-story-s-sub-agent-write-an-adr-with-flai-adr-new-topics-and.md). This rule does not judge them.
- A read-only sub-agent keeps scratch files in `/tmp`, as it may today.
- The guard reads the sub-agent's definition on each of its `Bash` calls. A project that edits a definition to add `Edit` or `Write` makes that sub-agent a writer, and the guard follows.
- Only Claude Code's sessions are covered, as with ADR-0060 and ADR-0102.

## Alternatives considered

- Refuse `Bash` to the verifier outright: it needs `flai verify --last`, `git diff`, and `go doc`.
- Refuse every sub-agent's shell writes: task sub-agents use `sed -i` and `gofmt -w` as part of their task, and would lose them.
- Run the sub-agent in a sandbox that cannot write the checkout: the harness offers no sandbox per sub-agent.
- Leave it to the definitions' prompts: the guarantee would rest on a model following its prompt, which ADR-0060 rejected.
