---
title: Documentation
updated: 2026-09-15
audience: agent
order: 50
status: active
---

# Documentation

Where documents live, what every file carries, how it is written, and how it stays true. The full standard is `design/system/documentation-standard.md`.

## Rules

- Everything is markdown with YAML front matter. Diagrams are Mermaid. No binary images under `design/` or `wip/`.
- Put each document where its audience looks: `design/system` for how the system is, `design/adrs` for why, `design/tech` for what it is built with, `design/conventions` for how we work, `docs/<audience>` for outward-facing guides, `wip/` for work state.
- Prefer editing the living document over adding a new one. A new file needs a reason a reader would agree with.
- Front matter minimum for design and docs files: `title`, `updated` (ISO date, bumped on every meaningful edit), `status`. `README.md` files are indexes and are exempt.
- Titles containing `: ` or starting with a YAML-special character are double-quoted.
- Dates are ISO 8601 UTC. Timestamps are `YYYY-MM-DDTHH:MM:SSZ` and are real: never write a time you did not observe.
- Files and folders are lowercase kebab-case. Work items are `<ID>-<slug>.md`. ADRs are `NNNN-slug.md`.
- Links between documents are relative paths so they work in git hosting and in the dashboard. Do not link to line numbers.
- Every folder a reader might land in has a `README.md` saying what it is for.
- When behavior changes, the documentation changes in the same commit: `docs/` for users, `design/system` for builders, the convention file for agents.
- Write plainly. Say what is, not what might be. Tables for parallel facts, prose for argument. Short sentences. No filler, no marketing.
- Do not duplicate content. Link to the source of truth instead; a copied paragraph is a future contradiction.
- Keep `design/tech` honest: add a file when a dependency is added, bump the version when it moves, remove it when it goes.
- Operator documentation covers every configuration setting, with an index so an operator can find a setting quickly whether it is an argument, an environment variable, stored configuration, or a stored secret.
- Operator documentation includes runbooks for installation, updates, deletion, backups, restores, and migrations, and says where dashboards, tools, and other operational information are.

## When in doubt

- If you cannot name the audience, you have not found the folder.
- If a document has grown past what its audience will read, split by audience, not by length.

<!-- system-flow:end-of-baseline -->

## Project additions
