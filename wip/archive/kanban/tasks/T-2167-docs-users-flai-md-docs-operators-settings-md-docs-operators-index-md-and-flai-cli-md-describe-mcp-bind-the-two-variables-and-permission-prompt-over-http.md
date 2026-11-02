---
id: T-2167
type: task
nature: feature
title: docs/users/flai.md, docs/operators/settings.md, docs/operators/index.md, and flai-cli.md describe mcp.bind, the two variables, and permission_prompt over HTTP
status: done
parent: S-0506
owner: alex
created: 2026-10-26T04:50:15Z
updated: 2026-11-02T06:29:01Z
transitions:
  - to: ready
    at: 2026-11-02T06:19:10Z
    by: agent-S-0506
  - to: in-progress
    at: 2026-11-02T06:19:12Z
    by: agent-S-0506
  - to: done
    at: 2026-11-02T06:29:01Z
    by: agent-S-0506
stream: S-0506
tags: [cli, mcp]
touches: [docs/users/flai.md, docs/operators/settings.md, docs/operators/index.md, design/system/flai-cli.md, docs/users/flai-reference.md]
after: [T-2150, T-2152, T-2160, T-2418, T-2419]
usage:
  source: log
  seconds: 589
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 102
      output: 32033
      cache_read: 5948911
      cache_write: 166361
      cost: 2.7967
---
# T-2167 docs/users/flai.md, docs/operators/settings.md, docs/operators/index.md, and flai-cli.md describe mcp.bind, the two variables, and permission_prompt over HTTP

## Work

Update `docs/users/flai.md` § MCP over HTTP: the tool parity with stdio, `flai serve mcp bind` and when to set it (a container on Docker's bridge), `FLAI_MCP_URL` and `FLAI_MCP_TOKEN` as what `flai serve` hands a run, and the sentence on `permission_prompt` over HTTP rewritten to what the ADR of T-2151 settles. Give `mcp.bind` and the two variables rows in `docs/operators/settings.md` so `settings_doc_test` passes, describe them in `docs/operators/index.md`, update `design/system/flai-cli.md` where it describes `flai mcp http` and `flai host`'s MCP servers, and regenerate `docs/users/flai-reference.md` with `make flai-reference`. Waits for T-2150, T-2152, and T-2160, whose behaviour it documents.

## Done when

`settings_doc_test` passes, the markdown lint passes on the four documents, and each sentence on the HTTP server's behaviour matches what the tests of the earlier tasks check.

## Notes

The ADR's own text is T-2151's; this task points to it.
