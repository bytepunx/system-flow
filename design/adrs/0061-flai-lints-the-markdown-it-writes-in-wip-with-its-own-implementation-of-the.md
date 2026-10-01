---
id: ADR-0061
title: "flai lints the markdown it writes in wip with its own implementation of the markdownlint rules it can break, configured by the project's markdownlint file"
status: accepted
date: 2026-10-01
supersedes: []
superseded_by: []
topics: [cli]
---

# ADR-0061 flai lints the markdown it writes in wip with its own implementation of the markdownlint rules it can break, configured by the project's markdownlint file

## Context

Work items, threads, and narratives are written by flai in the main checkout (ADR-0019), so the markdown lint a story runs in its worktree never sees them, and they reach `main` unlinted and turn CI red: I-0027 counts ten occurrences (MD004, MD009, MD012, MD024, MD026, MD029, MD036, MD037), each costing the next stories' smoke runs. S-0179 asked flai to lint what it writes, leaving open whether it embeds a linter or runs the project's own (`scripts/lint-md.sh`, markdownlint-cli2), and required that a missing linter not fail `flai check`. The lint has to run where flai writes: `flai check`, which `item_new`'s and `item_edit`'s refusal already runs twice per call, and every thread reply and narrative log entry.

## Decision

flai lints markdown with its own Go implementation, `flai/internal/mdlint`, of the markdownlint rules what flai and its agents write can break (MD001, MD004, MD009, MD010, MD012, MD018, MD019, MD022, MD024, MD025, MD026, MD029, MD031, MD034, MD036, MD037, MD040, MD047, MD049, MD050), configured by the project's own markdownlint file at its root, with markdownlint's options, aliases, tags, and inline comments. A project with no such file is not linted. `flai check` reports each finding under the wip folder as a warning, `markdown.MDnnn`, so `--strict` fails on it; item creation, thread entries, and narrative log entries refuse what would bring a finding, naming the rule and line; titles lose the punctuation MD026 rejects.

## Consequences

- The lint runs in-process in milliseconds over the whole wip folder, with no Node and no network, in any project made from the template, and the same in the MCP server, the CLI, and the dashboard's calls.
- markdownlint-cli2 stays the authority. A rule mdlint does not implement, or a case it cannot place, is not reported, and CI catches it as before; mdlint is built to report nothing rather than something markdownlint would not.
- Agreement with markdownlint is tested: `scripts/mdlint-fixtures.sh` records what markdownlint-cli2 reports on mdlint's fixtures, the test compares rule and line for every implemented rule, and the integration tier lints every markdown file in this repository and expects nothing. A new markdownlint version, or a rule agents start to break, means a fixture, a regenerated reference, and possibly a rule.
- A rule's options are read from the file; options it does not understand fall back to defaults, and `MD004` style `sublist` and `MD009` `strict` are not judged at all.

## Alternatives considered

- Run markdownlint-cli2 from `flai check`: needs Node and, through npx, the network on first use; a process per check, run twice per edit, is seconds, not milliseconds; a template project may have neither.
- Lint only in `scripts/lint-md.sh` before review: the files are in the main checkout and written after the story's last commit, by the agent and by the operator; nothing would run it there.
- A Go markdown parser library with a lint layer: none implements markdownlint's rules, so the rules would be written anyway, over a parser that is not CommonMark-and-GFM as micromark is, and a new dependency.
