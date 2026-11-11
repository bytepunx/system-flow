---
id: ADR-0288
title: "The MCP tool test answers a tier's command only for the tier that failed, gives every other tier the number of arguments its command ran with, and fits its answer to 40,000 bytes of JSON by cutting paths and then the failed command"
status: accepted
date: 2026-11-11
supersedes: []
superseded_by: []
refines: [ADR-0104, ADR-0179]
---

# ADR-0288 The MCP tool test answers a tier's command only for the tier that failed, gives every other tier the number of arguments its command ran with, and fits its answer to 40,000 bytes of JSON by cutting paths and then the failed command

## Context

The MCP tool `test` answered `verify.Result` as `flai test --json` prints it (S-0273). Each tier's `command` is its argv with `{packages}`, `{importers}`, or `{files}` filled in, so a tier given every file it selects, such as gofmt, flaiover's lint, or vitest, lists every file in its command. `paths` lists them once more. Claude Code saves a tool result over 50,000 characters to a file and shows only a preview ([ADR-0104](0104-the-mcp-prime-tool-returns-a-context-pack-in-parts-of-at-most-40-000-bytes-of.md)).

[I-0260](../issues/I-0260-the-mcp-tool-test-s-answer-on-a-pass-runs-to-68-kb-over-claude-code-s-tool.md) records S-0425, whose call over nine paths passed every tier and answered 68,784 characters: gofmt's `command` was 14 KB, flaiover-lint's 10 KB, and vitest's 10 KB. The agent wrote `flai test --json` to a file and read it with `jq`.

A tier that passed needs no command: nothing is rerun from it. The tier that failed does, since a rerun by hand needs the exact argv.

## Decision

The MCP tool `test` answers `flai test --json`'s result with three cuts, made in `flai/internal/mcpserver/test.go` alone.

- **The command.** Only the tier whose state is `failed` keeps `command`, as it ran. Every other tier, passed, not reached, or skipped, gives `args` in its place: the number of arguments its command ran with. Each tier keeps its name, `dir`, state, exit code, durations, findings, `omitted`, `covered_by`, and `log`.
- **The fit.** The whole answer fits 40,000 bytes of JSON, `context.PartLimit`, as `task_done`'s and `story_start`'s do ([ADR-0179](0179-task-done-and-story-start-carry-an-inbox-that-lists-the-story-s-threads-and.md)). When it would not, `paths` is cut to the most of its first entries that fit, and `paths_omitted` counts the rest. When it still would not, the failed tier's `command` is cut the same way, and its `command_omitted` counts the arguments left out. Each count is omitted when zero.
- **The CLI.** `flai test --json`, the dashboard's `test.run`, and the verify record keep every tier's whole `command` and every path: a shell caller can write them to a file. The tool's description says the answer is `flai test --json`'s with these cuts.

This refines ADR-0104, whose limit now bounds this answer too, and ADR-0179, whose count-and-omit form it takes, as ADR-0198 did for `task_done`'s notes.

## Consequences

- A passing run over many files answers within Claude Code's limit, and the agent reads it without `jq`. `TestTestAnswersACommandOnlyForTheTierThatFailedAndFitsTheLimit` reproduces I-0260 with a tier given 1,500 files.
- The findings are not cut: they are bounded by `max`, which the caller sets, 5 by default.
- An agent that wants a passing tier's argv runs `flai test --json`.

## Alternatives considered

- **Cut the commands in `verify.Result` itself.** Rejected: `flai/cmd/test_test.go` and `flai/internal/verify/run_test.go` rely on the whole command, and the verify record is read back by the close-out.
- **Keep every command, cut only to fit.** Rejected: it would still answer up to 40,000 bytes of argv the agent cannot act on, in place of the paths it can.
- **Give the number of files a tier was given rather than its arguments.** Rejected: a tier's arguments mix its own words with what `{packages}`, `{importers}`, and `{files}` expand to, and only the count of all of them is known without parsing the manifest's command again.
