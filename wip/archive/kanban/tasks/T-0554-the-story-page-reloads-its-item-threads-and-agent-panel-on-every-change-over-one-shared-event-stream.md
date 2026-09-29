---
id: T-0554
type: task
nature: feature
title: The story page reloads its item, threads, and agent panel on every change, over one shared event stream
status: done
parent: S-0154
owner: alex
created: 2026-09-29T07:08:18Z
updated: 2026-09-29T07:18:24Z
transitions:
  - to: ready
    at: 2026-09-29T07:08:32Z
    by: agent-S-0154
  - to: in-progress
    at: 2026-09-29T07:12:38Z
    by: agent-S-0154
  - to: done
    at: 2026-09-29T07:18:24Z
    by: agent-S-0154
stream: S-0154
tags: []
usage:
  source: log
  seconds: 346
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 90
      output: 25225
      cache_read: 6436067
      cache_write: 84067
      cost: 2.4646
---

# T-0554 The story page reloads its item, threads, and agent panel on every change, over one shared event stream

## Work

The dashboard forwards `agent` as a server-sent event of its own on `/api/events`. One shared EventSource per project serves every listener on a page (the inbox, the story page, its threads, its agent panel) instead of one each. The story page reloads its item, the Threads pane its threads, and StoryAgent the agent's state when the project changes, debounced.

## Done when

Behaviour tests show the item page and Threads reload on a change event, StoryAgent re-asks on an `agent` event, and one EventSource serves several listeners.

## Notes
