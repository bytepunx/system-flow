---
id: ADR-0187
title: "An example may carry the answers the component harness gives a framed component's fetch and EventSource to the app's own server, and the harness answers 501 only what none matches"
status: accepted
date: 2026-10-22
supersedes: []
superseded_by: []
refines: [ADR-0164, ADR-0172]
---

# ADR-0187 An example may carry the answers the component harness gives a framed component's fetch and EventSource to the app's own server, and the harness answers 501 only what none matches

## Context

Framed, the component harness answers every request a component makes to the app's own server itself, with 501 for a `fetch` and an `error` for an `EventSource` ([ADR-0172](0172-the-component-harness-answers-a-framed-component-s-requests-to-the-app-s-own.md)). That keeps the dev server's server code, which has no dashboard token and no host channel, from answering 500. But a component whose view depends on such an answer cannot be seen in its answered state. S-0449's headless pass found fourteen of flaiover's components so ([I-0164](../issues/I-0164-the-components-page-cannot-show-a-component-that-asks-the-app-s-server-in-its-answered-state-since-the-harness-answers-every-such-request-with-501.md)). ADR-0172 named "answer with fixtures" as the alternative to follow. An example is a named set of props, with no other key allowed ([ADR-0164](0164-a-component-s-examples-are-named-sets-of-props-in-one-json-file-per-component.md)).

## Decision

An example may carry `answers`, an array of the answers the harness gives the component's requests to the app's own server. The page sends them to the harness with the props, and the harness answers a request from the first that matches it, and as ADR-0172 says when none does.

- **An answer** has `path`, the request's path, starting with `/`, with a query when it answers that query alone, and optionally `method`, `status`, `headers`, and `body`, for a `fetch`, or `events`, for an `EventSource`.
  - `method` is an HTTP method, compared without case; absent, it matches every method.
  - `status` is an integer from 200 to 599, 200 when absent.
  - `headers` maps header names to strings.
  - `body` is any JSON value: a string is sent as it is, with `content-type: text/plain; charset=utf-8`, anything else as JSON, with `content-type: application/json`, unless `headers` names a `content-type`. Absent, the body is empty, as it must be with status 204, 205, or 304.
  - `events` is an array of `{event?, data, id?}`. `data` is sent as it is when a string, as JSON otherwise; `event` absent is `message`. An answer with `events` takes no `status`, `headers`, or `body`, and no `method` but `GET`.
- **Matching.** A `fetch` matches an answer without `events`, an `EventSource` one with them. Of those, the first answer, in the example's order, whose `method` matches and whose `path` is the request's path with its query wins; when none is, the first whose `path` has no query and is the request's path without its query.
- **A fetch answered** resolves with a `Response` of the answer's status, headers, and body, after the current task, as a request over the network would.
- **An EventSource answered** opens no connection. After the current task it fires `open`, with `readyState` 1, then each event in order, each in a task of its own, and stays open until the component closes it.
- **Unmatched.** A request no answer matches gets ADR-0172's answer: a `fetch` 501, an `EventSource` an `error` and closed. A `WebSocket` is never answered from an example.
- **The log.** Each request is still posted as an `event`: `fetch` with the method, the path, and the status it got, 501 when unmatched; `EventSource` with the path and 200 when an answer opened it, 501 when none did.
- **The protocol.** `select` and `props` take `answers`, an array, whole as props are: each message's answers replace the last's, and a message without them leaves none.
- **The catalogue and the check.** `flai components list` gives an example's `answers` when it has them. `flai check` reports `components.examples` for an answer that does not fit the format. Answers are fixtures, not checked against the app's routes: a story that changes a route's answer changes the examples that copy it, as one that changes a component's props changes its examples.

## Consequences

- A component whose view comes from the app's server can be shown answered, from an example, on the Components page and in a headless pass, with no dashboard token and no host channel.
- An example's file grows by the answers, which can be large; they stay out of the page's URL.
- A component whose answers change shape shows stale data until its examples follow. Nothing checks that.
- A component that writes, such as an accept, gets the example's answer and acts on nothing.
- What ADR-0172 leaves uncovered, `XMLHttpRequest`, `sendBeacon`, a form's submission, and a `WebSocket`, stays so.
- [component-library.md](../system/component-library.md) describes the format under § The file, the harness's answers under § The harness, and the protocol's field under § The protocol.

## Alternatives considered

- **Proxy to the dashboard.** ADR-0172 refused it: a preview's writes would act on the project.
- **Answers per component, in a file of their own.** Every example of a component would share them, but a component's states differ by its answers as much as by its props, such as a list that is empty or full; one example carries both.
- **Match by pattern or by order of calls.** A path pattern, or answers taken in turn, would cover more, but an exact path, with or without its query, covers what flaiover's components ask, and reads plainly in the file.
- **Answer a WebSocket too.** A socket's messages are a conversation, not an answer; no listed flaiover component needs one.
- **A form per answer on the page.** A JSON editor, as an object prop has, is the smallest that fits; a form can follow.
