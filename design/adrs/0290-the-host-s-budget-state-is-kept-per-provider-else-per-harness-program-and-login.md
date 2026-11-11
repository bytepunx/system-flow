---
id: ADR-0290
title: "The host's budget state is kept per provider, else per harness program and login, a harness change clears the old identity's state, and flai serve budget clear clears one by hand, which flai guard refuses to served agents"
status: accepted
date: 2026-11-11
supersedes: []
superseded_by: []
refines: [ADR-0209]
---

# ADR-0290 The host's budget state is kept per provider, else per harness program and login, a harness change clears the old identity's state, and flai serve budget clear clears one by hand, which flai guard refuses to served agents

## Context

ADR-0209 §5 made `flai serve` hold every start on a provider, on every project the host serves, once a run on it ends on the usage limit. A run that names no provider was kept under its harness's name, such as `claude-code`. On 2026-11-11 a run of the program `rmhclod`, one Claude login, ended on the limit, and the hold then held a project whose `claude-code` program was `meclod`, another login. The operator's change of that project's program from the dashboard did not lift it, and nothing but removing `~/.flai/serve/budget.json` by hand did (S-0720).

## Decision

The host's budget state is kept per provider, else per harness program and login; a change of a harness's program, arguments, or login clears the old identity's state; and `flai serve budget clear` clears one by hand, which `flai guard` refuses to every agent `flai serve` started.

1. A budget state is kept under the provider a run names; else under its harness, its program (the operator's, else the adapter's default), and, when the host sets one, its login: such as `claude-code meclod` or `container docker login /home/alex/.claude-rmh`. The program's arguments are not part of the key. The key is recorded on the run at its start's gate (`budget_key` in `agents.json`), so its end is kept under the same key.
2. A state of the shape before S-0720, kept under the bare harness name, is still read, and holds every start on its harness, whatever its program and login, until it clears: which login reached the limit cannot be told.
3. `flai serve agent harness` and `flai serve agent set`, the dashboard's harness setting among them, clear the state of the harness's old identity, and one of the old shape, when they change its program, its arguments, or its login, for a project or every project. Each clear is journaled under `serve.budget` with why and sent as a notice, as `flai serve`'s own clears are.
4. `flai serve budget` shows the states, and `flai serve budget clear [<key>]` clears one, or every state, by hand, journaled and notified the same way. A hold's message names that command.
5. `flai guard` refuses `flai serve budget clear` to every session `flai serve` started and every sub-agent, as ADR-0220 refuses a write of the host configuration: an agent the hold holds does not lift it for itself. The bare `flai serve budget` stays a read.
6. `budget.json` is changed under an OS file lock beside the in-process lock, since the clears run in a `flai` process other than `flai serve`'s.

## Consequences

- Two projects on two Claude logins of one harness hold apart; two on one login share the hold, as two on one provider do.
- A harness change of arguments clears the state of an unchanged key, which is the operator's say that the account may have changed.
- An old-shape state left from before an upgrade holds its whole harness once more until a probe ends clean, the operator clears it, or a harness change does.

## Alternatives considered

- Keeping the arguments in the key: a key with every argument is long and unreadable, and arguments rarely name the login; a change of arguments clears the state instead.
- Dropping old-shape states on read: it would let the login that reached the limit start again at once, against the limit.
- Letting served agents clear a hold: an agent could lift the hold that holds its own starts, which is the operator's to lift.
