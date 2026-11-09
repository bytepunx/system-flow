---
id: T-2737
type: task
nature: remediation
title: "An ADR and metrics.md § Close-outs say a serve journal line that is not an entry is skipped and warned of, and the section stays"
status: done
parent: S-0624
owner: alex
created: 2026-11-09T17:04:34Z
updated: 2026-11-09T17:12:14Z
transitions:
  - to: ready
    at: 2026-11-09T17:11:17Z
    by: agent-S-0624
  - to: in-progress
    at: 2026-11-09T17:11:19Z
    by: agent-S-0624
  - to: done
    at: 2026-11-09T17:12:14Z
    by: agent-S-0624
stream: S-0624
tags: [docs, metrics]
touches: [design/system/metrics.md, design/adrs/README.md, docs/users/flai.md]
usage:
  source: log
  seconds: 55
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 15
      output: 2085
      cache_read: 1800279
      cache_write: 29084
      cost: 0.6345
---
# T-2737 An ADR and metrics.md § Close-outs say a serve journal line that is not an entry is skipped and warned of, and the section stays

## Work

Record the decision before the code, since `design/system/metrics.md` is the contract between `flai stats` and the dashboard and changes only with an ADR.

- Write an ADR with `flai adr new`: the close-outs read every line of the serve journal. A line has its NUL bytes removed before it is decoded, and a line that still is not an entry is skipped and counted. `flai stats` warns once on stderr, naming the lines skipped and the first one's number, and keeps the section. The section is absent only when the journal cannot be opened or read. Give its context from I-0290 and the analysis in `design/analysis/2026-11-07-all.md` § Finding 8. Weigh the alternatives: stopping at the line as now, skipping without a warning as `flai serve journal` does, and a repair command that rewrites the host's journal. Accept it and let flai regenerate `design/adrs/README.md`.
- In `design/system/metrics.md` § Close-outs (S-0546), rewrite the `closeouts` row to say this, linking the ADR, and bump `updated`.
- In `docs/users/flai.md`, rewrite the sentence of the close-outs paragraph that says a journal that cannot be read leaves the section absent.

Waits for nothing.

## Done when

- The ADR is accepted, and `design/adrs/README.md` lists it.
- `metrics.md` § Close-outs and `docs/users/flai.md` say a damaged line is skipped and warned of, and the section is absent only for a journal that cannot be read.
- `flai check --strict` and the markdown lint pass on the changed files.

## Notes
