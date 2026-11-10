---
id: T-2944
type: task
nature: remediation
title: SetStreamState replaces Context and appends Decisions and Open questions bullets, keeping the threads block
status: done
parent: S-0668
owner: alex
created: 2026-11-10T08:50:37Z
updated: 2026-11-10T09:49:33Z
transitions:
  - to: ready
    at: 2026-11-10T09:36:42Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T09:36:42Z
    by: system-flow
  - to: done
    at: 2026-11-10T09:49:33Z
    by: agent-S-0668
stream: S-0668
tags: [flai, narrative]
touches: [flai/internal/workitem/narrative.go, flai/internal/workitem/narrative_test.go]
usage:
  source: log
  seconds: 770
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 44
      output: 285
      cache_read: 5226273
      cache_write: 52417
      cost: 1.9431
---
# T-2944 SetStreamState replaces Context and appends Decisions and Open questions bullets, keeping the threads block

## Work

In `flai/internal/workitem/narrative.go`, widen what `Repo.SetStreamState` writes, as the story's `### Planning` proposes:

- Take the texts as one value, such as a `StreamText` with `Current`, `Next`, `Context`, `Decisions []string`, and `Questions []string`, in place of the two strings. Keep a wrapper with the old two-string signature, or change every caller in this task, so that `flai/cmd`, `flai/internal/mcpserver`, and `flai/internal/itemedit` still build when this task closes.
- `Context` replaces what is under `## Context`, as `replaceSection` replaces Current state.
- Each entry of `Decisions` is appended as one bullet under `## Decisions`, in the shape `flai/internal/digest/sources.go` reads. The section is made where the template puts it when missing.
- Each entry of `Questions` is appended as a bullet that starts with today's date, under `## Open questions` and above the `<!-- threads:start -->` block. The block between the markers is kept byte for byte, as `OpenQuestions` and `AnswerOpenQuestion` skip it.
- Each text keeps today's refusals: a `#` or `##` heading, text the lint rejects, a story not in progress or in review, and a story with no narrative. Giving no text at all is still an error.
- `StreamState.Written` lists the sections written in template order. The same `Context` again writes nothing. A `Decisions` or `Questions` entry always appends.

Add tests in `flai/internal/workitem/narrative_test.go` beside `TestSetStreamStateReplacesBothSections`:

- Each new text is written, and the other sections stay as they were.
- A question lands above a threads block that holds a thread, and the block is unchanged.
- `OpenQuestions` lists the new question, and `AnswerOpenQuestion` moves it to Decisions.

Waits for no task: the CLI, the MCP tool, and the host channel all call it.

## Done when

- `flai test flai/internal/workitem/narrative.go flai/internal/workitem/narrative_test.go` passes.
- The threads-block test fails when the write replaces the whole `## Open questions` section.
- `go build ./...` in `flai/` passes.

## Notes
