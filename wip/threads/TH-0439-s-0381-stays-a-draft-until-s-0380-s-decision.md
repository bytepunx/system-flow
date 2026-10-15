---
id: TH-0439
title: S-0381 stays a draft until S-0380's decision
anchor:
  path: wip/kanban/stories/S-0381-a-project-s-design-tokens-are-a-dtcg-file-the-manifest-names-beside-its-guidance-and-flai-validates-them-renders-them-as-css-reports-their-contrast-and-serves-them-to-the-dashboard.md
  item: S-0381
status: resolved
participants: [orchestrator, alex, planner-S-0381]
created: 2026-10-13T17:53:25Z
updated: 2026-10-14T23:56:14Z
---

# TH-0439 S-0381 stays a draft until S-0380's decision

On wip/kanban/stories/S-0381-a-project-s-design-tokens-are-a-dtcg-file-the-manifest-names-beside-its-guidance-and-flai-validates-them-renders-them-as-css-reports-their-contrast-and-serves-them-to-the-dashboard.md.

## Entries

### 2026-10-13T17:53:25Z orchestrator
**Recommendation:** keep S-0381 a draft until S-0380 is decided, then finalize it as its planner replans it.

flai lists it complete (45m, 173.08 USD/week), but its criteria rest on assumptions S-0380 confirms or corrects (TH-0438 § Assumptions): the design system as one markdown document with tokens in its front matter, and the Tailwind v4 namespaces. If I finalize it now, its planner can no longer rewrite those words without asking. S-0380's T-1531 names on TH-0438 what each draft must change. It is held behind S-0380 by its `after` meanwhile, so nothing waits on it.

### 2026-10-14T23:11:49Z alex
I would prefer we not store front-matter in tokens. In the precursor research, my request was that we instead design file.get and file.save (analogs to doc.get and doc.save) that are capable of working outside of markdown to store the tokens in a different format and also allow for expansion into features that give the operator design tool surfaces, component viewers/testers, and the ability to preview UI/UX without the need to complete a story to see the outcome visually

### 2026-10-14T23:55:03Z planner-S-0381
S-0381 is replanned to match your answer. The tokens are a DTCG file the manifest names (`design.tokens`, default `design/system/design-system.tokens.json`), beside the guidance in `design/system/design-system.md`. No tokens go in front matter (ADR-0146).

`file.get` and `file.save` are the next story's, as § Epic of `design/system/design-system-tool.md` sets out (ADR-0147). Kept general, they can later serve design tool surfaces, component viewers, and previews of a branch, not only the token file. S-0381 needs neither of them. The plan and the task rewrites I propose are on TH-0453. S-0381 stays a draft.

Source: design/system/design-system-tool.md § Epic

### 2026-10-14T23:56:14Z orchestrator
Resolved: S-0380's decision is in (ADR-0146 to ADR-0152), and S-0381 is replanned to the operator's answer here: tokens in a DTCG file the manifest names, nothing in front matter (TH-0453). file.get and file.save (ADR-0147) get their own story from E-0020's planner (TH-0438).
