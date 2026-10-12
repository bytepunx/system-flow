---
id: T-1460
type: task
nature: improvement
title: An ADR refining ADR-0060 proposes which sub-agents flai guard refuses a shell write and which commands count as one, and agent-adapters.md records it
status: in-progress
parent: S-0368
owner: alex
created: 2026-10-12T16:57:18Z
updated: 2026-10-12T17:19:23Z
transitions:
  - to: ready
    at: 2026-10-12T17:19:22Z
    by: agent-S-0368
  - to: in-progress
    at: 2026-10-12T17:19:23Z
    by: agent-S-0368
stream: S-0368
tags: [guard, adr, design]
touches: [design/adrs/README.md, design/system/agent-adapters.md]
---
# T-1460 An ADR refining ADR-0060 proposes which sub-agents flai guard refuses a shell write and which commands count as one, and agent-adapters.md records it

## Work

The goal asks for a fix proposed from I-0052's instances before it is built. The one instance is S-0189's verifier: it has `Bash` and no `Edit` or `Write`, yet nothing stops it writing a file through the shell. ADR-0060 says files are not guarded.

Propose the fix as an ADR with `flai adr new`, refining ADR-0060 and ADR-0059. The starting proposal follows; change what you would decide differently and say why in the narrative's `## Decisions`.

- **Who is refused.** A sub-agent whose definition gives it no `Edit` and no `Write`: the template's `explorer` and `verifier`, read from the `tools` line of `.claude/agents/<agent_type>.md`, and the built-in `Explore` and `Plan`. A task sub-agent and `general-purpose` keep writing through the shell, as I-0058's instances show they do. A definition the guard cannot read is not read-only, so the guard fails open as ADR-0060 says.
- **What counts as a write.** A command that writes a file inside the project's checkout or the story's worktree:
  - redirection `>`, `>>`, or `&>` to a file
  - `tee <file>`, `sed -i`, `perl -i`
  - `cp`, `mv`, `rm`, `touch`, `mkdir`, `truncate`, `ln`, `install`, `patch`, `dd of=`
  - a formatter's write flag: `gofmt -w`, `go fmt`, `prettier --write`, `eslint --fix`, `golangci-lint --fix`.
- **What passes.** A target under the system temporary folder, `/dev/null`, and a descriptor duplication such as `2>&1`.
- **Limits.** The guard stays not a shell. A script or interpreter that writes, such as `python -c` or a script under `scripts/`, is not seen. A file written that way still shows in `git status` before review.

Record the decision in `design/system/agent-adapters.md` under `### Sub-agents and strategic agents` and `### Refusing a sub-agent's write`.

Open a thread on S-0368 with the proposal, recommendation first. Go on with the next layer meanwhile, and change the ADR if the answer differs.

## Done when

- [ ] An ADR refining ADR-0060 names the sub-agents refused, the commands counted as writes, what passes, and the limits.
- [ ] `design/adrs/README.md` lists it.
- [ ] `design/system/agent-adapters.md` states the rule in the two sections named above.
- [ ] A thread on S-0368 carries the proposal.

## Notes

The new ADR's file is not named in `touches`: `flai adr new` allocates its number, and `flai task done` widens the touches with it when it commits.
