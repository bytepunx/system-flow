---
id: T-0246
type: task
nature: research
title: "Try the options that need no account against a scratch service: SSE, streamed responses, bearer tokens"
status: done
parent: S-0069
owner: alex
created: 2026-09-20T04:07:27Z
updated: 2026-09-20T04:26:48Z
transitions:
  - to: ready
    at: 2026-09-20T04:26:48Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T04:26:48Z
    by: system-flow
  - to: done
    at: 2026-09-20T04:26:48Z
    by: system-flow
stream: S-0069
tags: []
---
# T-0246 Try the options that need no account against a scratch service: SSE, streamed responses, bearer tokens

## Work
Build a scratch service on this host that behaves like flaiover where it matters: a page, a server-sent events stream that runs for minutes, a streamed newline-delimited JSON response, and an endpoint behind a bearer token. Put it behind each option that needs no paid account, no domain, and none of the operator's accounts (for example a Cloudflare quick tunnel, and a self-hosted tunnel between two containers), and record what was run and what was seen: did SSE stay open, did the stream arrive line by line or buffered, did the bearer header pass. Never the operator's dashboards, accounts, keys, or this repository. Mark everything else as not tried, and turn anything that needs the operator's account or domain into a question for them.

## Done when
- Each trial is written up with the command and the observation
- Each claim in the finding is marked tried or not tried
- The scratch service, containers, and tunnels are removed

## Notes
