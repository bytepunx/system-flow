---
id: T-1941
type: task
nature: remediation
title: scripts/litellm.sh stops freezing its default haiku model into admin.env, and replaces a retired one on up
status: done
parent: S-0455
owner: alex
created: 2026-10-21T00:39:43Z
updated: 2026-10-21T02:14:32Z
transitions:
  - to: ready
    at: 2026-10-21T02:10:07Z
    by: agent-S-0455
  - to: in-progress
    at: 2026-10-21T02:10:07Z
    by: agent-S-0455
  - to: done
    at: 2026-10-21T02:14:32Z
    by: agent-S-0455
stream: S-0455
tags: [scripts, litellm]
touches: [scripts/litellm.sh, scripts/README.md, docs/contributors/index.md]
usage:
  source: log
  seconds: 265
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 13
      output: 4720
      cache_read: 1032419
      cache_write: 37839
      cost: 0.5786
---
# T-1941 scripts/litellm.sh stops freezing its default haiku model into admin.env, and replaces a retired one on up

## Work

The host keeps failing smoke after S-0451's T-1928 changes the default (I-0171). On the first `up`, `scripts/litellm.sh` writes `LITELLM_HAIKU_MODEL` into `.flai-cache/litellm/admin.env`. Every later `up` reads it from there and writes it into `litellm.env` as `FLAI_TEST_LITELLM_HAIKU`, which `scripts/env.sh` exports and which beats the test's built-in default. So a host keeps the retired `nvidia/nemotron-3-super-120b-a12b:free` after the default moves on.

- `admin()` records `LITELLM_HAIKU_MODEL` in `admin.env` only when the operator set it. Otherwise the script's current default applies on every `up`.
- On `up`, a recorded value equal to a default the script has retired is removed from `admin.env`, and the script says so. Keep the retired defaults in one list in the script, with the nemotron-3-super slug in it.
- `litellm.env` is rewritten on every `up` with the model in effect, keeping the virtual key as it does now.
- Update the script's header comment, its row in `scripts/README.md`, and the `FLAI_TEST_LITELLM_HAIKU` rows in `docs/contributors/index.md`.

Waits for nothing in this story. It runs beside the gateway test's probe, which touches no path in common. The story waits for S-0451, whose T-1928 changes the same default.

## Done when

- Run against a scratch copy of `admin.env` that names the retired slug, `scripts/litellm.sh up` leaves `admin.env` without it, and `litellm.env` names the current default.
- A `LITELLM_HAIKU_MODEL` the operator set is kept across `up`.
- `sh -n scripts/litellm.sh` passes, and the markdown lint passes on the two docs.

## Notes
