---
id: ADR-0303
title: "An acceptance's verify_behind_main sends back only the tiers selected by the main paths that overlap the story's change set: a path the story changes, its Go package, a Go file of main's that imports a package the story changes, or a package a Go file the story changes imports"
status: accepted
date: 2026-11-15
supersedes: []
superseded_by: []
refines: [ADR-0291]
---

# ADR-0303 An acceptance's verify_behind_main sends back only the tiers selected by the main paths that overlap the story's change set: a path the story changes, its Go package, a Go file of main's that imports a package the story changes, or a package a Go file the story changes imports

## Context

[ADR-0291](0291-every-acceptance-but-the-merge-queue-s-landing-waits-for-a-flai-verify-over.md) holds an acceptance while the paths main brought in since the story's last passing `flai verify` select a tier that verify ran. Every path main brought in counts. In this repository the `integration` and `smoke` tiers select all of `flai/**`. So every acceptance that changes flai sends back every other story in review whose verify ran them. Each one then waits for a sync and a close-out of about twelve minutes, though nothing it changes meets what main changed. A digest fix under `design/` sent S-0738 back a second time ([I-0378](../issues/I-0378-an-acceptance-sends-a-story-back-for-a-verify-that-missed-one-main-commit-landed-while-its-close-out-ran-though-the-close-out-passed-against-main.md), TH-1213, TH-1215, TH-1216).

S-0792 asks for a narrower gate. A story is sent back only when what main brings in overlaps its change set, so that its behaviour is undoubtedly affected.

The case ADR-0291 was written for must stay blocked ([I-0229](../issues/I-0229-two-stories-accepted-apart-leave-main-red-s-0535-reworded-to-backlog-s-condition-and-s-0416-s-testadefinedroleisheldtoitsfile-still-expects-the-old-words.md)). Main reworded `flai/internal/roles`, and S-0416 had changed `flai/internal/guard/roles.go`, which imports `flai/internal/roles`.

## Decision

`verify_behind_main` judges only the main paths that overlap the story's change set.

1. **The change set.** The paths the story's last passing verify recorded, its report's `Paths`: the files the branch changed against its base.
2. **Main's paths.** As ADR-0291 §2 lists them, from the report's merge base to the main branch's head, with the `wip` paths the branch does not change left out.
3. **Overlap.** A path of main's overlaps the change set when any one of these holds:
   - it is a path the change set holds;
   - it is in the same Go package as a path the change set holds;
   - it is a Go file that, at main's head, imports the package of a path the change set holds;
   - a Go file the change set holds, as the story's branch has it, imports the package that main's path is in.
4. **A path's Go package.** A path's package is the nearest directory, from the path's own up to its module's root (the directory of the nearest `go.mod`), that holds a `.go` file. So an embedded file or a `testdata` file belongs to the package that reads it. A path outside every Go module, or with no such directory, has no package and overlaps by being the same path alone. Imports are read per file, with test files counted. A file that does not parse, or cannot be read, imports nothing.
5. **The tiers sent back.** The tiers the overlapping paths select, chosen as ADR-0291 §3 chooses them, are sent back. Paths that do not overlap send nothing back. When no path overlaps, there is no blocker.
6. **The blocker's words.** The blocker names the tiers, how many of main's paths overlap out of how many main brought in, and the first overlapping paths.
7. **Unchanged.** Everything else in ADR-0291 stands: which acceptances are gated, the blocker's name, no override, and an unreadable report or main as a blocker. So does ADR-0156's reuse in the close-out a send-back asks for.

## Consequences

- An acceptance that changes another part of flai no longer sends a story in review back, nor one that changes only documents, the digest, or another sub-project. The story is rebased at its acceptance and lands on its last verify.
- I-0229 stays blocked: `flai/internal/guard/roles.go` imports `flai/internal/roles`.
- A main file that calls a function the story renames imports the story's package, so it overlaps.
- What the rule does not catch reaches main untested together: a change carried through a chain of packages, a symbol shared by files of different packages without an import between them, a change to a test tier's own scripts or the manifest, or code in another language that imports by path. The next close-out of any story whose tiers run over it, the merge queue's batch, and CI find it there. The operator chose that risk over the send-backs, in S-0792's second criterion.
- Two stories that change one Go package still send each other back, since files of one package share their identifiers and their tests.

## Alternatives considered

- **The package import closure, counted both ways, as TH-1217 planned.** Rejected: every story in review on 2026-11-15 changed a package that imports `flai/internal/workitem`, and `flai/cmd` imports nearly every internal package. So the closure, even one step deep at the package level, still sent almost every flai story back on any flai acceptance, which criterion 1 rules out.
- **File overlap alone.** Rejected: it would have let I-0229 through, since S-0416 and S-0535 changed different files.
- **Run the overlapping tiers inside the acceptance.** Rejected, as ADR-0291 rejected it: the tiers take past the dashboard's timeout.
