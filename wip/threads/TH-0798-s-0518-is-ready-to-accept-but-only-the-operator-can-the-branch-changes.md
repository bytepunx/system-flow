---
id: TH-0798
title: "S-0518 is ready to accept, but only the operator can: the branch changes protected .claude/settings.json files"
anchor:
  path: wip/kanban/stories/S-0518-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md
  item: S-0518
status: open
participants: [orchestrator]
created: 2026-11-02T01:07:14Z
updated: 2026-11-02T01:07:14Z
---

# TH-0798 S-0518 is ready to accept, but only the operator can: the branch changes protected .claude/settings.json files

On wip/kanban/stories/S-0518-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md.

## Entries

### 2026-11-02T01:07:14Z orchestrator
Recommendation: accept S-0518 yourself with `flai accept S-0518 --by alex`. The orchestrator cannot accept it or add it to the queue.

The only blocker comes from `flai accept S-0518 --by orchestrator --verified 444b35bf --dry-run`:

```text
only the operator (alex) accepts S-0518: its branch changes paths Claude Code protects: .claude/settings.json, template/root/.claude/settings.json
```

ADR-0106 says only the operator accepts a branch that changes those paths.

What was checked, at 444b35bfb995b33187e7cf7c88eac5e238aa0326 (the head of story/S-0518):

- `flai verify S-0518` passed every step. It ran gofmt, vet, golangci-lint, shellcheck, markdown, template, integration (3m17s), and smoke.
- Criterion 1 is met. Both `PreToolUse` hooks now refuse a call when the guard exits 126 or 127, with the message `flai guard did not run (exit N)`, and still refuse when it exits 2. `scripts/flai.sh` falls back to the installed flai when `go` is absent. `flai/cmd/guard_test.go` (`TestTheGuardHooksRefuseACallWhenTheGuardCannotStart`) reproduces I-0216 against both settings files. ADR-0235 records the decision, and the docs were updated with it.
- Criterion 2 is met. I-0216 is closed with a remediation line, and `design/issues/summary.md` no longer lists it.

Before you accept, review the two hook commands. The verifier found these edges:

- In a project made from the template, where `flai` is not on PATH, every `Bash` and `mcp__flai__*` call is refused. That includes an interactive session, because the first hook has no `FLAI_*` gate. ADR-0235 accepts this.
- These exits still let the call through, as designed: exit 1 (such as a failed `go build`), a signal kill (137 or 139), and a timeout.
- The in-image check inside `flai-agent-claude-code:local` rests only on the narrative. The verifier could not reproduce it.
