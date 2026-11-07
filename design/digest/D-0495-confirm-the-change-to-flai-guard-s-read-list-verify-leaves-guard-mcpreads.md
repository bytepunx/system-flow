---
id: D-0495
title: "Confirm the change to flai guard's read list: verify leaves guard.MCPReads: Recommended: yes"
kind: answer
sources: [TH-0849]
topics: [guard, testing, flai, cli, go, template, conventions, code, all]
paths: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/internal/guard/orchestrate_test.go, flai/internal/roles/baseline/orchestrator.md, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, design/conventions/delegation.md, template/root/design/conventions/delegation.md, template/CHANGELOG.md, design/system/flai-cli.md, design/system/strategic-agents.md, design/issues/I-0240-flai-guard-still-lists-the-verify-mcp-tool-among-a-sub-agent-s-reads-though-it.md, design/issues/summary.md, design/issues/I-0277-a-story-s-agent-in-a-container-cannot-write-its-narrative-s-decisions-context.md, flai/internal/metrics/testdata/good/design/workflows/roles/orchestrator.md, docs/users/flai.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0495 Confirm the change to flai guard's read list: verify leaves guard.MCPReads: Recommended: yes

The operator, alex, answered TH-0849 at 2026-11-07T06:51:54Z.

What was asked:

```text
S-0556's notes ask for your confirmation before the close-out. The change is committed on `story/S-0556`, and nothing has been closed out yet.
```

The answer:

```text
**Recommended: yes.** Confirm that `verify` leaves `guard.MCPReads` and that sub-agents are refused `flai verify --job`. agent-S-0556 can then ask for the close-out and move S-0556 to review.

**Why it is yours:** S-0556's Notes ask you to confirm this before the close-out. The change also narrows what S-0546 chose, which D-0446 records. So I recommend and do not answer.

**Why yes:** I checked the diff on `story/S-0556` against D-0446.

- **D-0446 expects this change.** It says S-0556 narrows the guard and must keep two calls. Both stay:
  - The orchestrator verifying a story in review. `verify` is now in `MCPOrchestrates` in `flai/internal/guard/guard.go`, and the three copies of `orchestrator.md` name it. The orchestrator's plain `flai verify <S-nnnn> --json` still passes as a read.
  - A reader asking at a commit the record already holds. That goes through `flai verify --last` now, which runs nothing.
- **The gap is real.** The only change to `flagReads` is `--job` joining `--record-issues` among `verify`'s writing flags. So no sub-agent can start the committing job of ADR-0229 any more. Until now only the delegation convention asked them not to.
- **The one cost is small.** A reader that used the tool at a verified commit now uses `flai verify --last`.

On TH-0827 I approved the same plan, and I stand behind it as the answer.

Source: design/digest/D-0446-s-0546-kept-the-mcp-verify-tool-among-flai-guard-s-reads-since-a-call-at-a.md
```

Source: TH-0849.
