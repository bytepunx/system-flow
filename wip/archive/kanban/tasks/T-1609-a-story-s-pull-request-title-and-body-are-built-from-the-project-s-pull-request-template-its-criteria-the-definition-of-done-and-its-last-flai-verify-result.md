---
id: T-1609
type: task
nature: feature
title: A story's pull request title and body are built from the project's pull request template, its criteria, the definition of done, and its last flai verify result
status: done
parent: S-0393
owner: alex
created: 2026-10-16T04:35:21Z
updated: 2026-10-17T07:57:33Z
transitions:
  - to: ready
    at: 2026-10-17T07:44:47Z
    by: agent-S-0393
  - to: in-progress
    at: 2026-10-17T07:44:50Z
    by: agent-S-0393
  - to: done
    at: 2026-10-17T07:57:33Z
    by: agent-S-0393
stream: S-0393
tags: [cli, template]
touches: [flai/internal/forge/pullrequest.go, flai/internal/forge/pullrequest_test.go, template/root/.github/pull_request_template.md]
usage:
  source: log
  seconds: 763
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 26
      output: 10872
      cache_read: 1399798
      cache_write: 41241
      cost: 0.7328
---
# T-1609 A story's pull request title and body are built from the project's pull request template, its criteria, the definition of done, and its last flai verify result

## Work

- In `flai/internal/forge/pullrequest.go`, build the title `[S-nnnn] <title>` and the body. The body fills `.github/pull_request_template.md` with four parts:
  - the story's ID
  - its acceptance criteria, as checkboxes ticked as the story has them
  - the definition of done
  - the result of the story's last `flai verify`, its outcome and commit, as `flai verify --last` reads it
- Without the file, write the same four parts under headings of their own.
- Mark the places the template fills, so that a project's own wording around them stays. Add those marks to `template/root/.github/pull_request_template.md`.
- Test it with the template, without the template, and with no verify result.
- First in its layer: it waits for no task.

## Done when

The title and body are built in all three cases, the template file carries the places, and the tests pass.

## Notes
