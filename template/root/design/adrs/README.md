# Architecture Decision Records

One decision per file, numbered in order of acceptance. Never edit an accepted ADR except to set `superseded_by`, or its `topics` with `flai adr topics`. To change a decision, write a new ADR that supersedes it, then update the living design.

Record a new one with `flai adr new "<title>"`, or from the dashboard's ADRs page: the number, the file, its front matter, and the row below are made for you, from the sections in `0000-template.md`. From a story's worktree, `flai adr new --commit` also commits it on the story's branch and adds it to the story's touches. `flai check` warns when a file has no row here or a row has no file.

| ADR | Title | Status |
|-----|-------|--------|
