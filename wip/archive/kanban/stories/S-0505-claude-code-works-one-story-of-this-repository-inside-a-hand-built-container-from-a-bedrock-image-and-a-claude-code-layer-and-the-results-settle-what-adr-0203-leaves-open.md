---
id: S-0505
type: story
nature: experiment
title: Claude Code works one story of this repository inside a hand-built container from a bedrock image and a claude-code layer, and the results settle what ADR-0203 leaves open
status: done
parent: E-0029
owner: alex
created: 2026-10-26T04:42:20Z
updated: 2026-10-28T12:52:31Z
transitions:
  - to: ready
    at: 2026-10-28T12:04:36Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-28T12:04:42Z
    by: agent-S-0505
  - to: review
    at: 2026-10-28T12:35:04Z
    by: system-flow
  - to: done
    at: 2026-10-28T12:52:31Z
    by: alex
tags: [cli, docker]
topics: [harness, docker, experiments]
touches: [flai/images/bedrock/Dockerfile, flai/images/claude-code/Dockerfile, scripts/agent-image.sh, scripts/agent-container.sh, Makefile, scripts/README.md, ".dockerignore", design/experiments/S-0505-claude-code-works-one-story-of-this-repository-inside-a-hand-built-container-from-a-bedrock-image-and-a-claude-code-layer-and-the-results-settle-what-adr-0203-leaves-open.md, design/experiments/README.md, design/system/agent-adapters.md, design/tech/docker.md, design/issues/summary.md, design/issues]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2014
  priced_by: harness
  turns:
    - day: 2026-10-28
      ceremony: 4
      hand_edits: 2
      work: 85
  models:
    - model: claude-opus-5-5
      input: 346
      output: 109490
      cache_read: 20250852
      cache_write: 481642
      cost: 9.6306
  strategic:
    - kind: orchestrator
      seconds: 630
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 88
          output: 1170
          cache_read: 16753733
          cache_write: 37433
          cost: 3.684
        - model: claude-sonnet-5-5
          input: 10
          output: 59
          cache_read: 127995
          cache_write: 51021
          cost: 0.1224
cost_of_delay:
  value: 33.33
  by: planner-S-0505
  at: 2026-10-26T12:27:51Z
forecast:
  duration: 3h
  delivery: 2026-10-28T02:53:00Z
  basis: "Its own forecast of 3h; 23rd in the pull order with an in-progress limit of 5, behind S-0503, S-0504, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425, S-0499 and S-0502."
  by: flai
  at: 2026-10-26T16:03:52Z
finalized:
  by: orchestrator
  at: 2026-10-26T12:28:25Z
---
# S-0505 Claude Code works one story of this repository inside a hand-built container from a bedrock image and a claude-code layer, and the results settle what ADR-0203 leaves open

## Goal

The spike ADR-0203 names as its first story, built the way E-0029 refines it: two image layers built locally and by hand, no adapter, no pilot. The bedrock layer is a small Linux base, `git`, and the `flai` binary built from this tree; the claude-code layer adds Node and Claude Code on top of it. A container from the claude-code image is started by hand with the story's worktree mounted read-write at its host path, a folder of the story's own under `.flai-cache` as a writable `CLAUDE_CONFIG_DIR`, the operator's Claude login mounted read-only, and nothing else, and Claude Code inside it works one story of this repository through the HTTP MCP server `flai host` keeps for the project, with the bearer token and `X-Flai-Agent`, closing each task with `task_done`.

