---
id: T-2318
type: task
nature: remediation
title: "flai adr new writes <, >, and & in an ADR's title as themselves, and the title flai prints is the title"
status: done
parent: S-0538
owner: alex
created: 2026-10-29T02:57:13Z
updated: 2026-10-29T03:09:10Z
transitions:
  - to: ready
    at: 2026-10-29T02:59:47Z
    by: system-flow
  - to: in-progress
    at: 2026-10-29T02:59:49Z
    by: system-flow
  - to: done
    at: 2026-10-29T03:09:10Z
    by: agent-S-0538
stream: S-0538
tags: [flai]
touches: [flai/internal/adr/adr.go, flai/internal/adr/adr_test.go, flai/cmd/adr_test.go]
usage:
  source: log
  seconds: 561
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 2260
      cache_read: 830139
      cache_write: 22085
      cost: 0.388
---
# T-2318 flai adr new writes <, >, and & in an ADR's title as themselves, and the title flai prints is the title

## Work

- In `flai/internal/adr/adr.go`, make `yamlString` quote with a `json.Encoder` whose `SetEscapeHTML(false)` is set, trimming the newline it adds, so a quoted title keeps `<`, `>`, and `&` as themselves rather than `\u003c`, `\u003e`, and `\u0026`.
- Make `titleOf` read a double-quoted title as YAML does (an escaped quote or backslash comes out as the character), so the title `flai adr new`, `accept`, and `topics` print is the title, not its quoted form.
- Add `flai/internal/adr/adr_test.go`: `yamlString` on a plain title, on titles with `<name>-<item>`, `&`, `:`, a double quote, and a backslash; each result parses with goccy/go-yaml back to the title, and none holds `\u00`.
- Reproduce I-0222 in `flai/cmd/adr_test.go`: `flai adr new` with a title holding `<name>-<item> & co`, then `flai adr topics` on it; the file's front matter and its H1 hold the angle brackets and the ampersand as written, the printed title matches, and `flai check` reports nothing on the ADR or the index.

Waits for nothing: it is the first layer.

## Done when

- The new tests fail on the code before the change and pass after it.
- `flai test flai/internal/adr flai/cmd` passes.

## Notes

The cause is in S-0538's Notes under Planning: `json.Marshal` escapes HTML characters by default, and `flai adr topics` only showed what `flai adr new` wrote.
