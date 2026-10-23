---
id: D-0016
title: "Should Retry queue a ready story's agent when the in-progress limit is full: No restart button appears on the failed story I have on the board now (for the record)"
kind: answer
sources: [TH-0011]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flaiover/src, flai/cmd]
date: 2026-09-26
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0016 Should Retry queue a ready story's agent when the in-progress limit is full: No restart button appears on the failed story I have on the board now (for the record)

The operator, alex, answered TH-0011 at 2026-09-26T03:05:33Z.

What was asked:

```text
The story page already has a Restart agent button for a failed agent (S-0116). S-0118 changes it to Retry: it moves to the top right of the agent pane, and it hides once clicked. One criterion says Retry should "queue the agent as it would normally do". Today, restarting a failed agent on a story in ready is refused when the in-progress limit is full ("the in-progress limit leaves no room for S-…"). Nothing is queued.
```

The answer:

```text
No restart button appears on the failed story I have on the board now (for the record). Queue the story as recommended but change the indicator to yellow so it's clear that the story agent is waiting for something.
```

Source: TH-0011.
