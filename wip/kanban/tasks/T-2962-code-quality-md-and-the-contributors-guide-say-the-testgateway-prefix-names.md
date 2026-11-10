---
id: T-2962
type: task
nature: remediation
title: code-quality.md and the contributors guide say the TestGateway prefix names only the paid checks the gateway tier runs
status: in-progress
parent: S-0708
owner: alex
created: 2026-11-10T10:27:42Z
updated: 2026-11-10T10:54:02Z
transitions:
  - to: ready
    at: 2026-11-10T10:54:01Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T10:54:02Z
    by: system-flow
stream: S-0708
tags: [docs, testing]
touches: [design/conventions/code-quality.md, docs/contributors/index.md]
---
# T-2962 code-quality.md and the contributors guide say the TestGateway prefix names only the paid checks the gateway tier runs

## Work

Below the marker in `design/conventions/code-quality.md`, where the project additions name the paid gateway tests, add one sentence: a test named `TestGateway*` is a paid check that `scripts/gateway-smoke.sh` runs by name, since the go-test and integration tiers skip `^TestGateway`; a gateway test that needs no network takes another name, and a test in `flai/internal/serve` fails on a `TestGateway*` function the script does not run. Add the same rule, in one sentence, to the gateway section of `docs/contributors/index.md`, where it says the go-test and integration tiers skip `^TestGateway`. Edit nothing above the marker. It waits for no task: it touches no file the test task touches, and the rule it states is the story's plan.

## Done when

- Both documents state the rule, each once, with no copied paragraph.
- `flai test design/conventions/code-quality.md docs/contributors/index.md` passes, the markdown lint included.

## Notes
