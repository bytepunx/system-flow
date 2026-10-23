---
id: D-0017
title: "Should Retry queue a ready story's agent when the in-progress limit is full: I have run that command but because I had the older version I was not seeing a way to retry the story"
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

# D-0017 Should Retry queue a ready story's agent when the in-progress limit is full: I have run that command but because I had the older version I was not seeing a way to retry the story

The operator, alex, answered TH-0011 at 2026-09-26T03:09:28Z.

What was asked:

```text
The story page already has a Restart agent button for a failed agent (S-0116). S-0118 changes it to Retry: it moves to the top right of the agent pane, and it hides once clicked. One criterion says Retry should "queue the agent as it would normally do". Today, restarting a failed agent on a story in ready is refused when the in-progress limit is full ("the in-progress limit leaves no room for S-…"). Nothing is queued.
```

The answer:

```text
I have run that command but because I had the older version I was not seeing a way to retry the story. Let's complete the changes requested as part of this story.
```

Source: TH-0011.
