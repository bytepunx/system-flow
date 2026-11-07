---
id: T-2498
type: task
nature: feature
title: The live planner run's preconditions are in place, and the operator has set the planner on harness container and asked for a plan of one backlog story
status: in-progress
parent: S-0567
owner: alex
created: 2026-11-07T03:43:08Z
updated: 2026-11-07T03:45:06Z
transitions:
  - to: ready
    at: 2026-11-07T03:45:05Z
    by: agent-S-0567
  - to: in-progress
    at: 2026-11-07T03:45:06Z
    by: agent-S-0567
stream: S-0567
tags: [docker]
usage:
  source: log
  seconds: 378
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 81
      output: 16922
      cache_read: 7087831
      cache_write: 198825
      cost: 3.3469
---
# T-2498 The live planner run's preconditions are in place, and the operator has set the planner on harness container and asked for a plan of one backlog story

## Work

Check what the run needs on this host, then ask the operator on one thread for the parts that are theirs. This task waits for nothing.

- The installed flai holds S-0509's start paths: `flai version`, and its release notes or `flai serve status`, name a release after S-0509's acceptance on 2026-11-02.
- The image is built: `docker image inspect flai-agent-sf-claude-code:local`, the project's layer the manifest's `agent.image.layer` names, over `flai-agent-claude-code:local`. When it is missing, `make agent-images` builds it.
- Choose the backlog story the planner plans: the first story `flai plan --candidates` lists with no plan, or one the operator names.
- Open one thread on S-0567 asking the operator, each command written out: the host's container harness login (`flai serve agent harness container --login <path>`), which the guard refuses an agent (ADR-0220); `planning.agent.harness: container` in `system-flow.yaml`, through `flai manifest set` or the dashboard's settings; and the plan of the chosen story from the dashboard. Say whether `planning.agent` goes back to the host harness after the run.

Record what was checked and what the operator set in the story's Notes under `### Live run`, with `flai edit S-0567`.

## Done when

- The installed flai's version, the image's tag, and the chosen story are recorded in S-0567's Notes.
- The operator has answered the thread: the host's container harness and the planner's harness are set, and a plan of the chosen story is asked for.

## Notes
