---
id: T-2956
type: task
nature: remediation
title: scripts/agent-image-test.sh reproduces I-0339 without Docker, and the smoke tier runs it
status: done
parent: S-0713
owner: alex
created: 2026-11-10T10:23:58Z
updated: 2026-11-10T10:55:45Z
transitions:
  - to: ready
    at: 2026-11-10T10:54:26Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T10:54:27Z
    by: system-flow
  - to: done
    at: 2026-11-10T10:55:45Z
    by: agent-S-0713
stream: S-0713
tags: [scripts, test]
touches: [scripts/agent-image-test.sh, scripts/smoke.sh]
after: [T-2955]
usage:
  source: log
  seconds: 78
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 1188
      cache_read: 667139
      cache_write: 11983
      cost: 0.2531
---
# T-2956 scripts/agent-image-test.sh reproduces I-0339 without Docker, and the smoke tier runs it

## Work

Write `scripts/agent-image-test.sh` the way `scripts/agent-container-test.sh` reproduces I-0279: a scratch copy of `env.sh` and `agent-image.sh`, a fixture `system-flow.yaml` with no layer, and a scratch `HOME` whose `.docker/config.json` names `credsStore: desktop.exe`. Put on `PATH` a fake `docker` that writes `DOCKER_CONFIG` and the contents of the configuration it names to a log and exits 0. Run the script three times:

- with a fake `docker-credential-desktop.exe` that exits 1, as in I-0339;
- with no such helper on `PATH`;
- with a helper that answers.

The first two must log a configuration of `{}` and print the line naming the helper. The third must leave `DOCKER_CONFIG` unset. Every run must leave the scratch `HOME`'s `config.json` byte for byte as it was. Then add the test to `scripts/smoke.sh`, beside `agent-container-test.sh`. It waits for T-2955, whose behaviour it checks.

## Done when

- `scripts/agent-image-test.sh` passes against T-2955's script.
- With T-2955's change reverted, the test fails on the I-0339 case.
- `scripts/smoke.sh` runs the test.
- ShellCheck passes on both scripts.

## Notes
