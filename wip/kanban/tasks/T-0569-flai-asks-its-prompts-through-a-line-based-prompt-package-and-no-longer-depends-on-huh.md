---
id: T-0569
type: task
nature: feature
title: flai asks its prompts through a line-based prompt package and no longer depends on huh
status: done
parent: S-0160
owner: alex
created: 2026-09-29T19:34:26Z
updated: 2026-09-29T19:39:17Z
transitions:
  - to: ready
    at: 2026-09-29T19:34:41Z
    by: agent-S-0160
  - to: in-progress
    at: 2026-09-29T19:34:42Z
    by: agent-S-0160
  - to: done
    at: 2026-09-29T19:39:17Z
    by: agent-S-0160
stream: S-0160
tags: []
touches: [flai/cmd, flai/internal/prompt, flai/go.mod, flai/go.sum, design/tech/go-libraries.md]
usage:
  source: log
  seconds: 275
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 52
      output: 338
      cache_read: 2375523
      cache_write: 38807
      cost: 0.9598
---

# T-0569 flai asks its prompts through a line-based prompt package and no longer depends on huh

## Work

Add `flai/internal/prompt`: `Input` (a default, a validation), `Confirm` (a default), and `Select` (numbered options, a default), which write a question to an `io.Writer`, read a line from an `io.Reader`, and ask again on an answer they cannot use. Replace every `huh` form in `cmd` (`new`, `import`, `upgrade`, `move ... cancelled`) with it, keeping each question's default and validation. Drop `huh` and what only it brought from `go.mod`, and say so in `design/tech/go-libraries.md` with the alternatives weighed.

## Done when

- `go mod why github.com/atotto/clipboard` says flai does not need it, and `go.mod` names no `charmbracelet` module.
- Behaviour tests cover each prompt: default on an empty line, a retry on a bad answer, end of input as an error.
- `make flai-test` passes.

## Notes
