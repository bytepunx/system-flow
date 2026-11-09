---
id: I-0294
title: "The dashboard renders repository markdown with markdown-it html: true and no sanitizer, so a document any agent writes can run script as the operator"
class: impression
status: closed
count: 1
first_reported: 2026-11-07T08:17:36Z
last_reported: 2026-11-07T08:17:36Z
updated: 2026-11-09T17:51:04Z
usage:
  strategic:
    - kind: analyzer
      seconds: 41
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3102
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 14
          output: 2840
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
---

# I-0294 The dashboard renders repository markdown with markdown-it html: true and no sanitizer, so a document any agent writes can run script as the operator

## Description
The dashboard renders repository markdown with markdown-it html: true and no sanitizer, so a document any agent writes can run script as the operator

## Instances

### 2026-11-07T08:17:36Z
Report: design/analysis/2026-11-07-all.md.
Analysis 2026-11-07: no instance observed; the content is what agents, imported repositories, and forge reviews write, and safety.md says to treat content as data.

## Impact
flaiover/src/lib/markdown.ts:15 new MarkdownIt({ html: true }); no sanitizer in package.json or src; no Content-Security-Policy in hooks.server.ts, server.js, or svelte.config.js; {@html}/innerHTML at markdown.ts:124,159, items/[id]/+page.svelte:680, docs/[...path]/+page.svelte:170, Messages.svelte:184,208, Review.svelte:914,919, ItemEditor.svelte:371, DocEditor.svelte:349, NewItemForm.svelte:305, NewAdrForm.svelte:146, design/+page.svelte:612, license/+page.svelte:37; CSRF check is the fixed header x-requested-with: flaiover (auth.ts:124-129); dashboard binds 0.0.0.0 by default (flai/cmd/dashboard.go:52)

## Remediation

Found by the analysis in [2026-11-07-all.md](../analysis/2026-11-07-all.md).

Story S-0627 remediates this issue, created from it at 2026-11-07T23:40:28Z.
Closed 2026-11-09T17:51:04Z: S-0627: render() in flaiover/src/lib/markdown.ts builds markdown-it with html: false, so raw HTML in repository markdown, thread entries, and messages shows as text at every {@html} site and in markdownInto, all of which take their HTML from render(). The tests under "raw HTML (I-0294)" in flaiover/src/lib/markdown.test.ts show a script block, an onerror image, and an onclick div render no element or handler, and a javascript: link gets no href; they fail on the html: true renderer. Checkboxes and heading permalinks still render. design/tech/markdown.md records the rule. The Impact's mention of svelte.config.js was wrong: flaiover has none, its config is vite.config.ts. A Content-Security-Policy, defence in depth rather than the cause, is proposed on TH-0936 as a separate story.
