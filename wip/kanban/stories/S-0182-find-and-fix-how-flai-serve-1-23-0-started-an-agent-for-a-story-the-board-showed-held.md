---
id: S-0182
type: story
nature: remediation
title: Find and fix how flai serve 1.23.0 started an agent for a story the board showed held
status: in-progress
owner: alex
created: 2026-10-01T08:00:32Z
updated: 2026-10-01T09:23:32Z
transitions:
  - to: ready
    at: 2026-10-01T08:31:51Z
    by: alex
  - to: in-progress
    at: 2026-10-01T09:11:13Z
    by: agent-S-0182
tags: [flai]
touches: [flai/internal/serve, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0050-flai-serve-1-23-0-which-has-holds-started-an-agent-for-a-story-the-board-showed-held.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1079
  models:
    - model: claude-opus-5-5
      input: 278
      output: 80351
      cache_read: 23063664
      cache_write: 284355
      cost: 8.3905
---
# S-0182 Find and fix how flai serve 1.23.0 started an agent for a story the board showed held

## Goal

I-0050 (two occurrences, last 2026-09-29): `flai serve` 1.23.0, which already had holds (S-0128, ADR-0046), started an agent for a story the board showed as held. The launcher skips any story `holds.Of` returns a hold for (`flai/internal/serve/agents.go` ~417), and the board and `wait_for_work` use the same `Repo.Holds`, so the cause is not known. Two paths bypass holds and are suspects, though neither obviously fits a fresh start: `resume()` (`agents.go` ~524-538), which restarts an agent whose question was answered, and the operator's start (`serve/start.go` ~19, ~57-68), which ignores holds by design with a warning.

## Acceptance criteria
- [ ] The cause is found from the 2026-09-29 instances (the issue's instances, `~/.flai/serve/journal.jsonl`, `serve.log`, and `agents.json` where they survive, and the stories' transitions) and recorded in the story's narrative
- [ ] A regression test reproduces it, including an in-progress story whose claim is a file and a ready story whose claim is that file's directory, and passes after the fix
- [ ] No launcher path starts a held story except the operator's explicit start, and `resume()` either respects a hold or the design says why it need not
- [ ] I-0050 is closed with the cause and what fixed it

## Tasks
- T-0652 Find the cause of I-0050 from the 2026-09-29 evidence and record it in the narrative
- T-0653 resume starts again only an agent whose story is still open, so a story sent back to ready waits for its hold and the limit
- T-0654 An agent the operator starts past a hold or a full limit is told so in its prompt
- T-0655 Regression tests for every launcher start path with file and directory claims, the design and guide updated, and I-0050 closed

## Notes

If the evidence is gone and the cause cannot be found, the story ends with the regression tests for the suspect paths and a note in I-0050 instead of closing it.
