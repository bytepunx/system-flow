---
id: T-1755
type: task
nature: research
title: Evaluate prose linters against this repository's agent prose and recommend one
status: done
parent: S-0413
owner: alex
created: 2026-10-16T05:46:12Z
updated: 2026-10-16T06:02:53Z
transitions:
  - to: ready
    at: 2026-10-16T05:46:55Z
    by: agent-S-0413
  - to: in-progress
    at: 2026-10-16T06:01:44Z
    by: agent-S-0413
  - to: done
    at: 2026-10-16T06:02:53Z
    by: agent-S-0413
stream: S-0413
tags: []
touches: [design/system/agent-prose.md]
after: [T-1750]
usage:
  source: log
  seconds: 69
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 131
      output: 54410
      cache_read: 5882736
      cache_write: 206541
      cost: 3.5418
---
# T-1755 Evaluate prose linters against this repository's agent prose and recommend one

## Work

Evaluate the natural language linters for this use: Vale, proselint, write-good, textlint, coala, and any other the research finds (LanguageTool, alex, Harper, readability scorers). For each: what it checks, whether its rules are configurable, how it runs (runtime, install, speed), its licence and upkeep, markdown support, and its fit for scoring agent prose from a Go CLI on the host. Run the strongest candidates from a scratch install against real prose of this repository (narratives, thread entries, issues) and record what they found. End with a recommendation of one or more and the rule set to start with.

Waits for T-1750, which creates the document this section goes in.

## Done when

- `design/system/agent-prose.md` has a `## Prose linters` section with a comparison table, the trial's results, and a recommendation with its reasons.
- Acceptance criterion 1 is met.

## Notes
