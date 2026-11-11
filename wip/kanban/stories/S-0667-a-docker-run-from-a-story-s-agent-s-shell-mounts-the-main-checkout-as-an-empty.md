---
id: S-0667
type: story
nature: improvement
title: A docker run from a story's agent's shell mounts the main checkout as an empty folder, so a hand call in the agent image cannot read the project
status: review
owner: alex
created: 2026-11-10T08:34:40Z
updated: 2026-11-10T19:29:36Z
transitions:
  - to: ready
    at: 2026-11-10T08:58:35Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-10T10:58:37Z
    by: agent-S-0667
  - to: review
    at: 2026-11-10T19:29:36Z
    by: system-flow
blocked:
  - from: 2026-11-10T11:05:54Z
    until: 2026-11-10T19:28:15Z
    reason: awaiting the operator on TH-1047
tags: [docker, agent-container]
touches: [docs/operators/runbooks/agent-container.md, design/issues/I-0276-a-docker-run-from-a-story-s-agent-s-shell-mounts-the-main-checkout-as-an-empty.md, design/issues/summary.md, design/conventions/tooling.md, design/issues/I-0344-one-malformed-work-item-file-in-the-main-checkout-s-wip-makes-every-flai-read.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 519
  priced_by: harness
  turns:
    - day: 2026-11-10
      ceremony: 5
      hand_edits: 2
      work: 74
  models:
    - model: claude-opus-5-5
      input: 174
      output: 33258
      cache_read: 17100972
      cache_write: 505643
      cost: 8.1312
  strategic:
    - kind: orchestrator
      seconds: 161
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 34
          output: 240
          cache_read: 2514259
          cache_write: 250675
          cost: 0.7554
        - model: claude-sonnet-5-5
          input: 2
          output: 14
          cache_read: 6369
          cache_write: 19789
          cost: 0.0197
cost_of_delay:
  inputs:
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-11-10T08:34:40Z
  value: 50
  by: planner-S-0667
  at: 2026-11-10T08:49:29Z
forecast:
  duration: 1h30m
  delivery: 2026-11-10T13:11:00Z
  basis: "Its own forecast of 1h30m; 1st in the pull order with an in-progress limit of 3, behind S-0665, S-0708 and S-0713."
  by: flai
  at: 2026-11-10T10:55:00Z
finalized:
  by: alex
  at: 2026-11-10T08:51:21Z
---
# S-0667 A docker run from a story's agent's shell mounts the main checkout as an empty folder, so a hand call in the agent image cannot read the project

## Goal

This story remediates [I-0276](../../../design/issues/I-0276-a-docker-run-from-a-story-s-agent-s-shell-mounts-the-main-checkout-as-an-empty.md), "A docker run from a story's agent's shell mounts the main checkout as an empty folder, so a hand call in the agent image cannot read the project". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0276 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0276 is closed with `flai issue close I-0276 --reason` saying what fixed it

## Tasks
- T-2926 The cause of I-0276's empty mount is found on the host and a fix is proposed
- T-2927 A hand-run container sees the main checkout's files, and refuses with the cause when a mount comes up empty
- T-2929 The runbook says why a hand docker run typed in zsh can mount an empty folder, and I-0276 is closed
- T-3078 The tooling convention says to brace a shell variable before a colon, since zsh reads $name:r as a modifier

## Notes

Cost of delay inputs set by flai from I-0276. time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T05:43:57Z, 3.1 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0667 on 2026-11-10. The story declared no touches. `flai touches suggest` needs a start path, so it was given `scripts/agent-container.sh` and the runbook. It found 10 commits that changed them.

| Touch | Source | Why |
|-------|--------|-----|
| `scripts/agent-container.sh` | layout | The repository's hand-run form of a container start (S-0505), and the place a check for an empty mount goes (T-2927) |
| `docs/operators/runbooks/agent-container.md` | design, co-change | The runbook for hand-run containers, linked from `design/system/strategic-agents.md` § The live run (S-0567) (T-2929) |
| `design/tech/docker.md` | co-change (30%) | Where the fact about the host's Docker that the cause rests on belongs (T-2929) |
| `design/issues/I-0276-…empty.md` | design | The issue the goal names, bumped by T-2926 and closed by T-2929 |
| `design/issues/summary.md` | layout | `flai issue close` regenerates it |

Co-change also listed `design/system/flai-cli.md`, `docs/operators/index.md`, `docs/users/flai.md`, and `scripts/README.md` at 40%. They are left out: no task changes a command or adds a script. A story's agent widens the touches if T-2926's fix needs one of them. There are no folder touches.

The cause is not known yet. S-0567's agent ran on the host on `harness: claude-code`, and its `docker run -v <checkout>:<checkout>:ro` came up empty inside and outside the sandbox. A container `flai serve` started for the planner, with the same path, had the files. T-2926 compares the two, by docker context, `DOCKER_HOST`, and daemon, before anything is built.

Figures:

- Forecast: `flai forecast` gave 13m: 107 s a unit over 30 done improvement stories, times size 7 (2 criteria, 5 touches). Raised to 1h30m, since T-2926 runs docker live on the host from several shells before T-2927 can start. S-0567's live run was forecast at 1h. Delivery is moved from flai's 11:09Z by the extra 77 minutes.
- Cost of delay: `flai cod` gave 50.00 USD a week from flai's input: 20m lost per 168h cycle at 150 USD an hour. Kept as given: I-0276 has one occurrence, and the cost of a hand container that cannot read the project is about that.
- No `after`. No other story writes output this one reads. S-0670 and S-0682 change `scripts/agent-container.sh` too, which a trial merge judges.

### Cause

Found by T-2926 on 2026-11-10, from agent-S-0667's shell on the host. Docker is not the cause. The agent's shell is zsh. In zsh, `$name:r` is a history modifier, even inside double quotes: it gives the name's value with its extension taken off. So `-v "$ROOT:$ROOT:ro"` hands docker `<ROOT>:<ROOT>o`. The `:r` of `:ro` is eaten, and an `o` is left on the target. bash and `sh` keep the spec as written.

```text
$ echo "$0 $ZSH_VERSION"
/usr/bin/zsh 5.9
$ ROOT=/home/alex/git/bytepunx/system-flow
$ echo "$ROOT:$ROOT:ro"
/home/alex/git/bytepunx/system-flow:/home/alex/git/bytepunx/system-flowo
$ echo "${ROOT}:${ROOT}:ro"
/home/alex/git/bytepunx/system-flow:/home/alex/git/bytepunx/system-flow:ro
$ bash -c 'ROOT=…; echo "$ROOT:$ROOT:ro"'
/home/alex/git/bytepunx/system-flow:/home/alex/git/bytepunx/system-flow:ro
```

Docker then mounts the checkout read-write at `<ROOT>o`, and nothing at `<ROOT>`. `docker inspect` of the container shows it:

```text
HostConfig.Binds: ["/home/alex/git/bytepunx/system-flow:/home/alex/git/bytepunx/system-flowo"]
Mounts: Source /home/alex/git/bytepunx/system-flow -> Destination /home/alex/git/bytepunx/system-flowo, RW true
```

With `-w "$ROOT"`, docker makes `<ROOT>` as an empty working directory, which is the empty folder I-0276 saw. Without it, `ls "$ROOT"` says `No such file or directory`. The braced form lists the checkout:

```text
$ docker run --rm --entrypoint ls -w "$ROOT" -v "$ROOT:$ROOT:ro" flai-agent-sf-claude-code:local -la "$ROOT"
total 8
drwxr-xr-x    2 root     root          4096 Nov 10 11:01 .
$ docker run --rm --entrypoint ls -w "$ROOT" -v "${ROOT}:${ROOT}:ro" flai-agent-sf-claude-code:local -la "$ROOT"
total 160   (CLAUDE.md, LICENSE.md, Makefile, … owned by 1000)
```

The docker context (`default`), the daemon (`docker-desktop`, Docker Desktop 29.7.2, client 29.2.1), and `DOCKER_HOST` and `DOCKER_CONTEXT` (both unset) change nothing. So does `--mount type=bind,src=…,dst=…,readonly`, which has no `:r` to eat. Only an unbraced name right before `:r`, `:h`, `:t`, `:e`, or another zsh modifier letter is changed: `$ROOT:/x:ro`, `$ROOT:$ROOT/.git:ro`, and `/tmp:/tmp:ro` come through whole. That is why the sandbox made no difference, and why `flai serve`'s containers have the files: `flai serve` and `scripts/agent-container.sh` build docker's arguments in Go and in POSIX `sh`, which zsh never reads.

Fix: write `${name}` before a colon in a hand-written shell command. The fix goes in this repository's tooling convention, as a project addition below its marker, beside the zsh rule (T-3078). It also goes in the runbook, for an operator who types a `docker run` by hand (T-2929). No script or Go code holds the cause, so a test in the repository does not fit. The commands above reproduce it, and this note records them. Nothing outside the repository needs changing, so no thread was asked.
