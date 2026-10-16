---
id: ADR-0162
title: "The component library reaches sub-projects of kind sveltekit or svelte, finds their components by globs the manifest sets, and reads their props with a script flai embeds, run on the sub-project's own svelte/compiler"
status: accepted
date: 2026-10-16
supersedes: []
superseded_by: []
topics: [components, design, cli]
---

# ADR-0162 The component library reaches sub-projects of kind sveltekit or svelte, finds their components by globs the manifest sets, and reads their props with a script flai embeds, run on the sub-project's own svelte/compiler

## Context

E-0026 asks for a component library in the dashboard: every component of the operator's system, seen and worked with, its props changed, and changes to the design system or the source shown live. [design-system-tool.md § The project's own components](../system/design-system-tool.md#the-projects-own-components) left this out of E-0020, because it needs the project's own compiler and build, which only the host has ([ADR-0029](0029-the-dashboard-reaches-a-project-only-through-flai-on-the-host.md), [ADR-0031](0031-the-dashboard-s-container-holds-nothing-of-the-project-a-port-and-two-secrets.md)).

Before anything renders, flai must know which components a project has and what props each takes. Three questions decide that:

- Which sub-projects the library reaches. Storybook covers React, Svelte, and Vue. This repository has one UI sub-project, flaiover, a SvelteKit app on Svelte 5 and Vite.
- Where a sub-project's components are. A SvelteKit app keeps reusable components under `src/lib` and pages under `src/routes`, which do not render alone.
- How props are read. Svelte 5 declares them by destructuring `$props()`, typed in TypeScript. flaiover has 63 `.svelte` files under `src/lib/components`, all `lang="ts"`, 56 with `$props()`. `svelte` is a dependency of flaiover and ships `svelte/compiler`; `svelte2tsx` is not.

The planner asked these on TH-0483, questions 2 and 3, recommending Svelte 5 with Vite first and a script run on the project's own compiler. The operator answered "2. yes, Svelte 5 and Vite first" and "3. let's start with that".

## Decision

The component library reaches sub-projects of kind `sveltekit` or `svelte`, finds their components by globs the manifest sets, and reads their props with a script flai embeds, run with the host's `node` on the sub-project's own `svelte/compiler`.

- **Kinds.** A sub-project whose manifest `kind` is `sveltekit` or `svelte` is supported. `svelte` is a new kind, for a Svelte 5 app or library built with Vite but not SvelteKit. Every other sub-project is listed as unsupported, with its kind. Discovery, prop reading, and the harness sit behind a seam keyed by kind, so that React and Vue can be added later as adapters of their own.
- **Where components are.** A supported sub-project may carry `components` in the manifest:
  - `include`: globs relative to the sub-project's folder. Default `src/lib/**/*.svelte`, so routes are left out.
  - `exclude`: globs taken out of what `include` finds. Default none.
  - `setup`: a module of the sub-project the harness imports before any component, such as the app's stylesheet. Default none.
- **How props are read.** flai embeds a script, `props.mjs`. It runs it with the host's `node`, in the sub-project's folder, so that `svelte/compiler` resolves to the sub-project's own installed copy. The script parses each component with the compiler's `parse` and the modern AST. It reads the `$props()` destructuring: each prop's name, its default as written, whether it is required, and its TypeScript type annotation. A local type alias, a local interface, or an inline object literal is resolved. An imported type is not resolved: its kind is `unknown`, with its text kept. Each type maps to a kind: `string`, `number`, `boolean`, `union` with its literals, `array`, `object`, `function`, `snippet`, or `unknown`.
- **Never executed.** A component is parsed, never imported or run, to list it.
- **Missing tools.** A missing `node`, or a sub-project without `svelte` installed, is an error that names what is missing and how to fix it. A component the compiler cannot parse is listed with its error.

## Consequences

- Props are read as the project's own compiler reads them, at the version the project pins. A Svelte upgrade in the project needs no flai release, unless the AST changes shape.
- Listing a component runs nothing of the project but its compiler, so `flai components list` needs no host action.
- flai needs `node` on the host for this feature, and the sub-project's dependencies installed. Without them the catalogue is an error that says so, and the rest of flai is unaffected.
- An imported prop type gets a JSON editor, not a typed control, until a later story resolves imports.
- React and Vue wait for stories of their own, each an adapter behind the seam: a discovery default, a prop reader, and a harness.
- The manifest gains a key that `flai manifest set` does not reach, so it needs a writer of its own, `flai components config`.
- [component-library.md](../system/component-library.md) describes the whole epic.

## Alternatives considered

- **All three frameworks now.** Triples the harness and the prop extraction before one is proven on a real app.
- **A Go parser in flai.** Needs no Node, but cannot read TypeScript types faithfully, and drifts from the compiler with each Svelte release.
- **svelte2tsx.** Gives full TypeScript types, but it is not a dependency of flaiover, so flai would have to install it into the project or carry it.