The results document settles, each by measurement, what [ADR-0203](../../../design/adrs/0203-flai-may-run-a-story-s-agent-in-a-container-from-a-baseline-image-it-builds.md) and [agent-adapters.md § A container harness of flai's own](../../../design/system/agent-adapters.md#a-container-harness-of-flais-own) leave open: whether `claude` signs in and runs from a login mounted read-only with a writable per-story configuration folder; how fast the bind mount is on this WSL2 host; which tools of the stdio MCP server the HTTP server answers for an agent named by header, `task_done`, `stream_sync`, `permission_prompt`, and `wait_for_events` among them; whether `flai guard` runs as the hook inside the container; by which network mode the container reaches the host's listener; and which test tiers of the close-out cannot run for want of a tool in the image. Each limitation that needs a code change becomes an issue, so that the stories after this one, the HTTP MCP server for containers, pilot, and the container adapter, start from facts.

## Acceptance criteria

- [x] `scripts/agent-image.sh` builds, from this tree and locally, `flai-agent-bedrock:local` from `flai/images/bedrock/Dockerfile` (a small Linux base, `git`, and the `flai` binary built from the tree, as `flaiover/Dockerfile` builds it) and `flai-agent-claude-code:local` from `flai/images/claude-code/Dockerfile` on top of it (Node and Claude Code), `make agent-images` runs it, and no GitHub workflow builds or publishes either image
- [x] `scripts/agent-container.sh S-nnnn` starts a container from the claude-code image by hand, with the story's worktree mounted read-write at its host path, `.flai-cache/agents/S-nnnn` mounted as the writable `CLAUDE_CONFIG_DIR`, the operator's Claude login mounted read-only, `--user` the operator's uid and gid, no other mount and never the Docker socket, and Claude Code inside it reaching the project's HTTP MCP server on the host with `--mcp-config` naming a file in the mounted cache folder, the token from `flai mcp token`, and `X-Flai-Agent`
- [x] One story of this repository, chosen with the operator, is worked inside that container from its first task to its move to review, each task closed with `task_done` over HTTP, and the results record each tool that answered as the stdio server does, each that did not, and whether `flai guard` ran as the hook on the agent's calls
- [x] The results document under `design/experiments/`, written from its template, records each measurement: the login (what `claude` read and wrote, and where), the bind mount's speed as the time of `git status`, `flai check --strict`, and a file walk of the worktree inside the container against the same commands on the host, the network mode by which the container reached the host's listener and the modes that failed, and each test tier of `scripts/close-out.sh` that could not run inside the container for want of a tool, with the tool named
- [x] Each limitation found that needs a change in flai is recorded as an issue with `flai issue new`, linked from the results, and the results' recommendation names adopt, adapt, or drop for the container harness and the story of E-0029 that carries each limitation on
- [x] `design/system/agent-adapters.md` § A container harness of flai's own records in one paragraph what the spike settled, linking the results document, `design/tech/docker.md` lists the two images, and `scripts/README.md` names the two scripts

## Tasks
- T-2149 The bedrock and claude-code images build locally from Dockerfiles under flai/images with scripts/agent-image.sh and make agent-images
- T-2159 scripts/agent-container.sh starts a container by hand with the worktree, the cache folder, and the login mounted, and Claude Code inside it signs in and reaches the project's HTTP MCP server
- T-2166 One story of this repository is worked inside the container over HTTP, and the tools that answered, the guard's run, and the close-out tiers that could not run are recorded with issues
- T-2173 The bind mount's speed is measured inside the container against the host with git status, flai check --strict, and a file walk of the worktree
- T-2177 The results document is finished with its recommendation and indexed, and agent-adapters.md and docker.md record what the spike settled

## Notes

ADR-0203 § Decision, First story, names this spike; E-0029's notes refine it to layered images built locally. The images carry no release stamp and are not signed: the signed, published image is a later story, outside this epic. The story works a story of this repository that changes documentation only unless the operator names another, since the bedrock and claude-code layers carry no Go toolchain; what the Go tiers needed is a finding for the project layer story, not a reason to add tools here.

### Planning

Planned by planner-E-0029 on 2026-10-26 with the epic. Revisited by planner-S-0505 on 2026-10-26 against the story's criteria: the five tasks, their touches, and their `after` stand as drafted, the forecast is kept, and the cost of delay value is written now that E-0029 has its input.

Touches. `flai touches suggest S-0505` on the declared touches (696 of 2427 commits changed them) lists `design/system/flai-cli.md` (8%), `docs/users/flai.md` (7%), `docs/operators/index.md` (5%), and the dashboard's design among the co-changes; none is added, since the story changes no command, setting, or dashboard behaviour. Every touch is declared and a task's (layout), with these sources:

- `flai/images/bedrock/Dockerfile`, `flai/images/claude-code/Dockerfile`, `scripts/agent-image.sh`, `scripts/agent-container.sh`, `Makefile`, `scripts/README.md`, `.dockerignore` (layout): the images live under the flai sub-project because they ship flai, as `flaiover/Dockerfile` lives under flaiover; `.dockerignore` at the root is the build context's, as for flaiover. T-2149 and T-2159.
- `design/experiments/S-0505-….md`, `design/experiments/README.md` (design): the experiment results document and its index, as `design/experiments/README.md` requires. T-2159, T-2166, T-2173, T-2177.
- `design/system/agent-adapters.md` (design, criterion 6): the section the ADR cites gains the spike's paragraph. T-2177.
- `design/tech/docker.md` (co-change, design): changed with the flaiover image in 3 of 15 seed commits, and criterion 6 names it. T-2177.
- `design/issues/summary.md` and the folder `design/issues` (layout, folder kept): `flai issue new` in T-2166 writes issue files whose numbers and slugs no task can name; the folder is kept for them alone, and `flai task done` widens the touches with the committed paths (ADR-0144).

Forecast. `flai forecast S-0505` gives 1h12m (median 227 s per unit over 3 done experiment stories; size 19: 6 criteria, 13 touches) and delivery 2026-10-27T11:14Z, 23rd in the pull order. Kept at 3h, the epic planner's figure, for its reason: the spike builds two images (a Go build and a Claude Code install, minutes each, built several times), runs a whole story of this repository inside the container, a session of its own of 30 to 60 minutes on this repository's recent stories, and measures the login, the network modes, and the mount by hand, where the three experiment stories in the history measured inside one process. Delivery 2026-10-27T23:16Z is flai's delivery for 1h12m played out for 3h as flai's replan of 2026-10-26T04:57Z did (1h12m to 3h moved it 12h02m, from 03:48Z to 15:50Z); flai's next replan plays it out again from the pull order.

Cost of delay. `flai cod S-0505` gives 33.33 USD a week: the story's share of E-0029's 150 USD a week of penalty, set by the operator on 2026-10-26 after TH-0657, 3h of the 13h30m forecast over the epic's six open stories without inputs. Written as it stands: the share by forecast duration is flai's rule, and nothing argues for inputs of the story's own, since the spike brings no revenue and loses no work per cycle apart from the epic's.
