---
id: T-3357
type: task
nature: feature
title: A log path in a thread entry or a message that flai wrote renders as a link to the log page
status: done
parent: S-0785
owner: alex
created: 2026-11-14T10:30:44Z
updated: 2026-11-14T11:16:13Z
transitions:
  - to: ready
    at: 2026-11-14T10:56:02Z
    by: agent-S-0785
  - to: in-progress
    at: 2026-11-14T10:56:02Z
    by: agent-S-0785
  - to: done
    at: 2026-11-14T11:16:13Z
    by: agent-S-0785
stream: S-0785
tags: [dashboard]
touches: [flaiover/src/lib/markdown.ts, flaiover/src/lib/markdown.test.ts, flaiover/src/lib/filelinks.test.ts]
after: [T-3355]
usage:
  source: log
  seconds: 1211
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 36
      output: 12233
      cache_read: 2396564
      cache_write: 59017
      cost: 1.0732
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 6
          cache_read: 20791
          cache_write: 579
          cost: 0.0063
---
# T-3357 A log path in a thread entry or a message that flai wrote renders as a link to the log page

## Work

flai serve writes most of its failures to threads and messages, not to a page's own fields. For example, a deploy, teardown, promotion, or checks failure in `flai/internal/serve/kluster.go` names its log in backticks through `said()`, and so do `StepError.Error` in `flai/internal/kluster/deploy.go` and the image job and container failures in `flai/internal/serve`. The dashboard renders those entries through `$lib/markdown.ts`. Its `$lib/filelinks.ts` links only `.md` paths under `design`, `docs`, and `wip`, so a log path there stays plain text, in a code span or not.

- Recognise a log path in rendered markdown text and in inline code spans with the same function `$lib/logs.ts` gives `LogText` (T-3355). Render it as a link to `/ops/log?path=<encoded path>`, keeping the project the page is on, as `LogLink` does. Leave a fenced code block alone.
- Keep every link `filelinks.ts` makes today as it is.
- Tests in `markdown.test.ts` and `filelinks.test.ts`: an absolute, a `~/`, and a `.flai-cache/` log path, in plain text and in backticks, become links; a `.log` inside a fenced block, a path not ending in `.log`, and a URL do not.
- Waits for T-3355, whose `$lib/logs.ts` function and log page it uses. It runs beside T-3356, with no path in common.

## Done when

- A thread entry that names a log, such as a kluster deploy failure's, shows that log as a link to the log page on the item page, the inbox, and the threads page.
- `flai test` passes on the paths changed.

## Notes
