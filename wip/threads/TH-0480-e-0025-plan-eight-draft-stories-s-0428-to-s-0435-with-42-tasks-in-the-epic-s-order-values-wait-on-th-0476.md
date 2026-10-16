---
id: TH-0480
title: "E-0025 plan: eight draft stories, S-0428 to S-0435, with 42 tasks, in the epic's order; values wait on TH-0476"
anchor:
  path: wip/kanban/epics/E-0025-a-digest-of-the-archive-kept-by-an-agent-and-decayed-as-the-design-absorbs-it-primes-a-story-s-agent-with-what-past-stories-decided-and-what-the-operator-answered-so-the-same-question-is-not-asked-twice.md
  item: E-0025
status: open
participants: [planner-E-0025]
created: 2026-10-16T06:14:47Z
updated: 2026-10-16T07:51:59Z
---

# TH-0480 E-0025 plan: eight draft stories, S-0428 to S-0435, with 42 tasks, in the epic's order; values wait on TH-0476

On wip/kanban/epics/E-0025-a-digest-of-the-archive-kept-by-an-agent-and-decayed-as-the-design-absorbs-it-primes-a-story-s-agent-with-what-past-stories-decided-and-what-the-operator-answered-so-the-same-question-is-not-asked-twice.md.

## Entries

### 2026-10-16T06:14:47Z planner-E-0025
I drafted eight stories for E-0025, all in the backlog as drafts. Each has its tasks, file-level touches, and a forecast. The cost of delay values wait on TH-0476, question 1. Every story carries its own docs and template changes, so there is no docs story. Each story's `### Planning` notes say where every touch came from and why each figure stands.

## Stories and order

| Story | What | After | Forecast | flai's |
|-------|------|-------|----------|--------|
| S-0428 | Entry format and home, ADR, `flai/internal/digest`, `flai digest new/edit/list/show`, `digest.entry` and `digest.secret` checks | none | 51m | 51m |
| S-0429 | `flai digest harvest` and a one-off `flai digest seed` from the operator's answers and closed issues, run here | S-0428 | 45m | 32m |
| S-0430 | Priming: the digest step in a story's pack, the orchestrator's pack, citing an entry as a `source`, conventions | S-0429 | 48m | 48m |
| S-0431 | The curator: host action `curate`, run, prompt, guard, definition, activity document, on demand | S-0430 | 1h23m | 1h23m |
| S-0432 | The curator after each acceptance and on `digest.schedule`; `absorb`, `cited`, `compact`; `digest.compact_after` | S-0431 | 55m | 55m |
| S-0433 | An operator's answer becomes an entry at once, from the CLI, MCP, and the dashboard | S-0432 | 29m | 29m |
| S-0434 | Metric by an ADR: threads to the operator per completed story, repeats; `flai stats` and a chart | S-0433 | 39m | 39m |
| S-0435 | Experiment: two weeks before against two weeks after the digest primes | S-0434 | 25m | 13m |

The total is 6h15m of agent time. flai plays the last code story out to 2026-10-17T06:02Z, behind the stories of E-0021 to E-0024 ahead of these in the pull order. S-0435 cannot start until 14 days after S-0430 is accepted, so its delivery is 2026-10-31 at the earliest.

## Tasks and layers

- **S-0428:** T-1807; then T-1808; then T-1809, T-1810, T-1811; then T-1812.
- **S-0429:** T-1813, T-1814; then T-1815; then T-1816; then T-1817.
- **S-0430:** T-1818; then T-1819; then T-1820; then T-1821, T-1822.
- **S-0431:** T-1823; then T-1824; then T-1825, T-1826; then T-1827, T-1828; then T-1829; then T-1830.
- **S-0432:** T-1831; then T-1832, T-1833; then T-1834, T-1835; then T-1836.
- **S-0433:** T-1837; then T-1838, T-1839, T-1840; then T-1841.
- **S-0434:** T-1842; then T-1843; then T-1844, T-1845; then T-1846.
- **S-0435:** T-1847; then T-1848.

No two tasks of one layer share a path. `flai check --strict` reports nothing.

## Assumptions

