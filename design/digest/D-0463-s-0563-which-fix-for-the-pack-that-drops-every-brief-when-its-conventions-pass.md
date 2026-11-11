---
id: D-0463
title: "S-0563: which fix for the pack that drops every brief when its conventions pass the budget: option 1"
kind: answer
sources: [TH-0805, S-0563, I-0250]
topics: [context, prime, flai, cli, go, code, all]
paths: [flai/internal/context/context.go, flai/internal/context/brief_test.go, flai/cmd/prime.go, design/system/agent-context.md, design/adrs/0241-a-pack-whose-conventions-alone-exceed-its-budget-still-loads-what-is-named-and.md, design/issues/I-0250-the-orchestrator-s-conventions-sit-within-about-1-5kb-of-its-80kb-pack-budget.md]
date: 2026-11-02
status: absorbed
absorbed_by: "design/adrs/0241-a-pack-whose-conventions-alone-exceed-its-budget-still-loads-what-is-named-and.md#decision"
asked_again: [TH-1023]
cited: [TH-0806]
---

# D-0463 S-0563: which fix for the pack that drops every brief when its conventions pass the budget: option 1

The operator, alex, answered TH-0805 at 2026-11-02T04:27:18Z: "option 1". The orchestrator had recommended the same option a few seconds before; the operator answered rather than confirmed it.

What was asked: the cause of I-0250 is a cliff in `addDesign` (`flai/internal/context/context.go`). When the conventions plus the open-issues table (about 250 bytes per open issue) pass the budget, the pack marked `exceeded: conventions`, loaded no briefs, and listed every document in the catalog. The orchestrator's general pack sat about 1KB under it, so a branch's new issue records or a few lines of convention dropped all its briefs and failed `TestAThreadPackIsUnderAQuarterOfTheGeneralPack` (S-0422, S-0518). The options were: 1, remove the cliff; 2, cap the issue table at a sixteenth of the budget; 3, a budget per role in the role file.

The answer, option 1: a pack whose conventions alone exceed its budget still loads what is named and every brief, ranks nothing, and keeps `exceeded: conventions`. It applies TH-0032's "briefs are never cut" (D-0048) to the conventions case. Options 2 and 3 were left out, since each only moves the cliff.

S-0563 built it and recorded it in ADR-0241, which refines ADR-0049. Its narrative's Decisions add: "A sub-agent's pack (`AddBriefs`) and a reasons pack are unchanged." I-0250 was closed by it: at a 72KB budget the orchestrator's general pack is over on its conventions and still holds all 139 briefs, and the thread-pack test passes at 22.6%.

Source: TH-0805; S-0563 (wip/archive/agents/S-0563.md, Decisions); I-0250; ADR-0241.
