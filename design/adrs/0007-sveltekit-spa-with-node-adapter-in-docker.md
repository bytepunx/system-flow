---
id: ADR-0007
title: SvelteKit SPA with node adapter, shipped as a Docker image
status: accepted
date: 2026-09-15
supersedes: []
superseded_by: []
---

# ADR-0007 SvelteKit SPA with node adapter, shipped as a Docker image

## Context

The brief asks for a SvelteKit single-page app with Tailwind, delivered as a Docker image with a writable mount of the target repo. A browser cannot read a filesystem, so something must serve the repo to the SPA and accept writes.

## Decision

`flaiover` is one SvelteKit project built with `adapter-node`. The root layout sets `ssr = false` so the UI is a client-rendered SPA. Server endpoints under `/api` read and write the repo mounted at `/project`. The Docker image is multi-stage on `node:24-alpine`, listens on 3000, runs as a non-root user, and is published to `ghcr.io/bytepunx/flaiover`. `flai dashboard` pulls and runs it with the repo mounted read-write.

## Consequences

- One deployable, one language on the dashboard side, no separate API service.
- Writes from the browser are ordinary file edits that appear in `git status`; the dashboard never commits.
- The server ports the validation and metrics rules from `flai`; a fixture test keeps them aligned.
- No authentication. The container is bound to localhost by `flai`, and operator docs say not to expose it.

## Alternatives considered

- Static SPA plus a Go API served by `flai`: two processes and a CORS story; puts UI concerns into the CLI.
- Server-side rendering: no benefit for a local single-user tool, and complicates client state for the board.
