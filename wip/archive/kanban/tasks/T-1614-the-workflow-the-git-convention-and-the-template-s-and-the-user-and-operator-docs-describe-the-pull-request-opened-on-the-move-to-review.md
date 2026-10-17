---
id: T-1614
type: task
nature: feature
title: The workflow, the git convention and the template's, and the user and operator docs describe the pull request opened on the move to review
status: done
parent: S-0393
owner: alex
created: 2026-10-16T04:35:50Z
updated: 2026-10-17T08:45:16Z
transitions:
  - to: ready
    at: 2026-10-17T08:33:47Z
    by: agent-S-0393
  - to: in-progress
    at: 2026-10-17T08:33:50Z
    by: agent-S-0393
  - to: done
    at: 2026-10-17T08:45:15Z
    by: agent-S-0393
stream: S-0393
tags: [cli, template]
touches: [design/system/workflow.md, design/system/flai-cli.md, design/conventions/git.md, template/root/design/conventions/git.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, docs/operators/settings.md, template/CHANGELOG.md, design/system/devex.md]
after: [T-1611, T-1613]
usage:
  source: log
  seconds: 685
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 131
      output: 54963
      cache_read: 7076363
      cache_write: 208487
      cost: 3.7044
---
# T-1614 The workflow, the git convention and the template's, and the user and operator docs describe the pull request opened on the move to review

## Work

- `design/system/workflow.md`, § Transitions and who makes them: under pull-request acceptance, in-progress to review also has flai serve push the branch and open or update the pull request.
- `design/conventions/git.md` and `template/root/design/conventions/git.md`:
  - Baseline: under pull-request acceptance, flai opens the pull request from `.github/pull_request_template.md`, and an agent never opens one itself.
  - Baseline: the story-branch exception to never force pushing, as TH-0463 settles it. Change the template's baseline first, then this repository's, in this story.
- `design/system/flai-cli.md`, `docs/users/flai.md`, and the regenerated reference: `flai forge push` and `flai forge status`.
- `docs/operators/index.md` and `settings.md`: the `forge` host action, and what turning it on lets flai serve do.
- Add a `template/CHANGELOG.md` entry.
- Waits for T-1611 (the ADR the docs link) and T-1613 (the commands' final flags).

## Done when

Every document above describes the pull request on the move to review, the settings and reference tests pass, and the markdown lint is clean.

## Notes
