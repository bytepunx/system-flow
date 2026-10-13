---
id: T-1520
type: task
nature: feature
title: A refused acceptance lists the markdown the fix left as editor links, and the operator accepts again from the same page once they are saved
status: done
parent: S-0378
owner: alex
created: 2026-10-13T16:57:27Z
updated: 2026-10-13T18:32:10Z
transitions:
  - to: ready
    at: 2026-10-13T18:19:45Z
    by: agent-S-0378
  - to: in-progress
    at: 2026-10-13T18:19:45Z
    by: agent-S-0378
  - to: done
    at: 2026-10-13T18:32:10Z
    by: agent-S-0378
stream: S-0378
tags: [dashboard]
touches: [flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, "flaiover/src/routes/api/items/[id]/accept/+server.ts", "flaiover/src/routes/api/items/[id]/accept/accept.test.ts", flaiover/src/routes/board/+page.svelte, "flaiover/src/routes/items/[id]/+page.svelte", flaiover/src/lib/filelinks.ts, flaiover/src/lib/filelinks.test.ts, flaiover/src/lib/components/DismissibleNotice.svelte, flaiover/src/lib/components/DismissibleNotice.svelte.test.ts]
after: [T-1516, T-1517, T-1519]
usage:
  source: log
  seconds: 745
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 170
      output: 64760
      cache_read: 9869257
      cache_write: 264545
      cost: 4.7704
---
# T-1520 A refused acceptance lists the markdown the fix left as editor links, and the operator accepts again from the same page once they are saved

## Work

- Pass the refusal's data from `accept.run` (T-1516) through the error event of `/api/items/[id]/accept` (`+server.ts`), beside its message.
- In `Review.svelte`, show a refusal that carries findings as a list: each file and line opens the editor (T-1518) at that line, with the story for a `checkout: story` finding. Show the rule and the message beside it. Keep the refusal's text for any other refusal, with its paths linked through `filelinks.ts` (T-1519).
- Keep the Accept button on the page after a refusal. Mark each finding's file as saved when the editor's `onsaved` fires, and enable Accept again then, so the operator accepts from the same page with no reload.
- Link the paths in the refusal notices of the board (`routes/board/+page.svelte`) and the item page (`routes/items/[id]/+page.svelte`) the same way.
- Waits for T-1516, whose findings it shows, for T-1517, whose endpoints the editor calls, and for T-1519, whose `filelinks.ts` it uses.

## Done when

- `accept.test.ts` covers a refusal whose error event carries the findings.
- `Review.svelte.test.ts` covers a refusal listed as editor links, a story finding opening the editor with the story, and Accept enabled again after `onsaved`.
- `scripts/flai.sh test` on these paths passes.

## Notes

- Criterion 4 of S-0378, and criterion 1 for refusals.
