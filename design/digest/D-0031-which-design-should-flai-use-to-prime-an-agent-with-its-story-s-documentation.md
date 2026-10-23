---
id: D-0031
title: "Which design should flai use to prime an agent with its story's documentation: Could we have conventions optionally include tags in the front matter and then have headings that tagged the following…"
kind: answer
sources: [TH-0020]
topics: [flai, cli, go, code, all]
paths: [flai/cmd]
date: 2026-09-26
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0031 Which design should flai use to prime an agent with its story's documentation: Could we have conventions optionally include tags in the front matter and then have headings that tagged the following…

The operator, alex, answered TH-0020 at 2026-09-26T17:20:34Z.

What was asked:

```text
The findings are in `design/system/agent-context.md` on `story/S-0125`. Today every agent reads all 12 conventions (47 KB) and the whole issues table, then has to find the relevant ADRs and design itself.
```

The answer:

```text
Could we have conventions optionally include tags in the front matter and then have headings that tagged the following bullets as mapping to those tags? flai prime could then determine which conventions were relevant vs. including them all for every story.
```

Source: TH-0020.
