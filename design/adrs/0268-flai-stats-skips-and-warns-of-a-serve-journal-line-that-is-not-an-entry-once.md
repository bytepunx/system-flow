---
id: ADR-0268
title: "flai stats skips and warns of a serve journal line that is not an entry once its NUL bytes are removed, and keeps the close-outs section, which is absent only when the journal cannot be opened or read"
status: accepted
date: 2026-11-09
supersedes: []
superseded_by: []
refines: [ADR-0229]
---

# ADR-0268 flai stats skips and warns of a serve journal line that is not an entry once its NUL bytes are removed, and keeps the close-outs section, which is absent only when the journal cannot be opened or read

## Context

`flai stats` counts close-outs from flai serve's host journal, `~/.flai/serve/journal.jsonl` ([ADR-0229](0229-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-asks-for-with-the.md) §7). `statsread.journalEntries` decoded every line and stopped at the first that was not an entry, and `metrics.md` § Close-outs said an unreadable journal leaves the section absent. So one damaged line dropped every week of close-outs.

That happened on this host ([I-0290](../issues/I-0290-one-line-holding-a-nul-byte-in-flai-serve-s-journal-drops-the-whole-close-outs.md)). Line 874 of the journal is an `orchestrate` entry of 2026-10-08 preceded by NUL bytes: a write the host's stop cut short left the file's tail zero-filled, and the next append landed after it. flai 1.84.0's `flai stats --json` warned `line 874: invalid character '\x00' looking for beginning of value` and left the section out, so the journal's 166 `verify.job` lines were unreported.

The journal is appended to by serve, one whole line per write, and nothing rewrites it. A torn write after a crash or a stop is the filesystem's doing, which no writer change prevents. `flai serve journal` already skips a line it cannot decode.

## Decision

flai stats skips and warns of a serve journal line that is not an entry once its NUL bytes are removed, and keeps the close-outs section, which is absent only when the journal cannot be opened or read.

1. **NUL bytes.** Each line of the journal has its NUL bytes removed, then its surrounding space, before it is decoded. A line that is NUL bytes alone is an empty line, skipped with no count. A line with NUL bytes before a whole entry, as I-0290's line 874, is that entry.
2. **A line that is not an entry.** A line that still does not decode as a journal entry is skipped and counted. The entries around it are read as before.
3. **The warning.** When one or more lines were skipped, `flai stats` warns once on stderr, at `warn`, with the count and the number of the first skipped line, and its decode error, and keeps the section.
4. **Absent.** The section is absent when the host has no journal, with no warning, and when the journal cannot be opened or read, as a line longer than the reader's limit, with a warning, as before.
5. **No new field.** The JSON shape of `closeouts` does not change. The skipped lines are in the warning only, so the dashboard has nothing to change.

## Consequences

- One damaged line no longer hides a month of close-out verdicts and durations; the section counts every entry it can read.
- A close-out whose end entry is on a damaged line is not counted, and the warning says how many lines were skipped, so the gap is visible.
- `flai stats` and `flai serve journal` now agree on what a damaged line is: neither stops on it.
- `metrics.md` § Close-outs and `docs/users/flai.md` say so.

## Alternatives considered

- **Repair the journal by hand.** It clears this host once, and the next torn write drops the section again.
- **Rewrite the journal on serve's start, dropping damaged lines.** It changes the operator's record of host actions, and the reader would still stop on a line torn after the start.
- **A `skipped` field under `closeouts`.** The dashboard would have to show it for it to matter, and it is a property of the host's file, not of the flow; the warning names it where the operator reads `flai stats`.
- **Stop at the first damaged line but keep the entries before it.** It still drops every close-out after the damage, which on this host is all of them.
