---
id: ADR-0217
title: "A pack primed for the orchestrator's reasons loads the role's conventions with each section tagged all and other words kept only when one of those words is a reason of the run, and per reason the thread, item, figures, and digest answers it needs, with no design briefs and no catalog"
status: accepted
date: 2026-10-29
supersedes: []
superseded_by: []
refines: [ADR-0047, ADR-0075, ADR-0194, ADR-0209]
---

# ADR-0217 A pack primed for the orchestrator's reasons loads the role's conventions with each section tagged all and other words kept only when one of those words is a reason of the run, and per reason the thread, item, figures, and digest answers it needs, with no design briefs and no catalog

## Context

[ADR-0209](0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md) §3 starts the orchestrator for named reasons and primes it for them: the role's conventions once and, per reason, what it needs, instead of the general pack. On 2026-10-29 the general orchestrate pack here was 175,725 bytes. The role's seven conventions alone were 64,370 bytes of that, already over a quarter, which S-0537's second criterion sets as the bound for a pack primed for one thread. So the conventions must be cut by reason, not only the briefs left out.

The cut has to come from the conventions' own text, as everything else in a pack does ([ADR-0047](0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md)). Heading topics would do it, but a topic on a heading narrows the section for every pack that reads the file. `work-management.md` is read by every story's agent. A heading topic such as `review` would take its rules out of every story's pack whose topics lack it.

## Decision

**A pack primed for the orchestrator's reasons keeps a convention section tagged `all` beside other words only when one of those words names a reason of the run, or the role's topic. It loads, per reason, what that reason needs. It leaves out the design briefs, the ranked sections, and the catalog.** This refines ADR-0047, [ADR-0075](0075-the-planner-the-orchestrator-and-the-analyzer-prime-by-role-plan-orchestrate-or.md), [ADR-0194](0194-a-story-s-pack-loads-the-digest-entries-its-topics-or-claim-select-ranked.md), and ADR-0209 §3; the rest of each stands.

1. **Reasons.** A reason is `<kind>:<id>` or a bare `<kind>`. The kinds are `thread:TH-nnnn`, `review:S-nnnn`, `planned:E-nnnn` or `planned:S-nnnn`, `draft:S-nnnn`, `release`, and `ready`. `flai prime --role orchestrate --reason <reason>`, repeatable, and MCP `prime` with role `orchestrate` and `reasons` take them. An unknown kind, a kind given without the ID it needs or with an ID it does not take, and an ID of the wrong type are refused, naming the known kinds. Reasons are refused with any other role, and with `--epic` or `--story`.
2. **Sections.** In a reason pack, a convention section whose effective topics are `all` alone is kept, as in every pack. A section whose effective topics are `all` and other words is kept only when one of those words is a kind of the run's reasons, or `orchestration`, the role's topic. `all` keeps the section in every other pack exactly as before, since `all` matches every pack (ADR-0047). The words are the reason kinds and the built-in strategic roles' topics (`planning`, `orchestration`, `analysis`, `digest`). `flai check` accepts them on a heading as it accepts a story's topics, so a project made from the template, whose baseline carries them, checks clean.
3. **What each reason loads.** A `thread` loads the thread with its entries and the item it is on. A `review` loads the story and its last `flai verify` result. A `planned` loads the item planned, the items the planner wrote under it, and the planner's open thread on it. A `draft` loads the draft story. A `release` loads `release_evaluate`'s figures. A `ready` loads the board and `promote_candidates`' figures. Each item loads once, however many reasons name it.
4. **Digest.** The digest's answers and recommendations ranked against the reasons' threads, the threads of `thread` reasons and the planner threads `planned` loads, are taken as ADR-0194 takes them for the threads awaiting the operator, within an eighth of the budget.
5. **Nothing else.** The pack holds no design or tech briefs, no ranked sections, no open-issues table, and no catalog: `doc_search` and `doc_get` find the rest. It is fitted to the budget and cut into parts as every pack is ([ADR-0104](0104-the-mcp-prime-tool-returns-a-context-pack-in-parts-of-at-most-40-000-bytes-of.md)).
6. **The start prompt.** `flai serve` gives each reason the form `prime` takes beside its line. A move to review is `review:S-nnnn`; a move to backlog and a draft's edit are `draft:S-nnnn`; a plan is `planned:<id>`; a thread is `thread:TH-nnnn`; a release due is `release`; the ready column draining and the order changing are `ready`. The prompt tells the orchestrator to prime with them. A run that holds an `asked` or `schedule` reason, or an event with no form, looks at the whole project, so it primes with the general pack, as a run by hand does.

## Consequences

- A run for one thread reads its conventions cut to what threads need, the thread, and its item: under a quarter of the general pack, as S-0537's test measures.
- The baseline conventions gain heading comments naming `all` and a reason or a role's topic. Every pack but a reason pack prints them as before.
- A reason kind the trigger gains later needs its loader and its words in the conventions. The refusal naming the known kinds makes a missing one visible at the first run.

## Alternatives considered

- **Heading topics without `all`.** Simplest. But they narrow `work-management.md`'s rules for every story's pack, and the strategic-agents sections for the planner, the analyzer, and the curator.
- **A second heading comment, `<!-- reasons: ... -->`.** No other pack would see it. But the topics parser, the heading slugs, and every reader of headings would need to strip a new comment, for what `all` beside a word already expresses.
- **A table in flai of which headings each reason needs.** Keeps the conventions unchanged. But the rule would live away from the text it cuts, and a project's own headings could not take part.
