---
id: D-0301
title: A story's agent that runs flai stream open is left with the story in ready, so permission_prompt refuses its .claude/ write until it moves the story itself
kind: cause
sources: [I-0097]
topics: [flai, cli, go, code, all]
paths: [design/adrs/0153-flai-stream-open-starts-a-ready-story-as-flai-story-start-does-moving-it-to-in.md, design/adrs/README.md, design/system/flai-cli.md, docs/users/flai-reference.md, docs/users/flai.md, flai/cmd/move_cancel_test.go, flai/cmd/stream.go, flai/cmd/stream_open_test.go, flai/internal/storystart/start.go, design/issues/I-0147-the-mcp-server-s-flai-lags-a-fresh-publish-so-task-done-refuses-and-a-story-s-agent-falls-back-to-the-tree-s-cli.md, design/issues/summary.md, design/issues/I-0097-a-story-s-agent-that-runs-flai-stream-open-is-left-with-the-story-in-ready-so-permission-prompt-refuses-its-claude-write-until-it-moves-the-story-itself.md]
date: 2026-10-14
status: absorbed
absorbed_by: "design/adrs/0153-flai-stream-open-starts-a-ready-story-as-flai-story-start-does-moving-it-to-in.md#decision"
asked_again: []
cited: []
---

# D-0301 A story's agent that runs flai stream open is left with the story in ready, so permission_prompt refuses its .claude/ write until it moves the story itself

I-0097 was closed on 2026-10-14.

The cause:

```text
A story's agent that runs flai stream open is left with the story in ready, so permission_prompt refuses its .claude/ write until it moves the story itself
```

The fix:

```text
S-0304 fixed it. flai stream open on a ready story now starts it with storystart.Begin, the code flai story start uses: it moves the story to in-progress as FLAI_AGENT, its epic following, or refuses it with exit 4 when the board holds it (ADR-0153). An agent that opens a ready story's stream is never left with it in ready, so permission_prompt asks for its protected write as in any in-progress story. TestStreamOpenStartsAReadyStory reproduces the cause.
```

Source: I-0097.
