---
id: D-0449
title: "flai manifest set tests= rewrites system-flow.yaml's tiers but requotes unrelated paths lines, so a story that changes one tier hand-edits the manifest and validates it with flai check"
kind: decision
sources: [S-0547]
topics: [testing, cli, flai]
paths: [system-flow.yaml, flai/internal/manifest/settings.go]
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: []
cited: [T-2560]
---

# D-0449 flai manifest set tests= rewrites system-flow.yaml's tiers but requotes unrelated paths lines, so a story that changes one tier hand-edits the manifest and validates it with flai check

S-0547 added `when: ["flai/**", "template/**"]` to this repository's `integration` and `smoke` tiers. Its narrative (`wip/archive/agents/S-0547.md`, § Decisions) records why it did not use flai's own command: "It is hand-edited in `system-flow.yaml`: `flai manifest set tests=` wrote the same tiers but requoted four unrelated `paths` lines, which would add conflicts with S-0548. `flai check` validates the hand edit."

So a story that changes one test tier while another open story also changes `system-flow.yaml` edits the tier by hand, keeps the other lines as they are, and runs `flai check` on the result, rather than rewriting the whole `tests` list with `flai manifest set`. The design does not record the requoting; `design/system/project-manifest.md` describes the keys only.

Source: S-0547's narrative, § Decisions.
