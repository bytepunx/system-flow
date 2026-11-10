---
id: D-0561
title: "S-0682: may I run one capped probe of the agent image (haiku, 0.02 USD), and with which Claude login: 1 - use ~/.claude-me"
kind: answer
sources: [TH-1000]
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

What was asked:

```text
S-0682 adds `scripts/agent-container.sh --probe S-nnnn`, which runs the agent image as `probe-S-nnnn` with a prompt that only reads and reports, and with no `FLAI_STORY`. T-2900 checks it on the rebuilt image.
```

The answer:

```text
1 - use ~/.claude-me
```

Source: TH-1000.
