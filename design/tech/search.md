---
title: Search
updated: 2026-09-20
status: active
---

# Search

| | |
|-|-|
| Library | None: `flai/internal/search`, about 250 lines of Go |
| Where | In `flai serve` on the host, one index per project, answering `search.query` over the channel (S-0074). Until 0.19 it was MiniSearch 7.2 in the dashboard's process |
| Scope | `design/` and `wip/` including `archive/`, `docs/` opt-in via a toggle |

Fields: `itemId`, `title`, `tags`, `headings`, `body`, boosted 4, 3, 1, 2, 1. Terms are lower-cased runs of letters and digits. A query term matches itself, the terms it is a prefix of (weight 0.375, less the longer the completion; a single character is not followed), and the terms within a fifth of its length in edits (weight 0.45, less per edit, at most six): MiniSearch's defaults, kept so that results stay familiar. Ranking is BM25+ per field (k 1.2, b 0.7, d 0.5), terms combined with OR. A hit carries the path, kind, item ID, title, scope, status, type, nature, score, and a snippet of 160 characters around the first matched term; the dashboard adds its route.

The index is rebuilt when the files change. Each query first walks the three folders with `stat` only and hashes path, modification time, and size; an unchanged signature reuses the index. So search needs no word from the watcher, and an index can never be staler than the files. On this repository (544 documents) a rebuild takes about 60 ms and the walk a few.

**Why no library.** The dashboard no longer reads project files (ADR-0029), so the index had to move to flai. Bleve is a search engine with its own storage and many times the size of the problem; a repository of a few thousand Markdown files indexes in well under a second in memory, and what MiniSearch did for us fits in one file that is easy to test against.

**Compared with MiniSearch** on this repository on 2026-09-20, ten queries against the running 0.19.2 dashboard and `flai hostapi search.query`: the first hit was the same for nine, and 41 of the 50 top-five results were shared, in nearly the same order. The one that differed was `S-0052`: MiniSearch followed the prefix `s` into every word that starts with it and put three design documents first, with the story itself not in the top five; flai does not follow a one-character prefix and puts S-0052 first. Differences further down are ties broken differently (flai breaks them by path).
