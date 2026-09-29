---
id: T-0523
type: task
nature: feature
title: Briefs for topic-selected design and stepped ADRs, and the decision-sentence check
status: done
parent: S-0146
owner: alex
created: 2026-09-29T03:19:38Z
updated: 2026-09-29T03:33:54Z
transitions:
  - to: ready
    at: 2026-09-29T03:19:45Z
    by: agent-S-0146
  - to: in-progress
    at: 2026-09-29T03:33:54Z
    by: agent-S-0146
  - to: done
    at: 2026-09-29T03:33:54Z
    by: agent-S-0146
stream: S-0146
tags: []
touches: [flai/internal/context, flai/internal/check]
---
# T-0523 Briefs for topic-selected design and stepped ADRs, and the decision-sentence check

## Work

What the story, its epic, and its tasks name loads whole. Topic-selected design and tech files load as a brief: title, first paragraph, heading outline, and how to fetch a section. ADRs reached one step away (links, refines, refined by), or chosen by topics, load as ID, title, and decision sentence. `flai check` warns about an ADR whose `## Decision` does not open with a sentence.

## Done when

- Behaviour tests cover a brief of a design file, an ADR decision sentence, and the check warning.

## Notes
