---
id: D-0561
title: "S-0682: may I run one capped probe of the agent image (haiku, 0.02 USD), and with which Claude login: 1 - use ~/.claude-me"
kind: answer
sources: [TH-1000, S-0682]
topics: [flai, cli, go, code, all]
paths: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/pilot/request.go, flai/internal/pilot/request_test.go, flai/internal/pilot/run.go, flai/internal/pilot/run_test.go, scripts/agent-container.sh, docs/operators/runbooks/agent-container.md, design/system/flai-cli.md, design/tech/docker.md, design/issues/I-0224-scripts-agent-container-sh-can-run-only-a-story-s-own-agent-under-the-story-s-own-name-so-a-trial-run-of-the-image-acts-on-the-live-story.md, design/issues/summary.md]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0561 S-0682: may I run one capped probe of the agent image (haiku, 0.02 USD), and with which Claude login: 1 - use ~/.claude-me

The operator, alex, answered TH-1000 at 2026-11-10T08:53:31Z.

What was asked: S-0682 adds `scripts/agent-container.sh --probe S-nnnn`, which runs the agent image as `probe-S-nnnn` with a prompt that only reads and reports, and with no `FLAI_STORY`. Its agent asked for one capped paid trial of it, beside the free check (`--probe S-0682 -- env`), which needs no model:

```text
Recommendation: one run of `FLAI_AGENT_MODEL=haiku FLAI_AGENT_MAX_BUDGET_USD=0.02 scripts/agent-container.sh --probe S-0682`, then I check that S-0682's narrative and items are unchanged.

The script mounts `FLAI_CLAUDE_LOGIN`, default `~/.claude`. There is no `~/.claude/.credentials.json` on this host. Logins exist in `~/.claude-me`, `~/.claude-rmh`, and `~/.claude-virago`.

1. Yes, with `FLAI_CLAUDE_LOGIN=<folder>` (name it).
2. No: skip the paid trial; criterion 1 rests on the pilot test that reproduces I-0224 and the free check.
```

The answer:

```text
1 - use ~/.claude-me
```

So the operator allowed one paid probe of the agent image capped at 0.02 USD on haiku, with `FLAI_CLAUDE_LOGIN=~/.claude-me`. The login folder was already the answer of D-0494, which TH-1000 asked again; the spend was new and is the operator's call each time.

Sources: TH-1000, S-0682.
