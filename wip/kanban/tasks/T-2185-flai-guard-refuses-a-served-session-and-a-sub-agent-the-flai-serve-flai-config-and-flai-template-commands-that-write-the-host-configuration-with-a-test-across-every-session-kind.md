---
id: T-2185
type: task
nature: remediation
title: flai guard refuses a served session and a sub-agent the flai serve, flai config, and flai template commands that write the host configuration, with a test across every session kind
status: done
parent: S-0511
owner: alex
created: 2026-10-28T08:19:48Z
updated: 2026-10-30T00:24:00Z
transitions:
  - to: ready
    at: 2026-10-30T00:13:46Z
    by: agent-S-0511
  - to: in-progress
    at: 2026-10-30T00:13:47Z
    by: agent-S-0511
  - to: done
    at: 2026-10-30T00:24:00Z
    by: agent-S-0511
stream: S-0511
tags: [cli, guard]
touches: [flai/internal/guard/hostconfig.go, flai/internal/guard/hostconfig_test.go, flai/internal/guard/call.go, flai/internal/guard/guard.go]
after: [T-2184]
usage:
  source: log
  seconds: 613
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 47
      output: 13320
      cache_read: 2283632
      cache_write: 111126
      cost: 1.4356
---
# T-2185 flai guard refuses a served session and a sub-agent the flai serve, flai config, and flai template commands that write the host configuration, with a test across every session kind

## Work

Add the rule T-2184's ADR records, in a file of its own, `flai/internal/guard/hostconfig.go`, modelled on `shared` in `shared.go` and `forgePush` in `guard.go`: a `rules` value whose `flai` func matches `cmd == "serve"` with the writing sub-commands, `config set`, and `template use`, and a method `(g Guard) hostConfig(c Call) Verdict` that returns the zero Verdict unless `c.Subagent` or `g.served(c)`, and otherwise `deny` with a reason that names the host configuration, the ADR, and what to do instead: a sub-agent puts what it needs in its final message, a served session asks the operator on a thread with `thread_open`; the orchestrator's refusal sets `Call` so its activity document logs it. `refuse` hands the rule `(cmd, sub, rest)` and `subcommands()` reads two words, so read the third word and the flags from `rest`: `serve agent set`, `serve agent harness <name> --program ...`, `serve agent provider <name> --api ...` write, while `serve agent provider`, `serve agent provider <name>`, `serve agent harness`, `serve forge`, and `serve forge <forge>` read, as `flagReads` tells reads from writes by flags; `--reset` makes each of those a write. Wire it in `Decide` in `call.go` beside `shared` and `forgePush`, before the sub-agent and role split, so it binds the story's agent, whose role has no write rule today, and every strategic role. Add a paragraph to the package comment in `guard.go` after the `flai forge push` one. Write `hostconfig_test.go` as `TestFlaiForgePushIsTheOperators` is written: append `"serve"`, `"config"`, and `"template"` to `g.Commands`; a sessions table with the story's agent, the planner, the orchestrator with `allOn`, the analyzer, the curator, a session served with no role, each as a sub-agent too, and the operator; the writes refused for each served session and sub-agent in every form the forge test covers (`scripts/flai.sh`, an absolute path, `--config`, `env`, `bash -c`, a chain after a read), the reads passing, and the operator passing everything; and the orchestrator's refusal carrying `Call`. The story's agent's `flai serve notify set --url x` passing through `Decide` before the change is the reproduction I-0210 asks for. Waits for T-2184, which fixes the list of forms.

## Done when

`flai test flai/internal/guard` passes with the new test, every served session and sub-agent is refused each writing form and allowed each reading form named in the ADR, the operator's own session is held by nothing new, and `Decide`'s comment names the rule.

## Notes

`flai serve stop` and `flai host stop` stop the operator's dashboards but write no configuration; they are out of this rule, as the ADR says. Leave `cliReads` alone: a sub-agent's `flai serve status` is refused today as any write is, and widening sub-agent reads is not this story's.
