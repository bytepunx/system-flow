---
id: E-0025
type: epic
nature: feature
title: A digest of the archive, kept by an agent and decayed as the design absorbs it, primes a story's agent with what past stories decided and what the operator answered, so the same question is not asked twice
status: backlog
owner: alex
created: 2026-10-16T04:28:05Z
updated: 2026-10-16T06:00:41Z
transitions: []
tags: [cli, template]
topics: [cli, planning, analysis, conventions]
touches: [flai/internal/context, flai/internal/serve, flai/cmd, design/system/agent-context.md, design/system/strategic-agents.md, design/system/metrics.md, design/conventions/agent-narrative.md, docs/users/flai.md, docs/operators/index.md, template]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 353
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 99
          output: 1728
          cache_read: 36445899
          cache_write: 27120
          cost: 9.0233
---
# E-0025 A digest of the archive, kept by an agent and decayed as the design absorbs it, primes a story's agent with what past stories decided and what the operator answered, so the same question is not asked twice

## Outcome

The archive holds hundreds of accepted stories, their narratives' decisions and open questions, the threads the operator answered, and the issues closed, and nothing summarizes any of it into what primes the next agent. A story's pack selects conventions, design, and ADRs by topics and links within a budget ([ADR-0047](../../../design/adrs/0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md), [ADR-0049](../../../design/adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md), [agent-context.md § As built](../../../design/system/agent-context.md#as-built)); what a past story learned reaches an agent only when someone wrote it into a design document. The analyzer's report of 2026-10-15 shows the cost: recommendations wait a day or more for the operator, and the questions recur, a secret's location, a write under `.git`, a scope choice ([I-0151](../../../design/issues/I-0151-a-pending-recommendation-waits-a-day-or-more-for-the-operator-s-confirmation-while-its-story-idles-in-progress.md)). Beads decays closed issues into summaries and keeps durable memories that `bd prime` injects; Gas Town lets an agent query its predecessor sessions. This epic gives system-flow a digest: entries drawn from the archive and the threads, each citing its story, kept by an agent run, decayed as the design absorbs them, and primed by topic within the budget.

- [ ] A digest entry is one fact a story learned, with the story, thread, or issue it comes from, the topics and paths it bears on, and the date: a decision from a narrative's `## Decisions`, an answer the operator gave on a thread, a recommendation they confirmed or overruled and why, an issue's cause and fix. The design decides, by an ADR, where the digest lives within the layout and how an entry is written, so that `doc_search` finds it and the dashboard renders it.
- [ ] An agent run keeps the digest: after each acceptance, or on a schedule, it reads what was archived since its last run and the threads that closed, writes or merges entries, and records each run in an activity document as the analyzer does; it is a role under the strategic-agents convention, with its host action and permission, and writes nothing but the digest and its document.
- [ ] Entries decay: when a design document, a convention, or an ADR now records what an entry says, the run marks the entry absorbed and drops it from priming, naming the document; entries no story has needed within a window the manifest sets are compacted into one line each; nothing is deleted from history.
- [ ] `flai prime --story` and the MCP `prime` include the entries the story's topics and touches select, ranked with the rest within the budget and marked as the digest's, and `flai prime --role` gives the orchestrator the entries on the questions it answers, so a recommendation can cite the answer the operator gave last time as its source.
- [ ] An operator's answer on a thread becomes an entry at once, before the digest's next run, so the next story that meets the question finds it.
- [ ] `metrics.md` defines, by an ADR, the digest's measure: threads opened to the operator per completed story, and how many of them asked what an entry already answered, reported by `flai stats` per week and in the dashboard, so the analyzer can see whether the digest cut the waits I-0151 counts.
- [ ] An experiment story measures the effect on this repository: the threads-to-operator rate over two weeks before the digest primes and two weeks after, with its results document under `design/experiments/`.
- [ ] `docs/users/flai.md`, `docs/operators/index.md`, the convention `agent-narrative.md`, and `design/system/agent-context.md` describe the digest and how an agent uses and feeds it, and the template ships the role, its convention, and the manifest settings.

## Stories

To be drafted by the planner. Order: the entry format and its home first, with the ADR; then a one-off build of the digest from today's archive, so there is something to prime with; then priming by topic within the budget; then the keeping run as a role, with absorption and compaction; then the operator's answers as entries at once; then metrics and the experiment; documentation and the template with each story.
- S-0428 A digest entry is one markdown file under design/digest that flai writes, lists, and checks, and an ADR records where it lives and how it is written
- S-0429 flai digest harvest lists what the archive and the threads hold since a date, and a one-off seed fills the digest from the operator's answers and the closed issues so there is something to prime with
- S-0430 flai prime gives a story's agent the digest entries its topics and touches select, and the orchestrator the entries on the threads it answers, ranked within the budget and marked as the digest's
- S-0431 The curator, a strategic agent behind the curate host action, takes in what the archive and the threads gained since its last run as digest entries written through flai alone
- S-0432 The curator runs after each acceptance and on a schedule, marks an entry the design now records absorbed, and compacts an entry no one has cited within the window
- S-0433 An operator's answer on a thread becomes a digest entry at once, before the curator's next run, so the next story that meets the question finds it
- S-0434 flai stats reports per week the threads opened to the operator per completed story and how many asked what a digest entry already answered, and the dashboard charts it
- S-0435 Experiment: does priming with the digest cut the threads opened to the operator per completed story on this repository, two weeks before against two weeks after

## Notes

- From the comparison of system-flow with Beads and Gas Town on 2026-10-16, and from the analyzer's bottleneck report of 2026-10-15, finding 2. Split from "workflow formulas and a memory digest of the archive" because the two share no code; the formulas are their own epic, and this one can be built on today's three roles without waiting for it.
- The planner reads first: `design/system/agent-context.md` § As built and § Fitting the pack to a budget; `design/conventions/agent-narrative.md`; `design/system/strategic-agents.md` for how the analyzer runs and writes; ADR-0047, ADR-0049, ADR-0091, ADR-0099, ADR-0104; `flai/internal/context` for the pack.
- Open to the designer: whether the digest is a folder of markdown entries, one file per topic, or front matter on the archived items themselves; whether the analyzer keeps the digest or a role of its own does; what the compaction window is.
- Safety: an entry never carries a secret's value, only where the operator said it lives; the run's prompt says so and `flai check` refuses an entry that matches a key's shape.
