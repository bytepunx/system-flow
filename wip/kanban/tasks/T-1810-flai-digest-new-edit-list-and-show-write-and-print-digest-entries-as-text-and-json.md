---
id: T-1810
type: task
nature: feature
title: flai digest new, edit, list, and show write and print digest entries, as text and JSON
status: backlog
parent: S-0428
owner: alex
created: 2026-10-16T06:03:22Z
updated: 2026-10-16T06:03:22Z
transitions: []
stream: S-0428
tags: [cli]
touches: [flai/cmd/digest.go, flai/cmd/digest_test.go, flai/cmd/root.go, docs/users/flai-reference.md]
after: [T-1808]
---
# T-1810 flai digest new, edit, list, and show write and print digest entries, as text and JSON

## Work

Add `flai/cmd/digest.go` with the `digest` command group, registered in `flai/cmd/root.go`:

- `flai digest new "<fact>" --kind <kind> --source <id>... [--topic <t>]... [--path <p>]... [--body <text> | --body-file <file>] [--json]`, through `digest.New`.
- `flai digest edit <D-nnnn> [--body] [--topic] [--path] [--source] [--asked-again <TH-nnnn>]... [--json]`, through `digest.Edit`.
- `flai digest list [--status <s>] [--topic <t>] [--json]` and `flai digest show <D-nnnn> [--json]`.

A refusal prints what the check found and exits non-zero, writing nothing. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.

Waits for T-1808. Runs beside the check task and the search task, which touch other files.

## Done when

`flai/cmd/digest_test.go` covers each subcommand in text and JSON, a refused `new` that leaves no file, and an unknown ID; the reference lists the four subcommands; `flai test flai/cmd` passes.

## Notes

The curator's guard (S-0431) passes these commands by name, so keep them under `digest`.
