---
id: TH-0033
title: "flai check --strict fails on main only because TH-0032 is open on archived S-0146: may I move it to S-0149?"
anchor:
  path: wip/kanban/stories/S-0147-agents-fetch-a-design-section-on-demand-doc-search-and-a-heading-on-doc-get-and-flai-doc-show.md
  item: S-0147
status: resolved
participants: [agent-S-0147, alex]
created: 2026-09-29T05:09:11Z
updated: 2026-09-29T05:14:48Z
---

# TH-0033 flai check --strict fails on main only because TH-0032 is open on archived S-0146: may I move it to S-0149?

On wip/kanban/stories/S-0147-agents-fetch-a-design-section-on-demand-doc-search-and-a-heading-on-doc-get-and-flai-doc-show.md.

## Entries

### 2026-09-29T05:09:11Z agent-S-0147
S-0147 is done apart from this. Behavior, integration, lint, the template smoke, markdown lint, and the installer smoke all pass. The one failure is the repository check in the smoke tier:

```text
wip/threads/TH-0032-...md:7: warning: threads.archived: TH-0032 is answered but S-0146 is archived; resolve it or move it
```

TH-0032 is S-0146's thread. It is not S-0147's. Your last answer there asked for an ADR refining ADR-0049 and then a backlog story. S-0149 is that story. Its goal says the refinement "probably needs a superseding ADR", and that ADR has not been written yet.

Options:

A (recommended). I move TH-0032 to S-0149. The ADR request stays open where the work is, and the check passes.
B. I resolve TH-0032, noting that S-0149 carries the ADR. The check passes, but the thread is closed before the ADR exists.
C. Leave it. S-0147 goes to review with that criterion unticked, and the notes say why.

I will wait for your answer before moving S-0147 to review.

### 2026-09-29T05:14:34Z alex
A

### 2026-09-29T05:14:48Z agent-S-0147
Resolved: A: TH-0032 moved to S-0149; flai check --strict clean.
