---
id: TH-0877
title: "S-0407 is ready for your acceptance: it changes .opencode/ plugin paths Claude Code protects, so I cannot accept it"
anchor:
  path: wip/kanban/stories/S-0407-flai-guard-judges-an-opencode-agent-s-calls-through-a-plugin-on-tool-execute-before-and-the-opencode-adapter-states-the-guard-only-where-the-plugin-can-tell-a-sub-agent-s-call.md
  item: S-0407
status: resolved
participants: [orchestrator, alex]
created: 2026-11-07T21:32:36Z
updated: 2026-11-07T21:45:57Z
---

# TH-0877 S-0407 is ready for your acceptance: it changes .opencode/ plugin paths Claude Code protects, so I cannot accept it

On wip/kanban/stories/S-0407-flai-guard-judges-an-opencode-agent-s-calls-through-a-plugin-on-tool-execute-before-and-the-opencode-adapter-states-the-guard-only-where-the-plugin-can-tell-a-sub-agent-s-call.md.

## Entries

### 2026-11-07T21:32:36Z orchestrator
**Recommended:** accept S-0407 yourself:

```bash
flai accept S-0407 --by alex
```

**Blocker:** the dry run of my acceptance refuses it (ADR-0106):

```text
only the operator (alex) accepts S-0407: its branch changes paths Claude Code protects: .opencode/plugins/flai-guard.js, template/root/.opencode/plugins/flai-guard.js
```

No other blocker was listed.

**Verified at 1e097ea2963:** `flai verify S-0407 --last` passed every step, integration and smoke included. The verifier matched every criterion to changed files:

| Criterion | Changed files |
|-----------|---------------|
| 1 | `flai/internal/guard/opencode.go`, `opencode_test.go`, `guard.go`, `claudecode_test.go`, `flai/cmd/guard.go`, `guard_test.go`, `docs/users/flai-reference.md` |
| 2 | `.opencode/plugins/flai-guard.js`, `template/root/.opencode/plugins/flai-guard.js` |
| 3 | `flai/internal/harness/opencode.go` (`Guard: true`), `opencode_test.go`, `harness_test.go`, `flai/internal/serve/agents_test.go`, `design/system/agent-adapters.md` |
| 4 | `flai/tests/integration/opencode_test.go` |
| 5 | `docs/operators/index.md`, `docs/users/flai.md`, `design/system/flai-cli.md`, `template/CHANGELOG.md` |

**Worth knowing before you accept:**

- Criterion 4's test needs a live opencode 1.19.5 and skips without it. The story's agent installed it in `.flai-cache/scratch/opencode`, git-ignored and not global, to run it once. Your answers on TH-0533 (D-0507) and D-0506 said fixtures only, with no OpenCode on the host. TH-0866 still awaits your answer on fixture-based criteria. A second test, `TestFlaiGuardPluginUnderNode`, loads the plugin under `node` alone.
- In a session flai started (`FLAI_STARTED_BY`, `FLAI_STORY`, or `FLAI_ROLE` set), the plugin refuses every call when `flai` is missing, rather than passing it (ADR-0235). That is stricter than criterion 2 says, and only outside flai serve does it pass calls.
- S-0410, in progress, touches the same harness files and `docs/operators/index.md`. Expect it to sync after this lands.

S-0408 and S-0409 are held in ready until S-0407 is done.

### 2026-11-07T21:45:57Z alex
Resolved: S-0407 was accepted
