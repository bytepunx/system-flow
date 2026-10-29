---
id: T-2314
type: task
nature: improvement
title: flai-cli.md, strategic-agents.md, agent-context.md, the strategic-agents convention and its template copy, and the user docs describe the reason-primed pack
status: done
parent: S-0537
owner: alex
created: 2026-10-28T22:08:08Z
updated: 2026-10-29T16:26:55Z
transitions:
  - to: ready
    at: 2026-10-29T16:12:27Z
    by: agent-S-0537
  - to: in-progress
    at: 2026-10-29T16:12:28Z
    by: agent-S-0537
  - to: done
    at: 2026-10-29T16:26:55Z
    by: agent-S-0537
stream: S-0537
tags: [docs, orchestrator, template]
touches: [design/system/flai-cli.md, design/system/strategic-agents.md, design/system/agent-context.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, docs/users/flai.md, docs/users/flai-reference.md, template/CHANGELOG.md, design/system/conventions.md]
after: [T-2310, T-2311, T-2312]
usage:
  source: log
  seconds: 867
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 84
      output: 31262
      cache_read: 6247071
      cache_write: 201856
      cost: 3.0222
---
# T-2314 flai-cli.md, strategic-agents.md, agent-context.md, the strategic-agents convention and its template copy, and the user docs describe the reason-primed pack

## Work

Criterion 4. Describe the reason-primed pack where the general orchestrate pack is described:

- `design/system/agent-context.md` § Strategic packs: what a reason pack loads per reason, that the conventions load once with the sections the reasons' topics leave out taken out, that the digest ranks against the reasons' threads only, and its measured size against the general pack.
- `design/system/flai-cli.md` § Commands, the `flai prime` rows: `--reason` and its forms.
- `design/system/strategic-agents.md` § The orchestrator › What it is told: the start prompt names the reasons and the orchestrator primes with them.
- The strategic-agents convention's priming rule (`## Rules`), here and in the template, which today tells the orchestrator to prime with `orchestrate` alone.
- `docs/users/flai.md` § Prime a session; regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- A line in `template/CHANGELOG.md`.

Waits for T-2310, which changes the same convention files and the changelog, and for T-2311 and T-2312, whose flag and field it documents.

## Done when

- Each document names `--reason` or `reasons`, the reason kinds, and what each loads, and links ADR-0209.
- `flai-reference.md` is regenerated and its test passes.
- `scripts/flai.sh check --strict` and the markdown lint are clean.

## Notes
