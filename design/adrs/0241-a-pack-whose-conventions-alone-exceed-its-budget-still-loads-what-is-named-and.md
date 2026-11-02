---
id: ADR-0241
title: "A pack whose conventions alone exceed its budget still loads what is named and every brief, ranks nothing, and keeps exceeded conventions"
status: accepted
date: 2026-11-02
supersedes: []
superseded_by: []
refines: [ADR-0049]
---

# ADR-0241 A pack whose conventions alone exceed its budget still loads what is named and every brief, ranks nothing, and keeps exceeded conventions

## Context

ADR-0049 § The pack, step 1, said that when the conventions alone exceed the budget, the pack is the conventions and a catalog of every document. When what the story names or the briefs take the pack over, everything named still loads and every brief still prints (TH-0032, D-0048). So the conventions case was the only one that cut briefs.

This made a cliff. "The conventions" here includes the open-issues table, which `Build` adds before the design step, at about 250 bytes per open issue. The orchestrator's general pack sat about 1KB under the budget: 72.8KB of conventions and a 7.5KB issue table. A few more lines of a convention, or a branch's new issue records, took it over. It then shrank from 194KB to 155KB and lost all 138 briefs, and nothing said so but `TestAThreadPackIsUnderAQuarterOfTheGeneralPack` (I-0250, recorded on S-0422 and S-0518). Story packs with wide topics hit the same cliff: S-0563's own pack was 169KB with no briefs.

## Decision

A pack whose conventions alone exceed its budget still loads what is named and every brief, ranks nothing, and keeps `exceeded: conventions`. The operator chose this on TH-0805, option 1. It applies TH-0032's "briefs are never cut" to the conventions case too. The digest step still loads within its share, as it does whenever a pack is over (TH-0624). The catalog lists what is neither loaded nor briefed, as in every other case.

This changes `addDesign` (`flai/internal/context/context.go`) only. A sub-agent's pack (`AddBriefs`, ADR-0059), which briefs while its half budget has room, and a pack primed for the orchestrator's reasons (ADR-0217), which has no briefs, are unchanged.

## Consequences

- No number of open issues and no few lines of convention can drop a pack's briefs. The size of the orchestrator's general pack no longer depends on what the repository holds, and `TestAThreadPackIsUnderAQuarterOfTheGeneralPack` passes whatever issues a branch adds.
- A pack already over budget on its conventions gets bigger. S-0563's pack, with topics `[all]`, goes from about 169KB to about 250KB. Packs with narrower topics grow less. Since ADR-0104, packs are served in parts, so a large pack still reaches the agent.
- The budget now limits only the ranked step, in every case. Narrowing the conventions' topics is still what brings a pack back under it.

## Alternatives considered

- Cap the open-issues table at a sixteenth of the budget (TH-0805, option 2). This removes the dependence on the issue count but keeps the cliff: the next few KB of conventions tips it again.
- A budget per role in the role file (TH-0805, option 3). This buys margin but keeps the cliff, and it changes the role-file schema (ADR-0207) and the template.
