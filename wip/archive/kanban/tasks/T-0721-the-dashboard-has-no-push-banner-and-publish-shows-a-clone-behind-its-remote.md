---
id: T-0721
type: task
nature: improvement
title: The dashboard has no push banner, and Publish shows a clone behind its remote
status: done
parent: S-0195
owner: arobson
created: 2026-10-02T23:31:09Z
updated: 2026-10-03T00:18:11Z
transitions:
  - to: ready
    at: 2026-10-02T23:31:26Z
    by: agent-S-0195
  - to: in-progress
    at: 2026-10-03T00:07:09Z
    by: agent-S-0195
  - to: done
    at: 2026-10-03T00:18:11Z
    by: agent-S-0195
stream: S-0195
tags: []
touches: [flaiover/src]
after: [T-0719, T-0720]
usage:
  source: log
  seconds: 662
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 140
      output: 38685
      cache_read: 7376370
      cache_write: 147663
      cost: 3.1096
---
# T-0721 The dashboard has no push banner, and Publish shows a clone behind its remote

## Work

Remove the push banner over the board entirely: `UnpushedNotice.svelte`, its test, and `/api/unpushed` go, and the board and item pages stop rendering it. The review page stops asking `/api/unpushed` and stops telling the operator to run `flai push --pending` after an acceptance; it says accepted work reaches the remote when it is published. `PublishBanner` still shows what is accepted and not yet published, and shows a clone behind its remote branch as it shows missing tags: nothing to publish, and the commands to run. Waits for T-0719, whose publish preview shape it shows, and T-0720, whose removed host methods it stops calling.

## Done when

- The board has no push banner and nothing in the dashboard pushes accepted work or offers auto-publish
- The board's Publish banner still lists what is accepted and not yet published, and names the fetch and merge when the clone is behind
- The dashboard's tests cover both

## Notes
