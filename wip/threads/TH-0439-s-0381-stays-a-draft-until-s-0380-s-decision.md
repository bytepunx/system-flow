---
id: TH-0439
title: S-0381 stays a draft until S-0380's decision
anchor:
  path: wip/kanban/stories/S-0381-a-project-s-design-system-is-a-document-whose-front-matter-holds-its-tokens-and-flai-validates-it-renders-it-as-css-and-serves-it-to-the-dashboard.md
  item: S-0381
status: answered
participants: [orchestrator, alex]
created: 2026-10-13T17:53:25Z
updated: 2026-10-14T23:11:49Z
---

# TH-0439 S-0381 stays a draft until S-0380's decision

On wip/kanban/stories/S-0381-a-project-s-design-system-is-a-document-whose-front-matter-holds-its-tokens-and-flai-validates-it-renders-it-as-css-and-serves-it-to-the-dashboard.md.

## Entries

### 2026-10-13T17:53:25Z orchestrator
**Recommendation:** keep S-0381 a draft until S-0380 is decided, then finalize it as its planner replans it.

flai lists it complete (45m, 173.08 USD/week), but its criteria rest on assumptions S-0380 confirms or corrects (TH-0438 § Assumptions): the design system as one markdown document with tokens in its front matter, and the Tailwind v4 namespaces. If I finalize it now, its planner can no longer rewrite those words without asking. S-0380's T-1531 names on TH-0438 what each draft must change. It is held behind S-0380 by its `after` meanwhile, so nothing waits on it.

### 2026-10-14T23:11:49Z alex
I would prefer we not store front-matter in tokens. In the precursor research, my request was that we instead design file.get and file.save (analogs to doc.get and doc.save) that are capable of working outside of markdown to store the tokens in a different format and also allow for expansion into features that give the operator design tool surfaces, component viewers/testers, and the ability to preview UI/UX without the need to complete a story to see the outcome visually
