---
id: T-2433
type: task
nature: improvement
title: The board's Import prompt says the import is committed before its tests and names the story that runs them
status: done
parent: S-0424
owner: alex
created: 2026-11-02T06:13:03Z
updated: 2026-11-02T06:20:12Z
transitions:
  - to: ready
    at: 2026-11-02T06:13:33Z
    by: agent-S-0424
  - to: in-progress
    at: 2026-11-02T06:13:37Z
    by: agent-S-0424
  - to: done
    at: 2026-11-02T06:20:12Z
    by: agent-S-0424
stream: S-0424
tags: []
touches: [flaiover/src/lib/components/ImportPrompt.svelte, flaiover/src/lib/components/ImportPrompt.svelte.test.ts, flaiover/src/routes/api/import/+server.ts]
after: [T-1785]
usage:
  source: log
  seconds: 395
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 18
      output: 5456
      cache_read: 1183732
      cache_write: 33943
      cost: 0.5382
---
# T-2433 The board's Import prompt says the import is committed before its tests and names the story that runs them

## Work

`flaiover/src/lib/components/ImportPrompt.svelte` still says the board's Import runs the tests and commits only when they pass, and that a failed test leaves the import uncommitted. Since T-1785, `import.run` runs `flai import --formula`, which commits the import first and starts the `import` formula's instance, whose steps run `flai check --strict` and the tests.

- The prompt's list of what Import does says it commits the import, then starts a story in the new project that checks it and runs the tests it names.
- When it is done, the page names that story from `result.formula.story` beside the commit.
- An import committed whose instance was not made arrives as an `error` event with flai's message, which says how to make it; the page shows it as any error.
- The `committed: false` branch stays for an older host that still runs `--commit`.
- The comments in `flaiover/src/routes/api/import/+server.ts` describe `--formula`.

Second layer: waits for T-1785, which switches the host API to `--formula`. Shares no path with T-1786.

## Done when

`ImportPrompt.svelte.test.ts` passes with cases for the new wording and for the instance named, and `flai test` on the changed files passes.

## Notes

Added by agent-S-0424 on 2026-11-02: T-1785's review found the board's prompt describing the old `--commit` flow.
