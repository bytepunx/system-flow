---
id: ADR-0318
title: "flai digest new and edit refuse an entry whose markdown the project's lint rejects, as flai's wip writes do"
status: accepted
date: 2026-11-17
supersedes: []
superseded_by: []
refines: [ADR-0061, ADR-0191]
---

# ADR-0318 flai digest new and edit refuse an entry whose markdown the project's lint rejects, as flai's wip writes do

## Context

[ADR-0061](0061-flai-lints-the-markdown-it-writes-in-wip-with-its-own-implementation-of-the.md) has flai lint the markdown it writes in `wip/`, because that folder is written in the main checkout, where no story's lint runs. The digest ([ADR-0191](0191-the-digest-is-one-markdown-file-per-entry-under-the-design-folder-s-digest.md)) is written the same way: the curator runs `flai digest new` and `edit` in the main checkout. Those commands checked an entry's front matter, sources, and secrets, never its markdown, and `flai check` lints only `wip/`.

[I-0363](../issues/I-0363-flai-s-markdown-lint-reports-md038-on-a-code-span-holding-escaped-backticks.md) is what followed. The curator wrote D-0592 with a code span that markdownlint-cli2 0.20.0 rejects (MD038 at line 17). It reached main unlinted, and `TestRepositoryLintsClean`, which assumes every file in the monorepo passes markdownlint, read mdlint's finding as a false positive. S-0836 ran both tools on the line: they agree.

## Decision

`digest.New` and `digest.Edit` lint the entry they wrote with the project's markdownlint configuration, through `workitem.Repo.LintGuard`, as flai's wip writers do. Each finding the write brings is a refusal problem with rule `markdown.MDnnn`, its line, and markdownlint's message. `new` removes the file and any folder it made; `edit` puts the file back, and a finding the entry had before the edit does not refuse it. A project with no markdownlint configuration is not linted. A configuration that cannot be read refuses the write, naming the error.

## Consequences

- An entry the lint rejects no longer reaches main through flai digest. The curator, or `flai digest seed`, which writes through `new`, is told the rule and the line and writes the entry again.
- The digest folder is still not in `flai check`'s lint. A hand edit there is caught by the smoke tier's `scripts/lint-md.sh` and by `TestRepositoryLintsClean`.

## Alternatives considered

- **Lint the digest folder in `flai check`.** Lost: it reports a bad entry after it is on main; refusing the write keeps it off.
- **Fix mdlint's parser.** Lost: mdlint was right. Both tools close a code span at the backtick after a backslash, as CommonMark says.
