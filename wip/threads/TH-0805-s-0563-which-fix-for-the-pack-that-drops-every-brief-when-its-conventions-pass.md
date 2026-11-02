---
id: TH-0805
title: "S-0563: which fix for the pack that drops every brief when its conventions pass the budget"
anchor:
  path: wip/kanban/stories/S-0563-the-orchestrator-s-conventions-sit-within-about-1-5kb-of-its-80kb-pack-budget.md
  item: S-0563
status: resolved
participants: [agent-S-0563, orchestrator, alex]
created: 2026-11-02T04:25:41Z
updated: 2026-11-02T04:54:25Z
---

# TH-0805 S-0563: which fix for the pack that drops every brief when its conventions pass the budget

On wip/kanban/stories/S-0563-the-orchestrator-s-conventions-sit-within-about-1-5kb-of-its-80kb-pack-budget.md.

## Entries

### 2026-11-02T04:25:41Z agent-S-0563
The cause of I-0250 is a cliff in `addDesign` (`flai/internal/context/context.go`). When the conventions plus the open-issues table pass the budget, the pack marks `exceeded: conventions`. It then loads no briefs and lists every document in the catalog. That contradicts TH-0032's "briefs are never cut", which holds when the named items or the briefs take the pack over.

Measured in this worktree:

| Pack | Before design | Result |
|------|---------------|--------|
| orchestrator, general, 80KB | conventions 72.8KB + issues table 7.5KB + header | 194KB with 138 briefs, about 1KB under the cliff |
| orchestrator, same at 72KB | over | 153KB, no briefs |
| S-0563's own pack, topics [all] | over | 169KB, no briefs |

The issues table grows about 250 bytes per open issue. That is why S-0518's issue files tip the orchestrator, and why the test result depends on what the repository holds.

Options:

1. **Remove the cliff (recommended).** When the conventions pass the budget, still load what is named and every brief, rank nothing, and keep `exceeded: conventions`. This applies TH-0032 to the conventions case too. A new ADR refines ADR-0049, and `agent-context.md` changes with it. Effect: the orchestrator's general pack always holds its briefs, and the test no longer depends on the repository's issue count. Cost: a pack already over budget gets bigger. S-0563's [all] pack would go from 169KB to roughly 250KB. Stories with narrower topics gain much less.
2. **Narrow fix, budget semantics unchanged.** Cap the open-issues table at a sixteenth of the budget (5KB at 80KB): the costliest rows, plus a line counting the rest. This frees about 2.5KB now and stops the issue count from moving the cliff. The cliff stays, so the next roughly 3.5KB of conventions tips it again.
3. **A role budget.** The orchestrator's role file sets its own budget, such as 96KB. This changes the role-file schema (ADR-0207) and the template. It buys margin, but the cliff stays.

I will build option 1 unless you choose otherwise. Option 2 can be added to either.

### 2026-11-02T04:27:03Z orchestrator (recommendation)
**Recommended: option 1.** Remove the cliff. When the conventions pass the budget, still load what is named and every brief, rank nothing, and keep `exceeded: conventions`. Record it in a new ADR that refines ADR-0049, and change `agent-context.md` with it. Leave options 2 and 3 out of S-0563.

**Why this is a recommendation, not an answer.** ADR-0049 § The pack, step 1, records the cliff on purpose: "When the conventions alone exceed the budget … the pack is the conventions and a catalog." Changing that rule is a new decision, so the call is yours.

**Why option 1.**
- On TH-0032 (D-0048) you chose C: keep the 80KB budget, keep every brief, and accept packs over budget until the conventions' topics are narrowed. You added the condition that the agent reads a document whole when its brief shows it matters. Option 1 applies that same answer to the conventions case. ADR-0049's step 1 was written before it.
- The cliff is not consistent with itself. About 1KB more conventions makes the orchestrator's general pack shrink by about 40KB, from 194KB to 153KB, and lose all 138 briefs. Nothing warns about it except a test.
- On TH-0804 (D-0462) you asked for the test to stop depending on the repository's wip and issues, or for a real margin. Option 1 does the first: the issue count can no longer tip the pack.

**Why not 2 or 3.** Each only moves the cliff. The next few KB of conventions or issues hits it again. Option 3 also changes the role-file schema (ADR-0207) and the template, which is beyond S-0563's scope.

**The cost.** A pack that is already over budget gets bigger. S-0563's own [all] pack goes from about 169KB to about 250KB, and packs with narrower topics grow much less. Since ADR-0104, packs are served in parts, so a large pack still reaches the agent.

Source: design/digest/D-0048-s-0138-s-pack-cannot-fit-80-kb-under-adr-0049-s-rules-which-way-how-many-kb-are.md

### 2026-11-02T04:27:14Z alex
option 1

### 2026-11-02T04:27:18Z alex
option 1

### 2026-11-02T04:54:25Z alex
Resolved.
