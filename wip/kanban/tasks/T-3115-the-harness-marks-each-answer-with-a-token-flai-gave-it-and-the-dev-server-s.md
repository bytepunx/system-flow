---
id: T-3115
type: task
nature: remediation
title: The harness marks each answer with a token flai gave it, and the dev server's wait takes only an answer that carries it
status: backlog
parent: S-0714
owner: alex
created: 2026-11-11T09:05:31Z
updated: 2026-11-11T09:05:31Z
transitions: []
stream: S-0714
tags: [flai, components]
touches: [flai/internal/components/harness/plugin.mjs, flai/internal/components/harness.go, flai/internal/components/harness_test.go, flai/internal/components/serve.go, flai/internal/components/devserver.go, flai/internal/components/devserver_test.go]
after: [T-3114]
---
# T-3115 The harness marks each answer with a token flai gave it, and the dev server's wait takes only an answer that carries it

## Work

Nothing in an answer tells flai that its own vite gave it, so a 200 or a 404 from whatever holds the port is read as vite's. Give each start a random token. `WriteHarness` in `flai/internal/components/harness.go` writes it into the plugin's options, `flai/internal/components/harness/plugin.mjs` adds a middleware that sets a response header, such as `X-Flai-Harness: <token>`, on every answer the dev server gives, and `Serve` in `flai/internal/components/serve.go` hands the same token to `StartDevServer`. `await` in `flai/internal/components/devserver.go` takes a 200 only when the header carries the token, and `DevServer` gives the token, so that a caller can check its own requests.

This reproduces I-0342's cause in a test: in `flai/internal/components/devserver_test.go`, a Go server on the port `pickPort` answers, giving 200 at the harness path without the header, is not taken for vite, and the start fails or moves to another port saying another server answers there.

Waits for T-3114, which changes `devserver.go` and `devserver_test.go` too.

## Done when

- Every answer of the harness's dev server carries the token header, and `harness_test.go` checks the plugin's options carry the token.
- `StartDevServer` does not take a 200 without the token for its vite's, and a test in `devserver_test.go` shows it with a Go server on the picked port.
- `flai test` on the six touched files passes.

## Notes
