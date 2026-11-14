---
id: S-0786
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-14T17:32:28Z
updated: 2026-11-14T17:36:11Z
transitions:
  - to: ready
    at: 2026-11-14T17:32:29Z
    by: system-flow
  - to: in-progress
    at: 2026-11-14T17:32:37Z
    by: agent-S-0786
  - to: review
    at: 2026-11-14T17:34:39Z
    by: agent-S-0786
  - to: done
    at: 2026-11-14T17:36:11Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0786-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 130
  priced_by: harness
  turns:
    - day: 2026-11-14
      ceremony: 2
      work: 26
  models:
    - model: claude-opus-5-5
      input: 62
      output: 7171
      cache_read: 4119618
      cache_write: 164332
      cost: 2.2822
  strategic:
    - kind: orchestrator
      seconds: 84
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 23
          output: 200
          cache_read: 853348
          cache_write: 121687
          cost: 0.2893
        - model: claude-sonnet-5-5
          input: 4
          output: 40
          cache_read: 13141
          cache_write: 33514
          cost: 0.0364
formula: publish
---
# S-0786 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3358 Fetch the remote's branch and tags
- T-3359 Evaluate the release policy, stopping unless it is met
- T-3360 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

### Accepted by the orchestrator

- Verified: 13c53d75acee2072c7f59e39b97041a8e687f93a
- At: 2026-11-14T17:36:11Z

Verdict: S-0786 meets all three criteria; the publish step released flai/v1.95.5 and flaiover/v0.61.0, verified at 13c53d75acee2072c7f59e39b97041a8e687f93a.
- 1: T-3358 done
- 2: T-3359 done
- 3: T-3360 done
