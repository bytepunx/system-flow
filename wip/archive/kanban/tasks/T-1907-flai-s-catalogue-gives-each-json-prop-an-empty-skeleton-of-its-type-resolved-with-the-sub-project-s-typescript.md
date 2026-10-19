---
id: T-1907
type: task
nature: feature
title: flai's catalogue gives each json prop an empty skeleton of its type, resolved with the sub-project's typescript
status: done
parent: S-0449
owner: alex
created: 2026-10-19T13:58:33Z
updated: 2026-10-19T14:14:37Z
transitions:
  - to: ready
    at: 2026-10-19T13:58:48Z
    by: agent-S-0449
  - to: in-progress
    at: 2026-10-19T13:58:49Z
    by: agent-S-0449
  - to: done
    at: 2026-10-19T14:14:36Z
    by: agent-S-0449
stream: S-0449
tags: []
touches: [flai/internal/components/props.mjs, flai/internal/components/props.go, flai/internal/components/props_test.go, flai/internal/components/testdata, design/system/component-library.md, docs/users/flai.md, design/system/flai-cli.md, docs/users/flai-reference.md, flai/cmd/components.go]
usage:
  source: log
  seconds: 947
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 124
      output: 51444
      cache_read: 7078800
      cache_write: 212336
      cost: 3.7843
---
# T-1907 flai's catalogue gives each json prop an empty skeleton of its type, resolved with the sub-project's typescript

## Work

Criterion 3, as TH-0515 settled it: option (A), which the operator confirmed. flai's catalogue carries, for each prop whose control is the JSON editor (kind `array`, `object`, or `unknown`), a `skeleton`: an empty value of the prop's type.

- `flai/internal/components/props.mjs`: resolve each such prop's type with the sub-project's own `typescript` package, loaded from the sub-project as the svelte compiler is. Include imported types and the sub-project's path aliases, such as SvelteKit's `$lib`, from its tsconfig. Use one program for all of a run's components, not one per component. Build the skeleton from the resolved type:
  - `string` gives `""`, `number` gives `0`, and `boolean` gives `false`.
  - A union of literals gives its first member. A union with `undefined` or `null` gives the skeleton of the rest. Another union gives its first member's skeleton.
  - An array or tuple gives one element, the item's skeleton.
  - An object or interface gives each property, optional ones included, with its skeleton. A `Record` or index signature gives `{}`.
  - A function gives nothing: leave the key out.
  - A cycle, or depth past a small bound, gives `null`, or `{}` for an object.
  - `Date` gives `""`.
  - A type that cannot be resolved falls back by kind: `[]` for array, `{}` for object, `null` for unknown. So does every prop when the sub-project has no `typescript`.
- `flai/internal/components/props.go`: `Prop` gains `Skeleton json.RawMessage` with the tag `json:"skeleton,omitempty"`, and a one-line doc comment. It is set for every array, object, and unknown prop, and for none of the others.
- Tests in `flai/internal/components/props_test.go`, with fixtures under `flai/internal/components/testdata/`: an imported interface through a path alias, an array of it, optional and union members, a cycle, an unresolvable type, and a sub-project without typescript. Each test fails without the change.
- Measure `flai components list --project flaiover --json` on this worktree before and after, and say the times. Keep it near today's time; if it is much slower, say how much.
- Docs, in the same change, in `design/system/component-library.md`:
  - § The catalogue › Props: imported types are now resolved for a json prop's skeleton. Its kind stays as today.
  - § Controls for props › What is sent, and the defaults text: a json prop's field starts with its skeleton, which is sent. A literal default still wins. Reset returns the field to its skeleton. Cite TH-0515.
  - § The protocol and catalogue tables: the `skeleton` field, where they list Prop's fields.
  - `docs/users/flai.md`: the `skeleton` field in `flai components list`'s output.

It runs with T-1908 (layer 3). They share no path. The contract between them: `skeleton` is a JSON value on the prop, present for array, object, and unknown props only.

## Done when

- `flai components list --project flaiover --json` gives every array, object, and unknown prop a `skeleton`, and one built from an imported type, such as `Story[]`, has its fields.
- The props tests pass with `flai test` on the changed paths, and fail without the change.

## Notes
