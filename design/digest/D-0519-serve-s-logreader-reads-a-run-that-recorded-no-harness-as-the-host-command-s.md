---
id: D-0519
title: serve's logReader reads a run that recorded no harness as the host command's, which S-0587 kept, since such a run predates S-0573's resolved-agent record and so was a command run
kind: decision
sources: [S-0587, S-0573]
topics: [flai, go, code]
paths: [flai/internal/serve/usage.go]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0519 serve's logReader reads a run that recorded no harness as the host command's, which S-0587 kept, since such a run predates S-0573's resolved-agent record and so was a command run

`logReader` in `flai/internal/serve/usage.go` picks the usage reader for a story's runs. It reads a run with no recorded harness as a run of the host's command. After S-0587, an agent that names no harness can fall back to `claude-code` ([ADR-0261](../adrs/0261-a-story-or-strategic-agent-that-names-no-harness-falls-back-after-the-project-s.md)). So the verifier asked whether `logReader` should change too.

S-0587 kept it: "`usage.go` `logReader` still reads a run with no recorded harness as the command's. Such a run predates S-0573's resolved-agent record, so it was a command run, and a story that never ran has no logs." (wip/archive/agents/S-0587.md, Decisions, Verifier review). Since S-0573, every run records the agent it resolved to (D-0478).

Source: S-0587's narrative, § Decisions.
