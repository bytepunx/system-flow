---
id: T-2166
type: task
nature: research
title: One story of this repository is worked inside the container over HTTP, and the tools that answered, the guard's run, and the close-out tiers that could not run are recorded with issues
status: done
parent: S-0505
owner: alex
created: 2026-10-26T04:50:07Z
updated: 2026-10-28T12:31:52Z
transitions:
  - to: ready
    at: 2026-10-28T12:22:33Z
    by: agent-S-0505
  - to: in-progress
    at: 2026-10-28T12:22:34Z
    by: agent-S-0505
  - to: done
    at: 2026-10-28T12:31:52Z
    by: agent-S-0505
stream: S-0505
tags: [cli, docker]
touches: [design/experiments/S-0505-claude-code-works-one-story-of-this-repository-inside-a-hand-built-container-from-a-bedrock-image-and-a-claude-code-layer-and-the-results-settle-what-adr-0203-leaves-open.md, design/issues/summary.md, flai/images/claude-code/Dockerfile]
after: [T-2159]
usage:
  source: log
  seconds: 558
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 105
      output: 34348
      cache_read: 6490648
      cache_write: 145995
      cost: 2.9996
---
# T-2166 One story of this repository is worked inside the container over HTTP, and the tools that answered, the guard's run, and the close-out tiers that could not run are recorded with issues

## Work

With the operator, choose one small story of this repository that changes documentation only (or take the one they name), start it on the host with `flai story start` so its worktree and narrative exist, then run `scripts/agent-container.sh` on it with the story's prompt as `flai serve` would give it, and let Claude Code inside the container work it to review: tasks written, each closed with `task_done` over HTTP, `stream_sync`, `inbox`, `wait_for_events`, the narrative's `stream_state`, and the close-out. Record in the results document's "Results": each MCP tool called, whether it answered as the stdio server does, and each refusal or missing tool; whether `flai guard` ran as the hook on the agent's calls inside the container (the hook's log lines name it) and what it refused; whether `permission_prompt` reached the operator on a protected write, tried on purpose once; and each tier of `scripts/close-out.sh` that could not run for want of a tool in the image, with the tool named. For each limitation that needs a change in flai, run `flai issue new` in the story's worktree with the finding and the story to carry it on, and link the issue from the results. Waits for T-2159, which gives the container and the MCP route.

## Done when

The chosen story is in review, the results document lists every tool called with its outcome, the guard's and the permission prompt's behaviour, and the tiers that could not run, and each limitation needing code has an issue linked from the results.

## Notes

If the story cannot reach review inside the container, say so in the results with the step it stopped at; a stopped spike is a result too.
