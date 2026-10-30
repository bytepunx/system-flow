---
id: ADR-0220
title: "flai guard refuses every sub-agent and every session flai serve started the flai serve, flai config, and flai template commands that write the host configuration or turn a host action on or off"
status: accepted
date: 2026-10-30
supersedes: []
superseded_by: []
refines: [ADR-0060]
topics: [cli]
---

# ADR-0220 flai guard refuses every sub-agent and every session flai serve started the flai serve, flai config, and flai template commands that write the host configuration or turn a host action on or off

## Context

The host configuration, `~/.flai/config.json` or the file `FLAI_CONFIG` names, holds what the operator lets `flai serve` do on this host: the host actions each project has on, the agent command and its harnesses and providers, the checks, the components server's limits, the forge's token variable, the reminder channel, and the folders offered for import ([ADR-0008](0008-cli-config-location.md), [ADR-0039](0039-a-settings-host-action-turned-on-only-in-a-shell-lets-the-dashboard-change-the.md)). Each setting has its own `flai serve` sub-command, run in a shell on the host.

`flai guard` had no rule for them ([I-0210](../issues/I-0210-flai-guard-lets-a-session-flai-serve-started-run-the-flai-serve-sub-commands-that-write-the-host-configuration.md)). A story's agent, whose role has no write rule of its own, could run `flai serve notify set`, `flai serve forge`, `flai serve agent provider`, or `flai serve enable`. `agentEnv` in `flai/internal/serve/agents.go` drops `FLAI_CONFIG` from a served agent (S-0183, I-0044), so such a call lands in the operator's own file. With it an agent could widen its own permissions, as `flai serve enable auto-approve` does, or name a program the host runs as the operator, as `flai serve notify set --command` and `flai serve agent set --` do. That breaks `safety.md`, which forbids an agent to change the operator's global tool configuration, and [ADR-0205](0205-flai-serve-reminds-the-operator-of-a-thread-that-awaits-them-past-threads.md), which leaves the reminder channel to the operator's shell. The `settings` host action is the dashboard's way to change the host's settings, not an agent's.

`flai config set` and `flai template use` write the same file through the same gap.

## Decision

`flai guard` refuses every sub-agent and every session `flai serve` started (`FLAI_STARTED_BY=flai-serve`, `FLAI_ROLE`, or `FLAI_STORY` set), whatever its role or permissions, each command that writes the host configuration or turns a host action on or off, in any form a shell line takes, as it refuses the shared paths ([ADR-0096](0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md)) and `flai forge push` ([ADR-0166](0166-under-pull-request-acceptance-flai-serve-pushes-a-story-s-branch-and-opens-or.md)). The rule runs before the sub-agent and role split, so it binds the story's agent and every strategic role alike.

Refused, as writes:

| Command | Forms |
|---------|-------|
| `flai serve enable`, `flai serve disable` | every form |
| `flai serve agent` | `set`, `clear`; `harness` with a name and a setting (`--program`, `--guard`, `--deny-protected`, or arguments after `--`) or with `--reset`; `provider` with a name and a setting (`--api`, `--base-url`, `--key-env`, `--spend-key-env`, `--model`) or with `--reset` |
| `flai serve checks` | `set`, `timeout`, `clear` |
| `flai serve components` | `idle`, `max-servers` |
| `flai serve forge` | with a setting (`--token-env`) or `--reset` |
| `flai serve notify` | `set`, `clear` |
| `flai serve import` | `add`, `remove` |
| `flai config` | `set` |
| `flai template` | `use` |

Left open, as reads: `flai serve status`, `actions`, `journal`, `agent show`, `agent harness` and `agent provider` bare or with a name alone, `checks show`, `components show`, `forge` bare or with a forge alone, `notify show`, `import list`, and `flai config get` and `path`. Left to their own rules, because they write no host configuration: `flai serve start`, `stop`, `project`, `usage`, `agent start`, `restart`, `commit`, `stop`, `stream`, and `orchestrate`, and the other `flai template` sub-commands. A read open here is still refused to a sub-agent by ADR-0060's rule for any flai command it does not list as a read.

The refusal names the host configuration and this ADR, and says what to do instead: a sub-agent puts what it needs in its final message, and a served session asks the operator on a thread with `thread_open`. The orchestrator's refusal names the call, so that its activity document logs it. The operator's own session, which `flai serve` did not start, is not held.

## Consequences

- No agent `flai serve` starts, and no sub-agent, changes what `flai serve` may do on the host, or names a program it runs; only the operator does, in a shell or through the dashboard's `settings` host action.
- An agent that needs a host setting changed asks for it on a thread, as it asks for a shared path.
- `flai serve project add` and `remove` change the registry of served projects, not the host configuration, and stay out of this rule, as do `flai serve stop` and `flai host stop`, which stop the operator's dashboards; a later issue may hold them.
- `flai guard`'s help, `docs/users/flai.md`, `docs/operators/index.md`, `docs/operators/settings.md`, `design/system/flai-cli.md`, and `design/system/agent-adapters.md` say which commands the guard refuses an agent.

## Alternatives considered

- Pass `FLAI_CONFIG` to a served agent, pointing at a copy: the agent's flai would then miss the providers and harnesses the host configures, and a write would still change what a later read of the copy says, without the operator.
- Refuse every `flai serve` and `flai config` command to an agent: it would close the reads an agent uses to see why a host action is off, and `flai serve start`, `restart`, and `commit`, which keep their own rules.
- Leave `flai config set` and `flai template use` to a separate issue: they write the same file through the same gap, so a guard that refuses one path to it and leaves the other open does not close the risk.
