---
id: TH-1258
title: Accept ADR-0312 as written?
anchor:
  path: wip/kanban/stories/S-0805-a-tier-in-the-manifest-belongs-to-a-suite-flai-test-suite-runs-one-and-the.md
  item: S-0805
status: resolved
participants: [agent-S-0805, orchestrator, alex]
created: 2026-11-16T23:52:50Z
updated: 2026-11-17T07:47:56Z
---

# TH-1258 Accept ADR-0312 as written?

On wip/kanban/stories/S-0805-a-tier-in-the-manifest-belongs-to-a-suite-flai-test-suite-runs-one-and-the.md.

## Entries

### 2026-11-16T23:52:50Z agent-S-0805
ADR-0312 (design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) is `proposed`. S-0805's criterion 1 needs it accepted, or superseded by an ADR with your changes.

**Recommendation: accept it as written.** It records what you decided on 2026-11-16.

S-0805 builds two points of it:

1. §1: a tier in `tests:` may give `suite: <name>`, kebab-case; a tier with none is the suite of its own name. `flai test --suite <name>` runs a suite's tiers; a `manual` tier still runs only when named.
2. The template's manifest labels its tiers `static` (today's `markdown`), `behavior` (today's `test`), `integration`, and `smoke`.

For this repository's manifest I will use, with no change to what any command runs:

- `static`: gofmt, vet, golangci-lint, shellcheck, markdown, flaiover-lint, svelte-check
- `behavior`: go-test, settings-index, vitest, flaiover-flai, design-routes, template, flaiover
- `integration`: integration, race
- `smoke`: smoke, gateway

Answer "accept", or say what to change. I start on §1 meanwhile and adjust to your answer.

### 2026-11-16T23:56:53Z orchestrator (recommendation)
**Recommended:** accept ADR-0312 as written. S-0805 then goes on with §1 and the template's four labels, as agent-S-0805 proposes.

**Why:**

- Its Context records your decisions of 2026-11-16: steps name suites, not tiers; a cluster step is skipped where no cluster is configured; this repository changes nothing until its own file says so.
- Its Decision is what E-0034's ten stories S-0805 to S-0815 are planned on. Their plans are approved and finalized, all but S-0815.
- One tension: §8 puts the new keys in the template's file now, and a Consequence says the template should not use them until the release that reads them is the oldest a project runs. S-0806's T-3473 settles it. Listing the keys raises `flai.minimum` on publish, so an older flai stops with a message naming the release instead of a bare parse error.

**Why I am asking you:** accepting an ADR is your decision.

**Alternative:** name what to change. agent-S-0805 then writes a superseding ADR under T-3444, and the stories after it are planned again where they read the change.

Source: design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md § Decision

### 2026-11-17T07:46:25Z alex
Confirmed the recommendation of 2026-11-16T23:56:53Z orchestrator.

Source: design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md § Decision

### 2026-11-17T07:47:56Z agent-S-0805
Resolved: The operator confirmed accepting ADR-0312 as written; flai adr accept 0312 committed ea1e2ea on story/S-0805.
