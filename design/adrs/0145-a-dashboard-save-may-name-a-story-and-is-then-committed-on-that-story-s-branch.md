---
id: ADR-0145
title: "A dashboard save may name a story and is then committed on that story's branch in its worktree, the editor lints with flai's subset of the markdown rules as the operator types, and an acceptance refused for unfixed markdown gives its findings as data"
status: accepted
date: 2026-10-13
supersedes: []
superseded_by: []
refines: [ADR-0023]
topics: [cli, dashboard]
---

# ADR-0145 A dashboard save may name a story and is then committed on that story's branch in its worktree, the editor lints with flai's subset of the markdown rules as the operator types, and an acceptance refused for unfixed markdown gives its findings as data

## Context

[ADR-0023](0023-documents-are-saved-through-flai.md) has the dashboard save a document through `flai doc save`, checked by `flai check` and committed on the main checkout's current branch. [ADR-0140](0140-flai-accept-refuses-an-acceptance-whose-files-a-tier-with-a-fix-command-still.md) refuses an acceptance when a tier with a `fix_command` still fails after the fix, and names each file, line, and error. The operator then had to correct the files by hand on the host: the files the branch changes are on `story/S-nnnn` in the story's worktree, not on the main checkout, and the refusal reached the dashboard only as text (S-0378).

The operator wants to fix that markdown from the dashboard: open the file a report, a refusal, a thread, or a page names, at its line, see the lint while typing, save, and accept again. The orchestrator recommended this decision on TH-0436. The operator's confirmation there was still pending when this was written.

## Decision

A document saved from the dashboard may name a story, and is then saved and committed on that story's branch in its worktree; the editor shows flai's subset of the markdown rules as the operator types; and `flai accept --json` gives a refusal for unfixed markdown as data the dashboard opens file by file.

1. **A save may name a story.** `flai doc show <path> --story S-nnnn` and `flai doc save <path> --story S-nnnn --hash <h>` read and save the file in that story's worktree. The hash, the edit mode, and `flai check` are the branch's, as ADR-0023 sets them for the main checkout. The save commits that one path on `story/S-nnnn`, with the subject `docs: [S-nnnn] <message>`. The main checkout is left as it was. The host methods `doc.show` and `doc.save` take an optional `story` and run in that worktree.
2. **What a story's save refuses.** A path under the wip folder, which lives on the main checkout ([ADR-0019](0019-story-branches-and-touches.md)); a story without a worktree; a worktree not on `story/S-nnnn`, as mid-rebase or detached; an item that is not a story. Each is a plain error, exit 1. A stale hash is still a conflict, exit 3, and refused content is still exit 4.
3. **The editor lints as the operator types.** `flai doc lint <path> --stdin` lints the content given as the markdown file at path, with the project's markdownlint configuration and flai's subset of its rules ([ADR-0061](0061-flai-lints-the-markdown-it-writes-in-wip-with-its-own-implementation-of-the.md)). The host read method `doc.lint` `{path, content}` answers the same in process: `{path, findings: [{line, rule, alias, message}]}`. This is advice. The project's markdown tier stays the authority at acceptance.
4. **A refusal for unfixed markdown is data.** `flai accept S-nnnn --json`, refused under ADR-0140, prints `{"refused": {"id", "reason", "findings"}}` and exits 4. Each finding has `tier`, `checkout` (`main`, or `story` with `story`), `path`, `message`, and, when markdownlint gives them, `line`, `column`, `level`, `rule`, and `alias`. Without `--json` the text refusal and its exit 1 are unchanged. `accept.run` maps exit 4 to a Refused error with `{id, reason, findings}` as its data, which the dashboard passes on, so the review page opens each file at its line, on the story's branch for a story finding.

## Consequences

- The operator fixes what blocks an acceptance in the dashboard, and accepts again without reloading.
- A save on a story's branch is a commit on that branch, not on main. The story's agent meets it at its next sync, as it meets any commit on its branch.
- The live findings can differ from the tier's: flai implements a subset of the rules. The editor says so, and the acceptance still decides.
- `doc.lint` runs on every pause in typing. It runs in process, so no process starts for it.
- A host flai older than `doc.lint` lacks the method: the dashboard's route answers 501, and the editor shows no findings.
- Work items and narratives, under the wip folder, still save on the main checkout only.

## Alternatives considered

- **Save on the main checkout and let a sync carry it to the branch.** The refused files are the branch's. A save on main would not change what the acceptance checks until a sync, and could conflict with the branch.
- **Parse the text refusal in the dashboard.** The wording would become a contract. JSON data keeps the fields stable.
- **Run the markdown tier on each pause in typing.** It starts a process, `npx` for `scripts/lint-md.sh`, and takes seconds. flai's subset answers at once, and the tier still runs at acceptance.
