---
id: T-0087
type: task
nature: feature
title: "Markdown rendering: markdown-it with anchors and task lists, relative links rewritten to explorer routes, mermaid and shiki lazy-loaded"
status: done
parent: S-0012
owner: alex
created: 2026-09-17T04:24:13Z
updated: 2026-09-17T04:29:50Z
transitions:
  - to: ready
    at: 2026-09-17T04:29:50Z
    by: agent
  - to: in-progress
    at: 2026-09-17T04:29:50Z
    by: agent
  - to: done
    at: 2026-09-17T04:29:50Z
    by: agent
stream: S-0012
tags: [dashboard]
---

# T-0087 Markdown rendering: markdown-it with anchors and task lists, relative links rewritten to explorer routes, mermaid and shiki lazy-loaded

## Work
src/lib/markdown.ts renders markdown on the client with markdown-it, markdown-it-anchor, markdown-it-task-lists; fenced mermaid becomes <pre class="mermaid"> rendered by a lazy-loaded mermaid; other fences highlighted by lazy-loaded shiki with a plain fallback; relative links to .md files become explorer routes and other relative links stay repo-relative; headings get ids.

## Done when
Unit tests cover link rewriting, the mermaid fence, task lists, and anchors.

## Notes
