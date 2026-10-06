---
id: T-0952
type: task
nature: feature
title: flai manifest set writes the strategic agents' settings to system-flow.yaml, validating each and refusing with the field and the reason
status: done
parent: S-0229
owner: alex
created: 2026-10-05T05:46:14Z
updated: 2026-10-06T21:41:09Z
transitions:
  - to: ready
    at: 2026-10-06T21:30:47Z
    by: agent-S-0229
  - to: in-progress
    at: 2026-10-06T21:30:47Z
    by: agent-S-0229
  - to: done
    at: 2026-10-06T21:41:09Z
    by: agent-S-0229
stream: S-0229
tags: [flai]
touches: [flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/cmd/manifest.go, flai/cmd/manifest_test.go, flai/cmd/root.go, docs/users/flai-reference.md, docs/operators/settings.md]
usage:
  source: log
  seconds: 622
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 83
      output: 381
      cache_read: 4172050
      cache_write: 181401
      cost: 1.905
---
# T-0952 flai manifest set writes the strategic agents' settings to system-flow.yaml, validating each and refusing with the field and the reason

## Work

ADR-0039 has every dashboard write run a flai command that hostapi builds from checked values, so the manifest's strategic keys need a command before the dashboard can save them. Add `flai manifest set <key>=<value>... [--unset <key>]... [--autocommit] [--trailer] [--json]`; the name is the planner's guess, and the story's agent may choose another.

- `flai/internal/manifest/settings.go`: a catalog of the keys the command may write, each with its kind (boolean, choice, number, duration, cron, agent), its values, its default, a sentence on what it does, and for each `orchestration.permissions` key a sentence on its risk. The keys:
  - `orchestration.permissions.*` (S-0218's T-0881), `orchestration.policy` and `orchestration.release.*` (S-0217's T-0809, with `whole_epics` from S-0222);
  - `planning.agent`, `planning.replan`, `planning.schedule`, `planning.hour_rate`, `planning.cycle`, `planning.default_duration`;
  - `analysis.agent` and `analysis.schedule` (S-0223).
- Take each key's values from the validation those stories wrote (`Planning.Errors()`, `Orchestration.Errors()`, the analysis block's), never from a second list, so the command and `flai check` cannot disagree.
- Write only the keys given, into their block, keeping every other key, the order, and the comments, as `WriteAgent` keeps the rest of the file. `--unset` removes a key, which brings back its default.
- Validate the whole manifest as it would be after the change before writing anything. A refusal names the field and the reason, such as `orchestration.release.value: must be a number of zero or more`, writes nothing, and exits with the code a rule's refusal has. `--json` gives the refusals as `{field, reason}`.
- Refuse a key outside the catalog, `planning.currency` among them: changing it re-denominates every amount on the items.
- `--autocommit` commits `system-flow.yaml` alone, as `flai agent set --autocommit` does.
- Register it in `flai/cmd/root.go` and regenerate `docs/users/flai-reference.md` with `make flai-reference`.

This task waits for nothing in this story. It runs with T-0954, whose paths it does not share.

## Done when

- a test writes a key of each kind into a fixture manifest with comments, and the rest of the file is unchanged byte for byte
- a test refuses a bad value of each kind and an unknown key, naming the field and the reason, and the file is unchanged
- `--unset` removes a key, and `--autocommit` commits only `system-flow.yaml`, each in a test
- `docs/users/flai-reference.md` is current, and `go test ./internal/manifest/ ./cmd/` passes

## Notes
