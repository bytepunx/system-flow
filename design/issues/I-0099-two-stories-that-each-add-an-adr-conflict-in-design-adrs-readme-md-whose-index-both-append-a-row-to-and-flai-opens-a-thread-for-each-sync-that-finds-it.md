---
id: I-0099
title: Two stories that each add an ADR conflict in design/adrs/README.md, whose index both append a row to, and flai opens a thread for each sync that finds it
class: efficiency
status: open
count: 1
cost: 5m
first_reported: 2026-10-06T23:31:35Z
last_reported: 2026-10-06T23:31:35Z
updated: 2026-10-06T23:31:35Z
---

# I-0099 Two stories that each add an ADR conflict in design/adrs/README.md, whose index both append a row to, and flai opens a thread for each sync that finds it

## Description

`flai adr new` adds a row for the new ADR to `design/adrs/README.md`. Two stories in progress that each add an ADR therefore change the same lines of that file, and git cannot merge them. flai's trial merge at sync reports the pair as conflicting and opens a thread to the two agents, which says to narrow or to wait; neither applies, since every row belongs. Whichever story is accepted second resolves the index by hand at its sync or its rebase, keeping both rows.

Since ADR-0096 lets stories in review stop holding others and the shared paths include `design/adrs`, more stories run at once, and most add an ADR. The thread churn is the same as the issue summary's before S-0278.

## Instances

### 2026-10-06T23:31:35Z
Story: S-0301.
On 2026-10-06 between 22:59Z and 23:27Z, three stories in progress at once (S-0299, S-0300, S-0301, then S-0261) each added an ADR. flai's trial merge found every pair conflicting on design/adrs/README.md alone, and opened TH-0220, TH-0221, TH-0222, TH-0226, and TH-0228, resolving two of them when a later sync merged cleanly and opening them again when it did not. Claude, watching the board for the operator, told each pair not to wait or narrow and to keep both rows; S-0301's agent kept both rows by hand at its last sync. The summary.md conflict (I-0074) was closed this way by S-0278, and the ADR index is the next file every story with an ADR edits.

## Remediation

The index is derived from the ADR files, as `summary.md` is from the issues: `flai adr new` could regenerate it whole, and `flai stream sync` and `flai accept` could treat it as a generated file (ADR-0098), writing it again from the ADR files when a rebase stops on it alone, and leaving it out of a pair's conflicts. `design/adrs/README.md` would then need to carry nothing but the generated table, or hold its hand-written part apart from it.
