---
id: D-0144
title: The dashboard drops flai's archived flag from board cards, so the done lane shows archived items it should hide while the clone lags its remote's tags
kind: cause
sources: [I-0060]
topics: [flaiover, dashboard, sveltekit, code, all]
paths: [design/issues/I-0060-the-dashboard-drops-flai-s-archived-flag-from-board-cards-so-the-done-lane-shows-archived-items-it-should-hide-while-the-clone-lags-its-remote-s-tags.md, design/issues/summary.md, flaiover/src/lib/server/board.ts, flaiover/src/lib/server/repo-channel.test.ts, design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md]
date: 2026-10-04
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0144 The dashboard drops flai's archived flag from board cards, so the done lane shows archived items it should hide while the clone lags its remote's tags

I-0060 was closed on 2026-10-04.

The cause:

```text
The dashboard drops flai's archived flag from board cards, so the done lane shows archived items it should hide while the clone lags its remote's tags
```

The fix:

```text
S-0247: flaiover's board.ts now carries flai's archived flag onto each board card, so doneLane() leaves out the done, archived cards while the clone lags its remote's tags; repo-channel.test.ts checks it through board() with a fake board.get rather than a mocked /api/board.
```

Source: I-0060.
