---
id: D-0466
title: "S-0508: the live container run (criterion 6, T-2169's first half) cannot run before S-0508 is accepted and published; move it to a follow-up story: Recommended: say yes, as you answered for S-0423 on…"
kind: answer
sources: [TH-0812]
topics: [harness, docker, flai, cli, go, code, all]
paths: [flai/internal/harness/container.go, flai/internal/harness/container_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/protected/protected.go, flai/internal/protected/protected_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/start.go, flai/internal/serve/start_test.go, flai/internal/serve/stop.go, flai/internal/serve/stop_test.go, flai/internal/serve/restart.go, docs/operators/settings.md, docs/users/flai-reference.md, docs/users/flai.md, docs/operators/index.md, docs/operators/runbooks/install.md, design/system/flai-cli.md, design/system/agent-adapters.md, design/tech/docker.md, design/issues/summary.md, design/issues, flai/internal/guard/call.go, flai/internal/guard/call_test.go, flai/cmd/serve_settings.go, flai/internal/serve/claudecheck.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go, docs/operators/runbooks/agent-container.md]
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: []
cited: [S-0509, S-0510, S-0565, TH-0814, TH-0815, TH-0820, TH-0911]
---

# D-0466 S-0508: the live container run (criterion 6, T-2169's first half) cannot run before S-0508 is accepted and published; move it to a follow-up story: Recommended: say yes, as you answered for S-0423 on…

The operator, alex, answered TH-0812 at 2026-11-02T18:08:27Z.

What was asked:

```text
**Recommended:** move criterion 6 and the live run in T-2169 to a draft follow-up story in E-0029, with `after: S-0508`.
```

The answer:

```text
**Recommended:** say yes, as you answered for S-0423 on TH-0803 (D-0461).

- S-0508's agent removes criterion 6 from S-0508 and narrows T-2169 to the documentation.
- It writes a draft story in E-0029 with `after: S-0508`. The draft holds criterion 6's words and the live-run half of T-2169's Work. Its Notes say why it waits: the adapter must reach the installed flai first.
- S-0508's Notes record the move.

**Why:** this is the same case as D-0461 and D-0422. The `container` adapter, its host settings, and the serve start path exist only on `story/S-0508`. The `flai serve` on this host is the installed release, so no story can run on `harness: container` until S-0508 is accepted, published, and the host upgraded.

**Why I recommend rather than answer:** this takes a criterion out of a finalized story and adds a story. That changes scope, so it is your call. Finalizing the draft is yours too.

**Alternative:** keep criterion 6. I don't recommend it: S-0508 could never meet the definition of done.

Source: design/digest/D-0461-s-0423-the-live-publish-run-t-1783-criterion-4-cannot-finish-before-s-0423-is.md
```

Source: TH-0812.
