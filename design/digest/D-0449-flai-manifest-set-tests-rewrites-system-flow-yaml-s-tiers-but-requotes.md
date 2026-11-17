---
id: D-0449
title: "flai manifest set tests= rewrites system-flow.yaml's tiers but requotes unrelated paths lines, so a story that changes one tier hand-edits the manifest and validates it with flai check"
kind: decision
sources: [S-0547, S-0644, S-0805]
topics: [testing, cli, flai]
paths: [system-flow.yaml, flai/internal/manifest/settings.go, flai/internal/verify/select_test.go]
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: []
cited: [S-0805, S-0839, T-2560, TH-1172, TH-1298]
---

# D-0449 flai manifest set tests= rewrites system-flow.yaml's tiers but requotes unrelated paths lines, so a story that changes one tier hand-edits the manifest and validates it with flai check

S-0547 added `when: ["flai/**", "template/**"]` to this repository's `integration` and `smoke` tiers. Its narrative (`wip/archive/agents/S-0547.md`, § Decisions) records why it did not use flai's own command: "It is hand-edited in `system-flow.yaml`: `flai manifest set tests=` wrote the same tiers but requoted four unrelated `paths` lines, which would add conflicts with S-0548. `flai check` validates the hand edit."

So a story that changes one test tier while another open story also changes `system-flow.yaml` edits the tier by hand, keeps the other lines as they are, and runs `flai check` on the result, rather than rewriting the whole `tests` list with `flai manifest set`. The design does not record the requoting; `design/system/project-manifest.md` describes the keys only.

S-0644 did the same when it added the `design-routes` tier: "The manifest's `tests` block is edited by hand, as S-0592 did for `settings-index` (38e021d62); no flai command writes it." A repository-manifest test in `flai/internal/verify/select_test.go` then holds the new tier in the manifest, as S-0644's `TestRepositoryManifestChecksTheDashboardDesignsRoutesWhenEitherChanges` does.

S-0805 held to it again when it gave every tier of this repository a `suite:` (`wip/archive/agents/S-0805.md`, § Decisions): "system-flow.yaml's tiers are edited by hand, not with flai manifest set tests=, which requotes unrelated paths lines (D-0449)."

Sources: S-0547's, S-0644's, and S-0805's narratives, § Decisions.
