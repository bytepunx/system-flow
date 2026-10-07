---
id: T-1079
type: task
nature: improvement
title: The conventions, the harness prompt, the design, and the user guide send the agent to flai test for test runs between tasks
status: done
parent: S-0273
owner: alex
created: 2026-10-06T22:52:38Z
updated: 2026-10-07T00:34:35Z
transitions:
  - to: ready
    at: 2026-10-07T00:28:55Z
    by: agent-S-0273
  - to: in-progress
    at: 2026-10-07T00:28:55Z
    by: agent-S-0273
  - to: done
    at: 2026-10-07T00:34:35Z
    by: agent-S-0273
stream: S-0273
tags: [conventions, docs, harness]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/delegation.md, design/conventions/work-management.md, design/conventions/tooling.md, template/root/design/conventions/delegation.md, template/root/design/conventions/work-management.md, template/root/design/conventions/tooling.md, template/CHANGELOG.md, template/template.yaml, design/system/devex.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, scripts/README.md, Makefile]
after: [T-1067, T-1069, T-1070, T-1073]
usage:
  source: log
  seconds: 340
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 104
      output: 45821
      cache_read: 6145362
      cache_write: 147909
      cost: 3.0004
---
# T-1079 The conventions, the harness prompt, the design, and the user guide send the agent to flai test for test runs between tasks

## Work

Criterion 3. This task waits for the command, the MCP tool, and the host method (T-1069, T-1070, T-1073), whose flags, arguments, and answers it describes. It also waits for T-1067, with which it shares `template/CHANGELOG.md`.

- Conventions:
  - `delegation.md`: the agent runs `flai test`, or the MCP tool `test`, for a test run between tasks, rather than `go test`, vitest, golangci-lint, or gofmt by hand, and rather than a sub-agent run. Today lines 23 to 32 hand test runs to the verifier.
  - `work-management.md`: the task cycle's "run the task's tests" step names `flai test`.
  - `tooling.md`: the project additions list `flai test` beside the `make` targets.
  - Copy each baseline change to its `template/root` copy and add a line to `template/CHANGELOG.md`, in the same story, as `CLAUDE.md` says.
- Harness prompt: in `flai/internal/harness/harness.go`, the start prompt's task cycle says to run the task's tests with `flai test` or the tool `test`. Cover it in `harness_test.go`.
- Design: in `design/system/devex.md`, the Test tiers row says the manifest's `tests` declare the tiers and `flai test` runs them by path, and the `make` targets stay for people. In `design/system/flai-cli.md`, add `flai test` under the commands, with the MCP tool and the host method.
- User guide: add `flai test` to `docs/users/flai.md`, and regenerate or extend its entry in `docs/users/flai-reference.md` as that file is maintained.
- `scripts/README.md` notes that `flai test` runs the same tiers by path. Check that every `Makefile` test target still runs, and change the `Makefile` only if one does not.

## Done when

- [ ] Each convention and its `template/root` copy say the same, and `template/CHANGELOG.md` records it.
- [ ] The harness prompt names `flai test`, and `harness_test.go` covers it.
- [ ] `devex.md`, `flai-cli.md`, `flai.md`, and `flai-reference.md` describe the command, the tool, and the method as built.
- [ ] `make test`, `make integration`, `make smoke`, and `make flai-test` run as before.
- [ ] `go test ./internal/harness/...`, `flai check --strict`, and the markdown lint pass.

## Notes

Drafted by the planner for S-0273.
