---
id: D-0505
title: "sendBackReason puts a CI check name and reviewer login unquoted into the agent's prompt: file this as an issue"
kind: answer
sources: [TH-0500]
topics: [flai, go, code, orchestration]
paths: [flai/internal/serve/forgestatus.go]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: [S-0711, T-3118, TH-1079]
---

# D-0505 sendBackReason puts a CI check name and reviewer login unquoted into the agent's prompt: file this as an issue

The gap in `sendBackReason`, which puts a CI check's name and a reviewer's login unquoted into a sent-back story agent's prompt, is to be filed as an issue of its own, not folded into S-0395.

The orchestrator found it when it accepted S-0394 (flai 1.60.0) and recommended on TH-0500, in `flai/internal/serve/forgestatus.go`, quoting both in a code span and stripping control characters, as the thread body already does. Review comments are already safe: cut to 500 characters, quoted, and kept out of the prompt. The risk is low, since a check name comes from the repository's CI configuration.

The operator, alex, answered at 2026-11-07T09:27:11Z: "file this as an issue". The orchestrator replied that no permission of its own lets it file an issue, that no issue recorded the gap yet, and gave the command: `flai issue new "sendBackReason puts a CI check name and reviewer login unquoted into the agent's prompt" --autocommit`. A story's agent, the analyzer, or the operator files it.

Source: TH-0500.
