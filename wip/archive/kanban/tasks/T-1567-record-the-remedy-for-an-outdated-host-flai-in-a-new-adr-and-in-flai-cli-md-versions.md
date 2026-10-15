---
id: T-1567
type: task
nature: improvement
title: "Record the remedy for an outdated host flai in a new ADR and in flai-cli.md § Versions"
status: done
parent: S-0387
owner: alex
created: 2026-10-14T20:54:28Z
updated: 2026-10-15T01:42:39Z
transitions:
  - to: ready
    at: 2026-10-15T01:41:24Z
    by: agent-S-0387
  - to: in-progress
    at: 2026-10-15T01:41:24Z
    by: agent-S-0387
  - to: done
    at: 2026-10-15T01:42:39Z
    by: agent-S-0387
stream: S-0387
tags: [adr, design]
touches: [design/adrs/0153-flai-serve-upgrades-the-host-s-flai-when-the-project-s-history-has-a-newer-published-release-behind-auto-upgrade.md, design/adrs/README.md, design/system/flai-cli.md, design/adrs/0158-flai-serve-checks-each-served-project-for-a-newer-flai-release-every-minute.md]
usage:
  source: log
  seconds: 75
  estimated: true
  priced_by: harness
  models:
    - model: claude-haiku-5-5
      input: 160
      output: 54996
      cache_read: 11528857
      cache_write: 238837
      cost: 0.7525
    - model: claude-opus-5-5
      input: 18
      output: 7035
      cache_read: 1029379
      cache_write: 29107
      cost: 0.5324
---
# T-1567 Record the remedy for an outdated host flai in a new ADR and in flai-cli.md § Versions

## Work

Confirm or change the remedy proposed under the story's `### Planning`, from I-0141's two instances, and record it with `flai adr new` and `flai adr accept`, refining ADR-0142 and ADR-0108. Decide:

- that `flai serve` checks `release.FlaiOutdated` for each served project on its reconcile loop, at most once a minute as the cache already allows, not only when it starts serving it;
- the new host action `auto-upgrade`, kept to the shell (`shellOnly`), with which `flai serve` asks the host to upgrade itself to the newest published, verifiable release when the running flai is behind the project's history, and when it does so (a minor behind only, or a patch too);
- what it does with the action off: one thread asking the operator to run `flai host upgrade`, its anchor, and that it opens no second while one is open;
- the message `flai test` and the MCP tools `test` and `verify` give when the manifest's tiers refuse and `FlaiOutdated` answers.

Update `design/system/flai-cli.md` § Versions: the host's flai and the tree to match, and the ADR index. The predicted ADR path is a guess at its number and slug; `flai task done` widens the touches with the real one.

It waits for nothing: every other task builds what it decides.

## Done when

- The ADR is accepted and listed in `design/adrs/README.md`.
- `design/system/flai-cli.md` § Versions describes the periodic check, `auto-upgrade`, the thread, and the version-aware tier refusal.
- `flai check --strict` and the markdown lint pass on both files.

## Notes
