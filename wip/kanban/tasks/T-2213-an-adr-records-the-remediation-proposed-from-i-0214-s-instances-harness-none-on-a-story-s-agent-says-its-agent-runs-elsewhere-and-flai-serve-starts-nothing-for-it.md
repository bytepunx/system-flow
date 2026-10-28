---
id: T-2213
type: task
nature: remediation
title: "An ADR records the remediation proposed from I-0214's instances: harness none on a story's agent says its agent runs elsewhere, and flai serve starts nothing for it"
status: backlog
parent: S-0516
owner: alex
created: 2026-10-28T13:13:11Z
updated: 2026-10-28T13:13:11Z
transitions: []
stream: S-0516
tags: [cli, serve, adr]
touches: [design/adrs/README.md]
---
# T-2213 An ADR records the remediation proposed from I-0214's instances: harness none on a story's agent says its agent runs elsewhere, and flai serve starts nothing for it

## Work

Read I-0214's two instances and the code they name: `look` and `restartEnded` in `flai/internal/serve/agents.go`, and `harness.For` in `flai/internal/harness/harness.go`, where today a story with no harness is started with the host's command and any other name but `claude-code` and `command` fails the start as a failed run. Propose the remediation from the instances and record it with `flai adr new`, then `flai adr accept`: a reserved harness name, `none`, on a story's agent (or the project's default) says the story's agent runs elsewhere, by hand or in a container nobody has made a harness for yet. `flai serve` starts nothing for such a story: not when it enters ready, not as an automatic restart (ADR-0108), and not on a resume or a send back; it names the reason in `waiting`, as it names a story with no harness and no command, and that is no failure. `flai serve agent start`, `restart`, and `commit` refuse it, saying its agent runs elsewhere. `flai story start` and the MCP tool `story_start` work as they do, since they start no agent. The manifest's shape check accepts the name as it accepts any; `harness.For` refuses it with the same words, so that nothing above the adapters starts it by mistake. Say in the ADR why a reserved harness name rather than a new field (a story's agent is the one place that says where it runs, and only `flai serve`'s adapters read it) and that `none` stays once S-0508 makes the container a harness. Add the ADR to `design/adrs/README.md`'s index as `flai adr new` does, and give it topics with `flai adr topics` so that `harness` stories prime with it. If the operator answered the plan's thread with another remediation, record that one instead and adjust the tasks that follow.

## Done when

- An accepted ADR under `design/adrs/` records the remediation, its context (I-0214's two instances), the alternatives (a new field on the story; a host setting; waiting for S-0508), and its consequences for `flai serve`, the operator's starts, and the docs
- `design/adrs/README.md` lists it, and `flai check --strict` is clean on the ADR

## Notes

The ADR's file is written under a number and slug no task can name in advance; the story keeps the folder `design/adrs` for it, and `flai task done` widens the touches with the committed path (ADR-0144). Layer 1: the tasks that build the remediation wait for this one, since it fixes the name and the words.
