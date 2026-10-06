# analysis

What the analyzer found: one report per run, named for the day it ran and its focus, `<date>-<focus>.md`, the date as `YYYY-MM-DD` and the focus `bottlenecks`, `intent`, or `risk`, or `all` for a run asked for none. A report's front matter has `title`, `updated`, `status` (`draft` while the analyzer writes it, `active` once the run has ended, and `deprecated` when a later report replaces it), `focus`, and the window its metrics cover, `from` and `to`, as dates. It has one section per finding, with its evidence, its severity, and its estimated impact.

Only the analyzer writes here, and it adds a row below for each report it writes. A finding the operator acts on becomes a story or an issue; this folder records what was found, and `../system/` says how the system is. `flai check` validates every report, and warns of one with no row here and of a row with no report. Convention: `../conventions/strategic-agents.md`.

| Report | Focus | Window | Status |
|--------|-------|--------|--------|
