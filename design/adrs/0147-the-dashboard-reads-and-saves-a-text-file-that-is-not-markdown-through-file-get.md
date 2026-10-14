---
id: ADR-0147
title: "The dashboard reads and saves a text file that is not markdown through file.get and file.save, under the design and docs folders or at a path the manifest names, by doc.save's rules for the hash, the check, the commit, and the story"
status: accepted
date: 2026-10-14
supersedes: []
superseded_by: []
refines: [ADR-0023, ADR-0145]
topics: [cli, dashboard]
---

# ADR-0147 The dashboard reads and saves a text file that is not markdown through file.get and file.save, under the design and docs folders or at a path the manifest names, by doc.save's rules for the hash, the check, the commit, and the story

## Context

[ADR-0023](0023-documents-are-saved-through-flai.md) has the dashboard save a document through `flai doc save`, checked by `flai check` and committed by flai. `doc.get` and `doc.show` refuse a file that is not `.md` (`docPath` in `flai/internal/hostapi`). [ADR-0145](0145-a-dashboard-save-may-name-a-story-and-is-then-committed-on-that-story-s-branch.md) lets a save name a story and commits it on that story's branch. [ADR-0029](0029-the-dashboard-reaches-a-project-only-through-flai-on-the-host.md) has the dashboard reach a project only through named host methods, with none that reads an arbitrary file.

[ADR-0146](0146-a-project-s-design-tokens-are-a-dtcg-file-the-manifest-names-beside-its.md) makes a project's design tokens a DTCG JSON file. On TH-0445, on 2026-10-14, the operator asked for `file.get` and `file.save` "as alternatives that we can leverage for different kinds of files", rather than forcing such files into markdown. The finding is [design-system-tool.md](../system/design-system-tool.md).

## Decision

The dashboard reads and saves a text file that is not markdown through the host methods `file.get` and `file.save`, under the manifest's design and docs folders or at a path the manifest names for a purpose, by `doc.save`'s rules for the hash, the check, the commit, and the story.

1. **The methods.** `file.get` `{path, story?}` answers the file's raw content and its SHA-256. `file.save` `{path, content, hash, story?, message?}` saves it. The CLI is `flai file show <path>` and `flai file save <path> --hash <h>`, the content on standard input, each with `--story`, as `flai doc show` and `flai doc save` are.
2. **What they reach.** A file under the manifest's design or docs folder, or a path the manifest names for a purpose: the design token file and the design stylesheet now.
3. **What they refuse.** Any other path; a path under the wip folder; a binary file; a file over a size limit S-0381 sets; and a markdown file, which keeps `doc.show` and `doc.save` and their edit modes.
4. **The rules of a save.** A stale hash is a conflict, exit 3, with the current content and a diff. flai runs `flai check` before and after writing; a finding the content introduces, or an error already on the file, puts the old file back and refuses the save, exit 4. The save commits that path on the main checkout's current branch as `docs:`, unless `dashboard.autocommit` is false. With `story`, it is read, checked, and committed in that story's worktree, on `story/S-nnnn`, and refused as ADR-0145 refuses.
5. **ADR-0029 stands.** There is still no method that reads an arbitrary file. A path outside the two folders enters `file.get`'s reach only by being named in the manifest.

## Consequences

- The Design page reads and saves the token file, on main or on a proposal's branch.
- A non-markdown file in the design or docs folder is checked only by the rules that read it, such as `design.tokens`. Others have the hash and the commit, and no validation.
- Each manifest key that names a path for a purpose widens what the dashboard can read and write. Such keys are added by a decision, not by the operator's configuration.
- A host flai older than these methods lacks them: the dashboard's route answers 501, as for `doc.lint`.

## Alternatives considered

- **Widen `doc.get` and `doc.save` to any file type.** They carry front matter, edit modes, and the markdown lint. A second pair keeps each contract plain.
- **A pair of methods per purpose, such as `design.save`.** Every new kind of file would need its own. The operator asked for a general pair.
- **A method that reads any file of the checkout.** ADR-0029 refuses it, and the token file does not need it.
