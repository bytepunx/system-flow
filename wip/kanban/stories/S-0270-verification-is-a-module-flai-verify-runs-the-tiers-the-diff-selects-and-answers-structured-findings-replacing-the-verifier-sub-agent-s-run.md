---
id: S-0270
type: story
nature: improvement
title: "Verification is a module: flai verify runs the tiers the diff selects and answers structured findings, replacing the verifier sub-agent's run"
status: backlog
parent: E-0017
owner: alex
created: 2026-10-05T01:35:29Z
updated: 2026-10-05T01:35:52Z
transitions: []
tags: []
after: [S-0273]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0270 Verification is a module: flai verify runs the tiers the diff selects and answers structured findings, replacing the verifier sub-agent's run

## Goal

Before review a story agent starts a Sonnet verifier sub-agent to run the close-out and read its log: 88 minutes over 22 stories, and in S-0248 three runs of the same two-minute script. `flai verify S-nnnn`, `verify` over MCP, and `verify.run` on the host channel run what `scripts/close-out.sh` runs for the branch's diff, cheapest tier first, and answer findings as data: step, path, line, message, and whether the step passed, failed, or was not reached, with the duration of each. A finding outside the story (S-0249) is a note, not a failure. The agent reads twenty lines, not a log, and the sub-agent is kept for judgement, checking the diff against the criteria and the conventions, if at all.

## Acceptance criteria
- [ ] `flai verify S-nnnn` runs the project's lint, test tiers, check, and narrative check for the story's worktree and answers one structured result, as text and `--json`, with each step's state and duration
- [ ] The same is `verify` over MCP and `verify.run` on the host channel; the dashboard's story page can show the last result
- [ ] `scripts/close-out.sh` calls it, or shares its implementation, so that the two never disagree
- [ ] `design/conventions/delegation.md`, the harness prompt, and the verifier definition send the agent to `flai verify` for the run and keep the sub-agent for the criteria and conventions review
- [ ] `design/system/flai-cli.md`, `devex.md`, and the user guide describe it

## Tasks

## Notes

Companion to S-0266 (verifier exit status) and S-0267 (duplicate tier), which stay worth doing until this lands.
