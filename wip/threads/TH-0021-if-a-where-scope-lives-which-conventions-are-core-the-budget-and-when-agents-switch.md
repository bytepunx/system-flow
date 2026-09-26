---
id: TH-0021
title: "If A: where scope lives, which conventions are core, the budget, and when agents switch"
anchor:
  path: wip/kanban/stories/S-0125-determine-how-to-prime-context-with-relevant-documentation-only.md
  item: S-0125
status: resolved
participants: [agent-S-0125, alex]
created: 2026-09-26T08:08:34Z
updated: 2026-09-26T17:45:33Z
---

# TH-0021 If A: where scope lives, which conventions are core, the budget, and when agents switch

On wip/kanban/stories/S-0125-determine-how-to-prime-context-with-relevant-documentation-only.md.

## Entries

### 2026-09-26T08:08:34Z agent-S-0125
These choices only matter if TH-0020 is answered A. Each one has my recommendation listed first.

**1. Where an ADR's scope comes from.** Accepted ADRs can't be edited.
- **a.** Work it out from the documents that link to the ADR. An ADR is in scope when a selected design section links to it or the story names it; any other ADR is reached by ranking or the catalog. Nothing new is written on ADRs. Design docs already link 37 of the 47 ADRs.
- **b.** Allow an `applies_to:` field to be added to an accepted ADR, the same way `superseded_by` can be.

**2. How the other documents declare scope.** Conventions, `design/system` and `design/tech` files would get an optional `applies_to:` front matter field.
- **a.** It lists paths or manifest project names. A convention marked `applies_to: [always]` is core. A document without the field is never selected by scope, only by links or ranking. `flai check` warns when a design doc or tech file has none. `flai upgrade` sets it on the template's conventions and existing projects add it themselves.
- **b.** Keep one central map file instead, `design/context.yaml`.

**3. Which conventions are core.**
- **a.** Every convention except three. `code-quality.md`, `logging.md` and `telemetry.md` apply only when the claim touches a code project. That saves about 12 KB on stories that only touch design or docs.
- **b.** All of them stay core, and scoping applies only to design, ADRs and tech.

**4. Budget.**
- **a.** 48 KB, about 12k tokens, for everything after the core conventions. It can be changed with `--budget` and a `context.budget` manifest key. Core and scoped conventions are never cut. Linked documents are trimmed to their best-ranked sections before any are dropped, and the catalog is always printed.
- **b.** No budget; load everything that is selected.

**5. When agents switch.**
- **a.** Build the pack and a replay test over the archive first. Switch the harness prompt, `CLAUDE.md` and the template's session-start convention in the story after that, once the replay reaches at least 90% of what archived stories named.
- **b.** Switch as soon as the pack exists.

Reply with something like `1a 2a 3a 4a 5a`, or change any of them.

### 2026-09-26T17:43:40Z alex
1b - instead of "applies_to", lets use "topics" 2a - use "topics" and "[all]" as the convention instead of "applies_to" 3 - default to "topics": [all] and I can go through and review them after 4b 5b

### 2026-09-26T17:45:33Z agent-S-0125
Resolved: 1b, 2a, 3 (every convention starts [all]), 4b, 5b, with topics in place of applies_to; recorded in ADR-0047.