1. TH-0476's recommended answers hold. Questions 2 to 5 shape S-0428, S-0431, S-0432, and S-0433; the ADR task of each names the question it rests on.
2. The convention `agent-narrative.md` that the epic names does not exist. The narrative's rules are in `work-management.md` and `design/system/agent-narrative.md`, so S-0430 changes those.
3. Design subfolders are not manifest keys. `design/digest` follows `layout.design`, as `design/analysis` does, so S-0428 adds no setting. It does add a row to the layout table in `CLAUDE.md` and its template.
4. `flai digest new` writes in the checkout it runs in. S-0429's seed is therefore committed on its branch. The curator and the operator's answers write in the main checkout, uncommitted, as the analyzer's reports and threads are. S-0431's ADR should say who commits them.
5. Decisions from narratives (381 archived) are harvested but never seeded mechanically. Only a model can judge which are already in the design, so the curator's first run takes them in. The seed writes only the operator's answers and the 138 closed issues.
6. The curator cannot run from a story's own build before release, since `flai serve` runs the installed flai. S-0431 proves the run with a stub harness command, and its first real run follows its release, with the operator's `curate` action on.
7. The `.claude/` writes in S-0431 and S-0432 go through `permission_prompt` and wait on the operator.
8. "Needed" means cited by ID in a narrative, thread, or work item (TH-0476 question 4). S-0430's conventions ask agents to cite `D-nnnn`, and S-0432 counts the citations.

## Proposals

- **Relax the order (not applied):** S-0433 needs only S-0429's answer rule, so it could follow S-0429 and put the operator's answers in the digest about 2h40m sooner. S-0434 reads only S-0428's front matter and the threads, so it could follow S-0428; the measure would then be running before the digest primes, which helps S-0435. I kept the epic's order.
- **Wait for S-0426 (not applied):** S-0426, the agent memory research now in progress, may bear on the entry format and retrieval. Adding S-0426 to S-0428's `after` would let the ADR take its findings in. S-0426 is already in progress, so the wait would be short.
- **Not drafted:** a Curator page on the dashboard, as `/workflow/analyzer` is (S-0228, S-0229), and a digest page listing entries by state. The epic does not ask for them. E-0024's S-0420 makes one page per role, which would give the curator a page once it is a role file.
- **E-0024 overlap (asks a change to finalized-style drafts of another epic):** the curator adds a fourth built-in role. E-0024's S-0414 (role files), S-0416 (guard rules), and S-0419 (the built-in roles on the role path) should then take in the curator too. Alternatively, S-0431 and S-0432 wait for S-0419 and the curator becomes a role file, which is TH-0476 question 3's second alternative. I have not changed E-0024's stories.

## Overlap with other epics

These stories share files with drafts of E-0021 to E-0024, so the holds run them one after the other:

- `flai/internal/guard/guard.go`: S-0431 and S-0432 with S-0393, S-0407, and S-0416.
- `flai/internal/harness/harness.go` and `adapters.go`: S-0430, S-0431, and S-0432 with S-0406, S-0417, and S-0419.
- `flai/internal/serve/agents.go` and `serve.go`: S-0431 and S-0432 with S-0393, S-0394, S-0399, S-0411, and S-0417 to S-0422.
- `flai/internal/hostapi/writes.go`, `flai/cmd/serve_actions.go`, and `flai/internal/manifest/manifest.go` and `settings.go`: S-0431, S-0432, and S-0433 with E-0021's, E-0022's, and E-0024's serve and settings stories.
- `flai/internal/context/strategic.go` and `flai/internal/conventions/conventions.go`: S-0430 and S-0431 with S-0415.

### 2026-10-16T07:51:59Z planner-E-0025
S-0426's finding (TH-0481, `design/system/agent-memory.md` on its branch) bears on this plan. It needs no change to the drafts now.

- It recommends adapting flai's own deterministic index over adopting a memory system, and reaching the archive and threads through this epic's digest under `design/` rather than indexing them directly (its decision 4). S-0428 and S-0430 build exactly that.
- Its follow-ons (section size bound, index kept in memory, stemmer and glossary) change `flai/internal/search` and `flai/internal/context/sections.go`. S-0428's T-1811 only adds tests beside them, so the overlap is small. If those stories land first, S-0430's ranking of entries gets their recall gains for free.
- Its decision 7 (pack room by narrowing the conventions, budget kept at 80 KB) is the constraint S-0430's ADR (T-1818) works within when it sets the digest step's share of the budget.
- It names S-0435 as the check on whether the digest misses what agents ask, which is what S-0435 measures.

The proposal to add S-0426 to S-0428's `after` still stands. The finding is written; what remains is S-0426's acceptance and the operator's answers on TH-0481.
