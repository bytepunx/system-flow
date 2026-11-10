---
title: Tooling
updated: 2026-11-10
audience: agent
order: 90
status: active
topics: [all]
---

# Tooling

Use the project's tools for the project's data. The tools keep the standard true; hand edits drift.

## Rules

- Create and move work items with `flai`: `flai epic new`, `flai story new --epic`, `flai task new --story`, `flai move`, `flai block`, `flai unblock`. Do not hand-write front matter for something a command creates.
- Open and log narratives with `flai stream open` and `flai stream log`. Write a narrative's `## Current state` and `## Next steps` with `flai stream state S-nnnn --current "<text>" --next "<text>"`, or the MCP tool `stream_state`, never by editing the file: it replaces the two sections, or the one you give, and leaves the rest alone. The same call records the other summary sections: `--context` replaces `## Context`, each `--decision` appends a bullet to `## Decisions`, and each `--question` appends a dated bullet to `## Open questions`, above the threads block. Use it where `wip/` is read-only, as in a container; where it is writable, those three sections may still be edited by hand. The `## Log` is appended by `flai stream log`.
- Hand-edit item bodies (goals, criteria, notes) freely. Hand-edit item front matter only for fields no command sets, and run `flai check` afterwards.
- Run `flai check --strict` before reporting a story as done and before any commit. Fix what it finds; do not explain it away.
- flai lints the markdown it writes in `wip/` with the project's markdownlint configuration, because that folder is written in the main checkout, where a story's own lint never runs: `flai check` warns on each finding there (`markdown.MDnnn`), and making an item, a thread entry, or a log entry is refused, with the rule and the line, when it would bring one. Fix what it names; do not work around it. The project's own markdown lint stays the authority for what flai does not check.
- Use `flai board`, `flai show`, and `flai stats` to answer questions about work state instead of grepping files.
- Set `FLAI_AGENT` and `FLAI_SESSION` at session start so narrative entries record who wrote them.
- Archive with `flai archive`, never by moving files.
- Run the dashboard with `flai dashboard` when the operator wants to see the board or the charts; do not build ad hoc views.
- Use the project's build, test, and lint entry points (`Makefile`, `package.json` scripts, `go` commands as documented) rather than reconstructing them. For the tests and lint of what you changed, that entry point is `flai test` with the paths: it runs the tiers the manifest's `tests` declare that those paths select, cheapest first, stops at the first that fails, and answers pass or its first findings; `flai test --all` runs every tier.
- When a tool is missing or the wrong version, say so and use a scratch install; do not silently fall back to hand edits or skip the step. Add or extend a script under `scripts/` so the operator can install it properly.
- If a command is missing for something you do repeatedly, propose it as a story instead of scripting around it.
- Put common tasks in `scripts/` as purpose-named shell scripts, so complex commands and command sequences have one name.
- Write shell commands for zsh as well as bash: the host shell may be zsh. Quote every glob, URL, and variable (`"$id"`, `'https://host/?a=b'`), compare with `[ "$a" = "$b" ]`, not `==`, and write a loop's list out rather than splitting an unquoted variable. Put a sequence of more than a few commands in a script under `scripts/`, POSIX `sh` under `set -eu`, instead of a long chain.
- Gate each step of a chain that must succeed on its own exit code: join the steps with `&&`, or run them in a script under `set -e`, never with `;` or through a pipe that hides a failure (`flai check --strict | tail -1` exits as `tail` does). Check what a commit records, such as a narrative or a generated file, before the commit, not after.
- The `Makefile` is the entry point for building, linting, quality checks, dependency updates, tests, and environment management. Make targets call the scripts in `scripts/` rather than inlining commands. The test targets are `test` (behavior), `integration`, and `smoke`, in that order of cost. They stay for people and CI; an agent runs the tiers through `flai test`.
- Scripts and Make targets work the same locally and in CI.
- Local testing, previews, and validation run in Docker and Docker Compose; when a cluster is needed, use k3d via `bytepunx/kluster`.

## When in doubt

- If `flai` can do it, `flai` does it.
- If a hand edit was unavoidable, `flai check` is not optional.

<!-- system-flow:end-of-baseline -->

## Project additions <!-- topics: all, code -->
- `make check`, `make flai`, `make test`, `make integration`, `make smoke`, `make flai-test` (all tiers plus lint), `make flai-snapshot`, and `make install-tools` call the scripts in `scripts/`.
- `flai test`, run as `scripts/flai.sh test` or `bin/flai test`, runs the tiers in `system-flow.yaml`'s `tests` for the paths given: gofmt, vet, golangci-lint, the Go tests of the changed packages and their importers with `-short` (all but the paid gateway tests and the integration tier's own, those that skip under `-short`: the real-git tests and the four whole-monorepo tests, `TestMonorepoIsClean`, `TestRepositoryLintsClean`, `TestPrimeOnThisRepository`, and `TestPrimeComesInPartsOnThisRepository`, S-0548), the `TestSettingsIndex` tests, `TestReferenceIsCurrent`, and `TestFlagIndexIsCurrent` under `-short` when `docs/operators/settings.md` or `docs/users/flai-reference.md` changes (`settings-index`, I-0304, I-0308), prettier and eslint on the changed flaiover files (`flaiover-lint`), vitest, flaiover's tests of the built flai, every test file that imports `flaiover/src/lib/server/testing.ts`, when a file under `flai/internal/hostapi/` or `flai/cmd/hostapi.go` changes (`flaiover-flai`, I-0304), `flaiover/src/lib/server/design-routes.test.ts` when `design/system/flaiover-dashboard.md` or a file under `flaiover/src/routes/` changes (`design-routes`, I-0322), svelte-check on the whole flaiover project when a changed TypeScript, JavaScript, Svelte, or configuration file under `flaiover` selects it (`svelte-check`, I-0206), ShellCheck on the changed scripts (`shellcheck`), and the markdown lint, then, under `--all` only, the whole `flaiover` tier, which covers `flaiover-lint`, `vitest`, `flaiover-flai`, `design-routes`, and `svelte-check`, integration, which covers the Go tests and `settings-index`, and smoke. The close-out runs integration and smoke too when their `when` matches the story, as it does for a story that touches `flai/` or `template/`, and the merge queue's batch runs them again ([ADR-0233](../adrs/0233-while-the-merge-queue-can-batch-the-close-out-runs-a-queue-tier-whose-when.md)). The `gateway` tier, `scripts/gateway-smoke.sh`, the paid gateway checks, is manual: it runs only by name, `flai test gateway`, or on the host's schedule, never for paths, under `--all`, in a close-out, or in the merge queue's batch, and the smoke tier no longer runs it (S-0548). A run spends about 0.09 USD at OpenRouter, about 0.32 USD a check in Claude Code's count, which `FLAI_TEST_GATEWAY_BUDGET` caps; how to run and schedule it is in `docs/operators/index.md`. Use it between tasks; the whole run stays `make flai-test` and the close-out. The `--all` tiers select the scripts they run too, so `flai verify` and the merge queue's batch run them for a story that changes only those scripts: `template` its `template-test.sh`, `flaiover` the `flaiover-*.sh` scripts and `env.sh`, `integration` its `integration.sh`, and `smoke` every script (I-0199).
- system-flow itself runs and tests with no cluster: flaiover is one Docker container run by `flai dashboard`. flai drives a project's kluster, from the `kluster:` block of that project's `system-flow.yaml` and behind the `kluster` host action, off by default ([design/system/kluster.md](../system/kluster.md), [ADR-0254](../adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md)).
- The dashboard does not exist yet (E-0003); until then `flai board` and `flai stats` are the views.
