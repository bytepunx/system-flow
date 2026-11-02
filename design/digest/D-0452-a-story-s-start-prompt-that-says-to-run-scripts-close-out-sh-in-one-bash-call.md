---
id: D-0452
title: A story's start prompt that says to run scripts/close-out.sh in one Bash call comes from a flai older than S-0546, and S-0549 followed it over work-management.md's close-out job
kind: decision
sources: [S-0549, S-0546, S-0420]
topics: [workflow, testing, flai]
paths: [design/conventions/work-management.md, flai/internal/harness/harness.go, scripts/close-out.sh]
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0452 A story's start prompt that says to run scripts/close-out.sh in one Bash call comes from a flai older than S-0546, and S-0549 followed it over work-management.md's close-out job

S-0549's narrative records: "Operator instruction over convention: the start prompt says to run `scripts/close-out.sh` myself in one Bash call, where `work-management.md` says to ask flai serve for the close-out job. I follow the prompt." (wip/archive/agents/S-0549.md, Decisions).

The two disagree because flai serve runs the installed flai, not this tree's. Since S-0546 (ADR-0229), the start prompt in `flai/internal/harness/harness.go` tells a story's agent to ask flai serve for the close-out with the MCP tool `verify`, to hold `wait_for_events` until `verified`, and to run `scripts/close-out.sh` itself only when `verify` says flai serve is not running. `design/conventions/work-management.md` says the same. A start prompt that tells the agent to run the script itself comes from an installed flai older than that change.

S-0549 treated the start prompt as the operator's instruction, which outranks a convention under CLAUDE.md's precedence, and recorded the conflict in its narrative's Decisions. A later story that meets the same conflict can do the same. Once the installed flai is newer than S-0546, the two agree.

S-0420 met the same conflict and did the same: "The close-out is run here with `scripts/close-out.sh`, as the operator's start prompt says, not as a flai serve job (work-management.md): the operator's instruction takes precedence." (wip/archive/agents/S-0420.md, Decisions).

Sources: S-0549, S-0546, S-0420.
