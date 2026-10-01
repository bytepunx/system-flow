---
id: S-0182
type: story
nature: remediation
title: Find and fix how flai serve 1.23.0 started an agent for a story the board showed held
status: backlog
owner: alex
created: 2026-10-01T08:00:32Z
updated: 2026-10-01T08:00:32Z
transitions: []
tags: [flai]
touches: [flai/internal/serve, flai/internal/workitem]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
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

## Notes

If the evidence is gone and the cause cannot be found, the story ends with the regression tests for the suspect paths and a note in I-0050 instead of closing it.
