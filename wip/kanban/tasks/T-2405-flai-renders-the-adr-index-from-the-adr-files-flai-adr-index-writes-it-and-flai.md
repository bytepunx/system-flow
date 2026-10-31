---
id: T-2405
type: task
nature: improvement
title: flai renders the ADR index from the ADR files, flai adr index writes it, and flai adr new, accept, and a supersede regenerate it
status: backlog
parent: S-0550
owner: alex
created: 2026-10-31T07:44:09Z
updated: 2026-10-31T07:44:09Z
transitions: []
stream: S-0550
tags: [cli, template]
touches: [flai/internal/adrindex/adrindex.go, flai/internal/adrindex/adrindex_test.go, flai/internal/adr/adr.go, flai/internal/adr/adr_test.go, flai/cmd/adr.go, flai/cmd/adr_test.go, design/adrs/README.md, template/root/design/adrs/README.md]
after: [T-2390]
---
# T-2405 flai renders the ADR index from the ADR files, flai adr index writes it, and flai adr new, accept, and a supersede regenerate it

## Work

Today the index's rows are edited one at a time. `adr.New` adds one with `addRow`, `adr.Accept` changes its status cell with `setRowStatus`, and a supersede notes it with `noteRow` (all in `flai/internal/adr/adr.go`). The MCP tool `adr_new` goes through `adr.New`. Nothing rebuilds the table, and there is no `flai adr index`.

- In `flai/internal/adrindex/adrindex.go`, add a render. It reads every ADR file under the ADR folder and writes the table, one row per file in number order, with the link, the title, and the status cell the story's ADR decides. It keeps the text above and below the table as it is. `Path` and the index's split stay.
- Make `adr.New`, `adr.Accept`, and the supersede path regenerate the whole index in place of `addRow`, `setRowStatus`, and `noteRow`, and remove those helpers if nothing else uses them.
- Add `flai adr index` to `flai/cmd/adr.go`, which writes the index from the folder and says whether it changed. Its help text feeds `docs/users/flai-reference.md`, which the documentation task regenerates.
- Regenerate this repository's `design/adrs/README.md` and read the diff row by row. Each status cell that changes carries a hand note the front matter does not hold. Record each in `## Decisions`, with what the ADR decided for it. Change an accepted ADR's front matter only where flai allows (`superseded_by`, `topics`), and never by hand.
- Rewrite the intro of `design/adrs/README.md` and `template/root/design/adrs/README.md` to say the table is generated from the ADR files by `flai adr index` and is never edited by hand.
- Test the render, including the "in part" and "refines" cells, in `adrindex_test.go`. Test the write paths in `flai/internal/adr/adr_test.go` and the command in `flai/cmd/adr_test.go`.

Second layer: it waits for the story's ADR task, whose decision sets the status cell and which writes `design/adrs/README.md` too. The ADR index check task waits for this one, for the render.

## Done when

`flai adr index` on this repository writes an index whose every row matches its ADR file, and a second run changes nothing. `flai adr new`, `flai adr accept`, and a supersede leave the index equal to what `flai adr index` writes. The tests pass with `go test -race ./internal/adrindex ./internal/adr ./cmd`. Story criterion 1, for the ADR index and its commands.

## Notes

Drafted by the planner of S-0550 on 2026-10-31.
