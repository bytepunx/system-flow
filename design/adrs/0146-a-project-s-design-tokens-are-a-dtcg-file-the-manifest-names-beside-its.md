---
id: ADR-0146
title: "A project's design tokens are a DTCG file the manifest names, beside its guidance in design/system/design-system.md, and a project has one design system"
status: accepted
date: 2026-10-14
supersedes: []
superseded_by: []
topics: [design, cli, dashboard, template]
---

# ADR-0146 A project's design tokens are a DTCG file the manifest names, beside its guidance in design/system/design-system.md, and a project has one design system

## Context

E-0020 adds a Design menu to the dashboard where a project's design system, its tokens and its guidance, is settled before UI is built. S-0380's finding, [design-system-tool.md](../system/design-system-tool.md), weighed where a project keeps it so that `flai check`, the dashboard, and prime read one source. It recommended tokens in the YAML front matter of a markdown document, because `doc.get`, `doc.save`, prime, and `flai check` handle only markdown today.

On TH-0445, on 2026-10-14, the operator declined that: "I would not like to force more things into a markdown document when it's not the obvious best choice simply because we currently only have API semantics for markdown documents. Introduce file.get and file.save as alternatives that we can leverage for different kinds of files." They gave no other number, so the recommended scope, one design system per project (choice 5), stands.

## Decision

A project's design tokens are a W3C Design Tokens Community Group (DTCG) file that the manifest names, its guidance is the markdown document `design/system/design-system.md`, and a project has one design system.

1. **The token file.** DTCG format, stable version 2025.10, with its Format, Color, and Resolver modules. flai reads the types `color`, `fontFamily`, `dimension` (spacing and radii), `shadow`, and, for the type scale, `typography` or `dimension`. Light and dark are one Resolver modifier, `theme`, with the contexts `light` and `dark`. flai implements that subset of the Resolver module and nothing more of it.
2. **Where it is.** A manifest key names the file. The default is `design/system/design-system.tokens.json`. S-0381 settles the key, and the exact layout against the spec: one file, or a base file, a dark set, and a resolver document, as the spec requires.
3. **The guidance.** `design/system/design-system.md` is a design document like any other: front matter, prose, and threads on its headings. Its `topics` say which stories it reaches.
4. **One per project.** A project has one design system now. A project whose second UI has a visual language of its own can have one per UI sub-project later, each whole, with topics naming its sub-project. That takes a decision of its own.
5. **The documentation convention.** The token file is an exception to the baseline rule "Everything is markdown" (`design/conventions/documentation.md`, [documentation-standard.md](../system/documentation-standard.md)). The operator's answer is followed, and the proposed edit is one exception: a project's design token file, in DTCG format, named by the manifest. S-0381 carries it to `template/` and to the standard.

## Consequences

- flai needs a reader, a validator, and a writer for the DTCG subset (S-0381).
- The dashboard reads and saves the token file through `file.get` and `file.save`, not `doc.get` and `doc.save`.
- Figma Variables, Style Dictionary, Tokens Studio, and Terrazzo read and write the file as it is; no export is needed.
- Tailwind does not read DTCG, so code gets the tokens through a stylesheet flai generates.
- Prime reads only markdown in `design/system` today. Giving agents the token file is a change to prime.
- The documentation convention and standard change, in the template first.
- The JSON holds no prose. When to use which token is said in the guidance only.

## Alternatives considered

- **Tokens in the front matter of `design-system.md`, guidance in its body.** The finding's recommendation: it fits every markdown method as it is. The operator declined forcing tokens into markdown because markdown is what the methods handle.
- **A fenced YAML block in the guidance's body.** Declined for the same reason; `doc.get` does not parse it, and it clutters the rendered guidance.
- **A Tailwind v4 theme stylesheet flai owns, as the source.** The runner-up: it needs no generated file and no `design.stale`. It ties the design system to Tailwind and to one way of writing dark values, and nothing outside Tailwind reads it.
- **A folder of documents under `design/`.** Prime does not read subfolders, and a check rule would have to join several files.
- **One design system per UI sub-project now.** The operator's projects have one UI each, and two documents would need a rule for which one an agent follows.
