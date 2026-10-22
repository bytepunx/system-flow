---
id: T-2035
type: task
nature: improvement
title: An ADR refining ADR-0172 lets an example carry the answers the harness gives a component's requests to the app's own server, and the design describes the format
status: backlog
parent: S-0476
owner: alex
created: 2026-10-22T20:58:53Z
updated: 2026-10-22T20:58:53Z
transitions: []
stream: S-0476
tags: [flai, flaiover]
touches: [design/adrs/README.md, design/system/component-library.md, design/components/README.md, template/root/design/components/README.md]
---
# T-2035 An ADR refining ADR-0172 lets an example carry the answers the harness gives a component's requests to the app's own server, and the design describes the format

## Work

Propose the solution I-0164 leaves open, as the story's goal asks, and record it as an ADR that refines [ADR-0172](../../../design/adrs/0172-the-component-harness-answers-a-framed-component-s-requests-to-the-app-s-own.md), whose "Answer with fixtures" alternative is the one to take. Write it with `flai adr new` (or the MCP tool `adr_new`), which names the file and adds it to `design/adrs/README.md`; the ADR's own path is not in this task's touches, since its number is not known until it is made, and `flai task done` widens the touches with it.

The decision to record:

- An example in `<path>.examples.json` ([ADR-0164](../../../design/adrs/0164-a-component-s-examples-are-named-sets-of-props-in-one-json-file-per-component.md)) may carry `answers`, an array, beside `name`, `props`, `theme`, and `width`. Each answer has `path` (string, not empty: the request's path, with its query when it has one), and optionally `method` (an HTTP method in upper case; absent, any method), `status` (integer 100–599; default 200), `headers` (object of string to string), `body` (any JSON value: an object or array is sent as JSON with `content-type: application/json` unless `headers` sets one; a string is sent as it is), and `events` (array of `{event?, data}` for an `EventSource` of that path).
- Framed, the harness answers a `fetch` of the app's own server from the example's answers in place of the 501 when one matches: the same method, or no method, and the same path with its query, else the same path without the query. The first match in the array's order wins. An `EventSource` whose path has an answer with `events` opens: it fires `open`, then each event in order, named `message` unless `event` names it, with `data` as a string (an object or array as JSON), and stays open until the component closes it. Everything unmatched keeps ADR-0172's behaviour: 501, or a failed `EventSource`, and the `event` message. A `WebSocket` is never answered from an example.
- The harness still posts an `event` for each request it answers, named `fetch` with `[method, path, status]` or `EventSource` with `[path]`, so the page's log shows what the component asked and what it got.
- The page sends the chosen example's `answers` with `select` and `props`, and keeps them when it saves an example; a flai older than this release ignores the field and still answers 501.

Then describe the same in `design/system/component-library.md`: the `answers` key, with its sub-keys, in § Examples › The file's table; the matching rule under § The dev server and its harness › The harness, beside the app's-own-requests bullet; the new `answers` field of `select` and `props` in § The protocol's table; and the `fetch` event's third argument. Say in `design/components/README.md`, and in its template copy `template/root/design/components/README.md`, that an example may carry the answers to the component's requests. Bump `updated` where the document has one.

This task waits for nothing: the format it fixes is what every later task builds.

## Done when

- [ ] A new ADR under `design/adrs/`, accepted, refining ADR-0172, records the `answers` key, the matching rule, what stays 501, and the page's part, and `design/adrs/README.md` lists it
- [ ] `design/system/component-library.md` describes `answers` in § The file, § The harness, and § The protocol
- [ ] `design/components/README.md` and `template/root/design/components/README.md` say an example may carry answers
- [ ] `flai check --strict` and the markdown lint pass on the changed files

## Notes

Planned by planner-S-0476. The format above is the planner's proposal; change it in the ADR if the code says otherwise, and keep the later tasks in step.
