---
id: T-0930
type: task
nature: feature
title: The analyzer's prompt and definition file each actionable finding as an issue and its report links each issue, and flai guard lets it
status: done
parent: S-0224
owner: alex
created: 2026-10-05T05:45:11Z
updated: 2026-10-06T21:25:30Z
transitions:
  - to: ready
    at: 2026-10-06T21:14:04Z
    by: agent-S-0224
  - to: in-progress
    at: 2026-10-06T21:14:04Z
    by: agent-S-0224
  - to: done
    at: 2026-10-06T21:25:30Z
    by: agent-S-0224
stream: S-0224
tags: [flai, template]
touches: [flai/internal/harness, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, ".claude/agents/analyzer.md", template/root/.claude/agents/analyzer.md, template/CHANGELOG.md, template/template.yaml, flai/cmd/guard.go, flai/cmd/guard_test.go, docs/users/flai-reference.md]
after: [T-0918, T-0925]
usage:
  source: log
  seconds: 686
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 132
      output: 60629
      cache_read: 5981554
      cache_write: 240721
      cost: 3.7734
---
# T-0930 The analyzer's prompt and definition file each actionable finding as an issue and its report links each issue, and flai guard lets it

## Work

Waits for T-0918 and T-0925: the prompt and the guard name their flags and tools. S-0223 made the analyzer's prompt, its definition, and its `analyze` guard rules, which this task extends.

- The analyzer's prompt in `flai/internal/harness` and its definition in `.claude/agents/analyzer.md` and `template/root/.claude/agents/analyzer.md` say: for each actionable finding, list the open issues (`flai issue list --json`), bump the one that records it with `--report` and its impact, else file one with `flai issue new` (or `issue_new`) giving the class from the finding (`defect`; `efficiency`; `impression` for a risk with no measured instance), the impact as time lost per cycle or revenue or penalty, the evidence, and `--report`; then link each issue it filed or bumped from the finding in its report. It never makes a story.
- `flai guard` with the `analyze` role passes `flai issue new` and `bump` and the MCP `issue_new` and `issue_bump`, and refuses `issue story`, `issue close`, and `item_new`, since the analyzer never authors stories.
- An agent's write under this repository's `.claude/` is refused: send the definition's whole contents to the operator on the story's thread to paste, then check what was pasted.
- Add a `template/CHANGELOG.md` entry for the template's definition.

## Done when

- The prompt and both definitions say how a finding becomes an issue, how a duplicate is bumped, and that the report links each issue
- Guard tests cover the analyzer's issue calls passing and its story calls refused; `scripts/flai-test.sh` passes

## Notes
