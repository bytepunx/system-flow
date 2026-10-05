---
id: T-0986
type: task
nature: research
title: Measure I-0068's packs and propose the fix as an ADR refining ADR-0049
status: backlog
parent: S-0261
owner: alex
created: 2026-10-05T05:52:07Z
updated: 2026-10-05T05:52:07Z
transitions: []
stream: S-0261
tags: [flai, docs]
touches: [design/adrs]
---
# T-0986 Measure I-0068's packs and propose the fix as an ADR refining ADR-0049

## Work

The story asks for a fix proposed from I-0068's instances before it is built. Measure first, then decide.

For each pack, record two sizes and set them against Claude Code's MCP tool result limit, 25,000 tokens unless `MAX_MCP_OUTPUT_TOKENS` raises it:

- the budgeted text size and its `exceeded`;
- the size of the JSON the MCP `prime` tool returns, in characters and in estimated tokens.

Measure these packs, priming each with `flai prime --json`, which gives an archived story the pack it would get today:

- the stories I-0068 names: S-0205, S-0210, and S-0211;
- the planner's pack for S-0261 (`--role plan --story S-0261`). It was 81,903 bytes against a budget of 81,920, but 92,948 characters as JSON, and Claude Code refused it too.

Keep the two causes the numbers show apart:

- **Briefs past the budget.** ADR-0049 and TH-0032 print every brief, so `exceeded: briefs` packs run 126 to 152 KB.
- **The encoding.** The JSON carries what the printed pack does not: each convention's `kept` and `left_out` pieces, the catalog's outlines, the topic sources, and the escaping.

Weigh at least these fixes, each by what it removes from the measured packs:

- fit the budget to the encoded result;
- have MCP `prime` return the printed text instead of the JSON;
- page the pack, with a part parameter that keeps each part under the limit, as the inbox's cap did for I-0023;
- brief fewer ADRs once the briefs exceed the budget, cataloguing the rest as the role pack's `briefs_left_out` does;
- lower the default budget.

Write the choice as a new ADR that refines ADR-0049, with the measurements in its context and the other fixes under its alternatives.

If the choice cuts briefs, it changes what the operator decided on TH-0032. In that case, ask on a thread on S-0261 before accepting the ADR, with your recommendation first, and wait for the answer.

This task waits for nothing. Every other task builds or describes what this ADR decides.

## Done when

- A new ADR under `design/adrs` refines ADR-0049, and its context gives the measured sizes of the four packs.
- The ADR names the fix, and its alternatives say why the other fixes were not chosen.
- If the fix overrides TH-0032, the operator's answer on the thread is quoted in the ADR.
- `flai check --strict` and the markdown lint pass.

## Notes
