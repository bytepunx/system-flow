---
id: T-1980
type: task
nature: improvement
title: AcceptFindings lists a refused acceptance's findings, each opened on its own checkout, and Review uses it
status: done
parent: S-0465
owner: alex
created: 2026-10-21T07:36:38Z
updated: 2026-10-21T07:54:43Z
transitions:
  - to: ready
    at: 2026-10-21T07:39:04Z
    by: agent-S-0465
  - to: in-progress
    at: 2026-10-21T07:39:05Z
    by: agent-S-0465
  - to: done
    at: 2026-10-21T07:54:43Z
    by: agent-S-0465
stream: S-0465
tags: [flaiover]
touches: [flaiover/src/lib/components/AcceptFindings.svelte, flaiover/src/lib/components/AcceptFindings.svelte.test.ts, design/components/flaiover/src/lib/components/AcceptFindings.examples.json, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts]
usage:
  source: log
  seconds: 938
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 42
      output: 10232
      cache_read: 2072972
      cache_write: 85370
      cost: 1.1302
---
# T-1980 AcceptFindings lists a refused acceptance's findings, each opened on its own checkout, and Review uses it

## Work

`flaiover/src/lib/components/Review.svelte` draws the findings of an acceptance refused for lint the fix left (`data-testid="accept-findings"`, S-0378): each file and line as a button that opens the editor on the finding's own checkout (`storyOf(f)`), its rule and message, and `saved` once the file is saved. The board and the item page need the same list, so move it into a component of its own rather than copy it.

- Add `flaiover/src/lib/components/AcceptFindings.svelte`. Props: the item's id and status, the findings (`Finding[]`, as `findingsOf` reads them), and the keys of the files saved so far, with a callback when a file is saved. It renders what `Review.svelte` renders today, with the same test IDs (`accept-findings`, `accept-finding`, `finding-saved`), and opens each finding with `openEditor` on its own checkout: the story's branch when `story` is set, the main checkout when it is not.
- Use it in `Review.svelte` in place of the inline list. Keep `Review.svelte.test.ts` passing as it is, and add nothing to it unless the move changes what it sees.
- Add `flaiover/src/lib/components/AcceptFindings.svelte.test.ts`: a finding on the story's branch opens the editor with that story, a finding in the main checkout opens it with none, and a saved file is marked.
- Add its examples, `design/components/flaiover/src/lib/components/AcceptFindings.examples.json`, as `design/components/README.md` says (ADR-0164): one set with a finding on the branch and one in the main checkout.

This task waits for no other: it changes only flaiover's components and reads the finding's shape S-0378 defined.

## Done when

- `AcceptFindings.svelte` renders the findings and opens each on its own checkout, and `Review.svelte` uses it.
- Its tests and `Review.svelte.test.ts` pass, and `flai test` passes on the changed paths, `components.examples` included.

## Notes
