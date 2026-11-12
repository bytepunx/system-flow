---
id: D-0505
title: "sendBackReason puts a CI check name and reviewer login unquoted into the agent's prompt: file this as an issue"
kind: answer
sources: [TH-0500, I-0340, S-0711]
topics: [flai, go, code, orchestration]
paths: [flai/internal/serve/forgestatus.go, design/system/flai-cli.md]
date: 2026-11-07
status: absorbed
absorbed_by: "design/system/flai-cli.md#reading-the-pull-request-back"
asked_again: []
cited: [S-0711, T-3118, TH-1079]
---

# D-0505 sendBackReason puts a CI check name and reviewer login unquoted into the agent's prompt: file this as an issue

The gap in `sendBackReason`, which puts a CI check's name and a reviewer's login unquoted into a sent-back story agent's prompt, is to be filed as an issue of its own, not folded into S-0395.

The orchestrator found it when it accepted S-0394 (flai 1.60.0) and recommended on TH-0500, in `flai/internal/serve/forgestatus.go`, quoting both in a code span and stripping control characters, as the thread body already does. Review comments are already safe: cut to 500 characters, quoted, and kept out of the prompt. The risk is low, since a check name comes from the repository's CI configuration.

The operator, alex, answered at 2026-11-07T09:27:11Z: "file this as an issue". The orchestrator replied that no permission of its own lets it file an issue, that no issue recorded the gap yet, and gave the command: `flai issue new "sendBackReason puts a CI check name and reviewer login unquoted into the agent's prompt" --autocommit`. A story's agent, the analyzer, or the operator files it.

It was filed as I-0340, and S-0711 fixed it, accepted 2026-11-12. From S-0711's Decisions: "sendBackReason quotes each check name and reviewer login through forgeName: white-space control characters become one space, every other control or format (Cf) character is dropped, runs of spaces are folded, and the result goes in a code span through codeSpan." Format characters are dropped because a bidirectional override can hide text from a human reader. Only the reason handed to the prompt changed: the thread's title and body keep the names as they are, since the body already puts check names in code spans and a GitHub login holds only letters, digits, and hyphens. `design/system/flai-cli.md` § Reading the pull request back now states the quoting.

Sources: TH-0500; I-0340; S-0711 (`wip/archive/agents/S-0711.md`, Decisions).
