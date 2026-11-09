---
id: TH-0938
title: Draft a story for a Content-Security-Policy on the dashboard, beside S-0627's fix?
anchor:
  path: design/analysis/2026-11-07-all.md
  heading: "Finding R1: the dashboard renders repository markdown as raw HTML with no sanitizer"
status: open
participants: [orchestrator]
created: 2026-11-09T17:16:24Z
updated: 2026-11-09T17:16:24Z
---

# TH-0938 Draft a story for a Content-Security-Policy on the dashboard, beside S-0627's fix?

On design/analysis/2026-11-07-all.md § Finding R1: the dashboard renders repository markdown as raw HTML with no sanitizer.

## Entries

### 2026-11-09T17:16:24Z orchestrator
**Recommended:** no story now. S-0627's fix, markdown-it with `html: false` in `render()`, removes the cause Finding R1 and I-0294 record. Ask again if a sink outside `render()` turns up.

**What was offered:** planner-S-0627 offered on TH-0936 to draft a separate story: a Content-Security-Policy on the dashboard's responses in `flaiover/src/hooks.server.ts`, as defence in depth. It is design work. Mermaid's injected `<style>` and shiki's inline styles need `style-src 'unsafe-inline'` or nonces.

**Why not now:**

- No instance has been observed. The analysis records the finding as an impression.
- Once S-0627 escapes raw HTML, the 13 sinks the analysis lists get no markup a document wrote.
- A CSP with `'unsafe-inline'` styles adds little against script. Nonces are design work across the shiki and Mermaid renderers.

**Alternative:** say "draft it", and I ask the planner for it once it is a backlog story. A new story is yours to create.

The analysis's two other causes are not part of this question: the fixed `x-requested-with` CSRF header and the `0.0.0.0` bind.
