---
id: TH-0936
title: "Plan for S-0627: render markdown with html: false; needs a cost of delay input (recommend penalty 100 USD a week)"
anchor:
  path: wip/kanban/stories/S-0627-the-dashboard-renders-repository-markdown-with-markdown-it-html-true-and-no.md
  item: S-0627
status: resolved
participants: [planner-S-0627, orchestrator]
created: 2026-11-09T17:12:19Z
updated: 2026-11-09T17:18:29Z
---

# TH-0936 Plan for S-0627: render markdown with html: false; needs a cost of delay input (recommend penalty 100 USD a week)

On wip/kanban/stories/S-0627-the-dashboard-renders-repository-markdown-with-markdown-it-html-true-and-no.md.

## Entries

### 2026-11-09T17:12:19Z planner-S-0627
**Question for the operator:** S-0627 has no cost of delay inputs and no epic, so `flai cod` refuses it.

Recommended: `penalty_per_week: 100` USD on S-0627.

- No instance has been observed. But any document, thread entry, message, or forge review an agent writes can run script in the operator's session. That session can accept, move, and edit items.
- It sits between S-0603's 100 and S-0598's 150 on TH-0888. The harm is a risk, not a cost paid every week, so I put it below S-0598.
- No revenue rides on it, so a penalty fits.

Alternative: `penalty_per_week: 30`, as S-0606 has, if you count the trusted-LAN boundary as enough for now.

## Plan

Fix: markdown-it with `html: false` in `flaiover/src/lib/markdown.ts`. Every `{@html}` site, and `markdownInto` (which `Threads.svelte` uses), takes its HTML from `render()`. So one change covers all 13 sinks. No document under `design/` or `docs/` uses raw HTML outside code. A sanitizer such as DOMPurify was weighed and not chosen: it adds a dependency and an allow-list for HTML nothing needs.

| Layer | Task | Touches | after |
|-------|------|---------|-------|
| 1 | T-2752 Render with raw HTML escaped; tests for `<script>`, `onerror`, `javascript:` links, checkboxes and permalinks kept | `flaiover/src/lib/markdown.ts`, `flaiover/src/lib/markdown.test.ts` | none |
| 1 | T-2753 Say so in `design/tech/markdown.md`, `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md` | those three | none |
| 2 | T-2754 Close I-0294 with what fixed it | the issue file, `design/issues/summary.md` | T-2752, T-2753 |

Forecast: 40m, raised from `flai forecast`'s 16m for the flaiover tiers the size does not count. Delivery 2026-11-09T18:33:00Z.

## Assumptions

- Escaping raw HTML is acceptable. Nothing in the repository renders differently except text that was raw HTML.
- T-2752 also installs the `link_open` rule once, rather than on every `render()` call. Today each call wraps the rule again, so the chain grows with every render.
- I-0294's Impact names `svelte.config.js`, but flaiover has none. T-2754's close reason says so.

## Proposed, not created

- A separate story: a Content-Security-Policy on the dashboard's responses, in `flaiover/src/hooks.server.ts`, as defence in depth. Mermaid's injected `<style>` and shiki's inline styles mean it needs `style-src 'unsafe-inline'` or nonces. That is design work beyond I-0294's cause. Say if you want it drafted.
- The analysis also names the fixed CSRF header (`flaiover/src/lib/server/auth.ts`) and the 0.0.0.0 bind (`flai/cmd/dashboard.go`). These are other causes, out of this story.

### 2026-11-09T17:15:52Z orchestrator
**Plan approved.** I set `penalty_per_week: 100` USD on S-0627, as you recommended. It is S-0627's first cost of delay input, and neither the story nor an epic had one. The planner works out the value on its next run.

**Why the plan fits:**

- One change in `render()` covers every `{@html}` site and `markdownInto`, so the fix meets criterion 1 at its cause.
- T-2752 and T-2753 share no path and run together. T-2754 waits for both, since the close reason names what fixed it.
- Each touch is a file, and each task's touches are in the story's.
- The forecast of 40m is `flai forecast`'s 16m, raised for the flaiover tiers, with the reason stated.

**Why 100, not 30:** nothing on this thread or in S-0627 gives a reason for the alternative. ADR-0018's trusted-LAN boundary does not cover content agents write: the analysis notes `flai dashboard` binds `0.0.0.0` by default.

**Left out of this plan:** the Content-Security-Policy story you offer is a new story, so it is the operator's call. I ask it on a thread of its own on the analysis's Finding R1, so it does not hold S-0627. The CSRF header and the bind are other causes, as you say.

Resolving this thread.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-09T17:16:17Z orchestrator
Resolved: Plan approved; penalty_per_week 100 USD set on S-0627 as the planner recommended. The CSP story offer is asked separately.

### 2026-11-09T17:16:26Z planner-S-0627
The orchestrator set `penalty_per_week: 100` on S-0627 at 2026-11-09T17:15:40Z. The cost of delay value, 100 USD a week, is now written, and the reason is under `### Planning` in the story's Notes. The plan is complete. The one open point is the Content-Security-Policy story proposed above: say if you want it drafted.

### 2026-11-09T17:18:00Z orchestrator
Noted: the value is 100 USD a week, and the plan is complete. The Content-Security-Policy offer is asked on TH-0938, on the analysis's Finding R1. A new story is the operator's call, and that thread keeps the question off S-0627. Resolving this thread again.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-09T17:18:29Z orchestrator
Resolved: Plan complete with value 100 USD a week; the CSP story offer is asked on TH-0938.
