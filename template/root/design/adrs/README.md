# Architecture Decision Records

One decision per file, numbered in order of acceptance. Never edit an accepted ADR except to set `superseded_by`, or its `topics` with `flai adr topics`. To change a decision, write a new ADR that supersedes it, then update the living design.

Record a new one with `flai adr new "<title>"`, or from the dashboard's ADRs page: the number, the file, and its front matter are made for you, from the sections in `0000-template.md`. From a story's worktree, `flai adr new --commit` also commits it on the story's branch and adds it to the story's touches.

The table below is generated from the ADR files by `flai adr index`, and written the same way by `flai adr new` and `flai adr accept`. Each row is one file's link, title, and status cell: its `status`, then the ADRs it supersedes, refines, and is superseded by. How much of an older ADR a newer one replaces is said in the newer one. Never edit the table by hand: `flai check` reports a table that differs from what the files render, and the next regeneration undoes the edit.

| ADR | Title | Status |
|-----|-------|--------|
