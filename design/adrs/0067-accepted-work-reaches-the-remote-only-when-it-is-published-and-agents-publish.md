---
id: ADR-0067
title: "Accepted work reaches the remote only when it is published, and agents publish only when the operator asks"
status: proposed
date: 2026-10-02
supersedes: []
superseded_by: []
refines: [ADR-0032, ADR-0048]
---

# ADR-0067 Accepted work reaches the remote only when it is published, and agents publish only when the operator asks

## Context

Since ADR-0032 accepting a story merges and commits and nothing else, and publishing (`flai release --pending`, or the board's Publish) tags and pushes. ADR-0048 kept a separate push of accepted work, `flai push --pending`, which releases nothing unless the `auto-publish` host action is on, and S-0063 made an unpushed acceptance an agent's duty: when `inbox` or `flai board` reports `unpushed`, the agent runs `git fetch` and `flai push --pending` before anything else. Two ways to send accepted work to the remote, one of them an agent's standing duty, make it unclear when work reaches the remote and who decides. The designer decided on 2026-10-02 that publishing is the one way accepted work reaches the remote.

## Decision

Accepted work reaches the remote only when it is published. Publishing is `git fetch`, then `flai release --pending` (or the board's Publish, which runs the same), which tags what was accepted since each component's last release and pushes the branch and the tags, never forcing; when the remote has moved it refuses, and the remote branch is merged, never rebased, before it runs again. Agents publish only when the operator asks them to. An unpushed acceptance is no longer an agent's duty: `inbox`, `flai board`, and the MCP server's instructions stop telling agents to push it.

## Consequences

- The `unpushed` duty in the work-management convention, the MCP server's instructions, and the harness prompt is removed; the board and the dashboard still show what is accepted and not yet published.
- `flai push --pending` and the `auto-publish` host action (ADR-0048) are no longer part of the workflow. Whether they are removed, or kept for an operator's own use, is decided by the story that builds this.
- Work accepted and not yet published stays local until the operator publishes; a clone that publishes must have the remote's tags (S-0174).
- `design/system/pushing-from-the-board.md`, `flai-cli.md`, `flaiover-dashboard.md`, and the operator and user guides change with the code.

## Alternatives considered

- Keep the agent's duty to push unpushed acceptances: two paths to the remote, and agents pushing on their own initiative.
- Have acceptance push again: ADR-0032 separated the two so that publishing is deliberate and batched.
