---
id: D-0318
title: "E-0026 needs five design answers: recommend the project's own Vite dev server on a loopback port, framed cross-origin; Svelte 5 first; examples beside each component; a components host action…"
kind: answer
sources: [TH-0483]
topics: [all]
paths: []
date: 2026-10-16
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0318 E-0026 needs five design answers: recommend the project's own Vite dev server on a loopback port, framed cross-origin; Svelte 5 first; examples beside each component; a components host action…

The operator, alex, answered TH-0483 at 2026-10-16T21:29:20Z.

What was asked:

```text
E-0026 asks for the project's own components, rendered live. `design/system/design-system-tool.md` § The project's own components left exactly this out of E-0020: it needs the project's own Svelte and Tailwind build, which only the host has, and ADR-0029 and ADR-0031 keep the dashboard's container away from the project. I am drafting the stories on the recommended answers below. The cost of delay input is already set (600 USD a week on the epic), so nothing else waits. Answer "all as recommended", or give the numbers you want otherwise.
```

The answer:

```text
1. perfect 2. yes, Svelte 5 and Vite first 3. let's start with that 4. lets put supporting files under design/components, 5. lets gate it behind the feature for now
```

Source: TH-0483.
