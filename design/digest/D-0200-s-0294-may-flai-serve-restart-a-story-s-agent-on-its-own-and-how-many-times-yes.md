---
id: D-0200
title: "S-0294: may flai serve restart a story's agent on its own, and how many times: yes"
kind: answer
sources: [TH-0208]
topics: [automation, flai, cli, go, code, all]
paths: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/restart.go, flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, design/adrs, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, docs/operators/settings.md, design/issues/I-0084-claude-code-ends-a-headless-agent-ten-minutes-after-its-turn-ends-even-while-its-background-sub-agent-is-still-working-and-flai-serve-leaves-the-story-in-progress-with-no-agent.md, design/issues/summary.md, design/issues/I-0063-flai-adr-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-adr-number.md, design/adrs/README.md, flai/internal/harness/harness.go]
date: 2026-10-06
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0200 S-0294: may flai serve restart a story's agent on its own, and how many times: yes

The operator, alex, answered TH-0208 at 2026-10-06T23:00:25Z.

What was asked:

```text
S-0294 has one part left: the second remediation from I-0084. When a story's agent ends with the story still in progress, flai serve marks the run failed and does nothing more. ADR-0043 gave restarting to a person, so whether flai serve may restart on its own is your call.
```

The answer:

```text
yes
```

Source: TH-0208.
