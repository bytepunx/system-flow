---
id: T-2754
type: task
nature: remediation
title: Close I-0294 with what fixed it
status: backlog
parent: S-0627
owner: alex
created: 2026-11-09T17:11:25Z
updated: 2026-11-09T17:11:25Z
transitions: []
stream: S-0627
tags: [dashboard]
touches: [design/issues/I-0294-the-dashboard-renders-repository-markdown-with-markdown-it-html-true-and-no.md, design/issues/summary.md]
after: [T-2752, T-2753]
---
# T-2754 Close I-0294 with what fixed it

## Work

Run `flai issue close I-0294 --reason "<reason>"`. The reason names the fix: `flaiover/src/lib/markdown.ts` renders with `html: false`, so raw HTML in repository markdown shows as text at every `{@html}` site; and it names the tests in `flaiover/src/lib/markdown.test.ts` that prove it. Add that the issue's note on `svelte.config.js` was wrong, since flaiover has none.

It waits for T-2752, whose fix and tests it names, and for T-2753, so the issue closes only once the rule is written down.

## Done when

- I-0294 is closed and its reason names the change and its tests.
- `design/issues/summary.md` no longer lists I-0294 as open.
- `flai check --strict` is clean.

## Notes
