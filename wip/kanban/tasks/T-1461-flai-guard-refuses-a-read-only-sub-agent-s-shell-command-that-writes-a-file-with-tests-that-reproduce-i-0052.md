---
id: T-1461
type: task
nature: remediation
title: flai guard refuses a read-only sub-agent's shell command that writes a file, with tests that reproduce I-0052
status: backlog
parent: S-0368
owner: alex
created: 2026-10-12T16:57:30Z
updated: 2026-10-12T16:57:30Z
transitions: []
stream: S-0368
tags: [guard, flai]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/internal/guard/writes.go, flai/internal/guard/writes_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, docs/users/flai-reference.md]
after: [T-1460]
---
# T-1461 flai guard refuses a read-only sub-agent's shell command that writes a file, with tests that reproduce I-0052

## Work

Build the rule T-1460's ADR decides. It waits for T-1460 because the ADR fixes who is refused and which commands count as writes.

- In `flai/internal/guard/writes.go`, find a file write among a command's words:
  - redirection to a file, in its attached forms too (`>f`, `2>f`, `&>f`)
  - the commands and flags the ADR lists.
  
  Return the target, so that `/dev/null`, the system temporary folder, and a descriptor duplication such as `2>&1` pass. Work on the words `commands` in `guard.go` already splits, heredoc bodies left out as S-0246 made them.
- Tell a read-only sub-agent by its `agent_type`: the `tools` line of the project's `.claude/agents/<agent_type>.md` with neither `Edit` nor `Write`, or a built-in the ADR names. Fail open when the definition cannot be read.
- In `Decide`, refuse such a sub-agent's `Bash` call that writes. Exit 2 with the command and the file. Say what to do instead: put the correction in the final message for the story's agent to make.
- Reproduce I-0052 in `flai/internal/guard/guard_test.go` or `writes_test.go`:
  - A verifier is refused `sed -i`, `perl -i`, `tee`, `echo x > f`, and `cat > f <<EOF`.
  - A verifier is allowed `go test ./... > /tmp/out`, `2>&1`, and `> /dev/null`.
  - A `general-purpose` or task sub-agent is still allowed `sed -i`.
- Add one hook-JSON case to `flai/cmd/guard_test.go`.
- Change the guard's help in `flai/cmd/guard.go` to say what it now refuses. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.

## Done when

- [ ] A read-only sub-agent's shell write inside the project is refused with exit 2 and a reason naming the file.
- [ ] Writes to `/dev/null` and the temporary folder, descriptor duplications, and every write by a sub-agent with `Edit` or `Write` pass.
- [ ] Tests reproduce each command I-0052 names and pass.
- [ ] `docs/users/flai-reference.md` matches the help.
- [ ] `flai test` on the changed paths passes.

## Notes

S-0353, in the backlog, moves `Decide` onto a neutral call (ADR-0130). This task adds to `Decide` as it stands; S-0353 carries the rule over when it is worked.
