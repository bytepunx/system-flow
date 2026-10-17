---
id: ADR-0164
title: "A component's examples are named sets of props in one JSON file per component under design/components, at its sub-project's name and its path with .svelte made .examples.json, saved through file.save, which creates an absent file on an empty hash, and checked by flai check as components.examples"
status: accepted
date: 2026-10-17
supersedes: []
superseded_by: []
refines: [ADR-0147]
topics: [components, cli, dashboard, conventions]
---

# ADR-0164 A component's examples are named sets of props in one JSON file per component under design/components, at its sub-project's name and its path with .svelte made .examples.json, saved through file.save, which creates an absent file on an empty hash, and checked by flai check as components.examples

## Context

A component with required props, such as a card that needs an item, renders on the Components page only once the operator has typed its props into the panel ([component-library.md § Controls for props](../system/component-library.md#controls-for-props-s-0440)). The panel keeps them in the page's URL, so they are lost with the link. S-0443 keeps them as examples: named sets of props, versioned in the repository, that the page lists and loads and that agents keep up to date.

On TH-0483, question 4, the planner recommended `<Component>.examples.json` beside each component, with `design/components/<sub-project>/` as an alternative. On 2026-10-16 the operator answered "lets put supporting files under design/components". Four things follow from that answer and the code as it stands:

- **Reach.** `file.get` and `file.save` reach any file under the design folder that is not markdown ([ADR-0147](0147-the-dashboard-reads-and-saves-a-text-file-that-is-not-markdown-through-file-get.md), `FileReach` in `flai/internal/docedit`). `design/components/` is under it.
- **Create.** `FileSave` saves only a file that exists: `locateFile` refuses an absent path as "does not exist". A component's first example needs a file that is not there yet.
- **Naming a component.** The catalogue names a component by its sub-project and its path in the sub-project, such as `src/lib/components/Card.svelte` ([ADR-0162](0162-the-component-library-reaches-sub-projects-of-kind-sveltekit-or-svelte-finds.md)). The host methods, `flai components list --project`, and the page's `?p=` name a sub-project by its manifest `name`. A sub-project's `path` may have several segments, or be `.`.
- **The panel.** The page shows one theme, `light` or `dark`, or both side by side, and frames at 375, 768, or 1280 CSS pixels, or at the panel's own width (`flaiover/src/routes/design/components/+page.svelte`). A prop's control is chosen by its kind; a function prop takes no value and is given a logger, and a snippet prop takes text.

The documentation convention says "Everything is markdown with YAML front matter", with one exception, the design token file ([ADR-0146](0146-a-project-s-design-tokens-are-a-dtcg-file-the-manifest-names-beside-its.md)).

## Decision

A component's examples are named sets of props in one JSON file per component under `design/components/`, at its sub-project's name and its path with `.svelte` made `.examples.json`. They are saved through `file.save`, which creates an absent file on an empty hash, and checked by `flai check` as `components.examples`.

1. **Where.** The file is `design/components/<sub-project>/<path>`, where `<sub-project>` is the sub-project's manifest `name` and `<path>` is the component's path in the sub-project with `.svelte` replaced by `.examples.json`. `flaiover/src/lib/components/Card.svelte` has its examples in `design/components/flaiover/src/lib/components/Card.examples.json`. The design folder is the manifest's `layout.design`.
2. **The format.** The file is a JSON object with one key, `examples`, an array in the order the page lists them. Each example is an object with these keys and no others:

   | Key | Type | Required | Meaning |
   |-----|------|----------|---------|
   | `name` | string, not empty | yes | Names the example; unique in its file. |
   | `props` | object, prop name to JSON value | yes, may be `{}` | The props the panel sends, as the panel sends them. A function prop is left out. |
   | `theme` | `"light"` or `"dark"` | no | The one theme the page shows. Absent, the page keeps the theme or themes it shows. |
   | `width` | integer, 1 or more | no | The frame's width in CSS pixels. Absent, the page keeps the width it shows. |

   ```json
   {
     "examples": [
       { "name": "A long title", "props": { "title": "Quarterly report for the north region", "open": true }, "theme": "dark", "width": 375 }
     ]
   }
   ```

   The page writes a file with two-space indentation and a final newline, one key per line, so a diff shows each change on its own line. It writes `theme` when it shows one theme and `width` when a width is picked rather than the panel's own.
3. **Reach and the save.** The files are under the design folder, so they are in `file.get`'s and `file.save`'s reach already: ADR-0147's reach is unchanged. ADR-0147's rules of a save gain one: a save with an empty hash, `file.save`'s `hash` empty or `flai file save` without `--hash`, creates a file that is absent, with its folders. An empty hash on a file that exists is refused as a conflict, exit 3, with its current content and hash. The check runs before and after as for any save; when it refuses the created file, the file and the folders the save created are removed. The save commits it as any save does, on main or on the story's branch.
4. **The convention.** The documentation convention gains a second exception to "Everything is markdown": the examples files, `*.examples.json` under `design/components/`. The folder's `README.md` is markdown, as every index is.
5. **The check.** `flai check` reports `components.examples`, a warning, which `flai check --strict` fails on, for:
   - a file that does not parse as JSON, or does not fit the format: a key it does not name, a `name` that is empty or used twice, `props` that is not an object, a `theme` that is not `light` or `dark`, a `width` that is not a positive integer;
   - an example that names a prop its component lacks, misses a required prop that is not a function, or gives a value that does not fit its prop's kind;
   - a file whose component does not exist: no supported sub-project of that name, or no component at that path in its catalogue;
   - a file whose component's props flai could not read, once, with the reader's error, since its examples cannot be checked.

   Each finding names the example and its line: the line of the prop's key for a finding about one prop, the line the example's object opens on for one about the example, the line of the syntax error for JSON that does not parse, and line 1 for a file without its component. A value fits its prop's kind as the panel's control sends it:

   | Kind | Fits |
   |------|------|
   | `string`, `snippet` | a string |
   | `number` | a number |
   | `boolean` | `true` or `false` |
   | `union` | one of its literals |
   | `array` | an array |
   | `object` | an object |
   | `unknown` | any value |
   | `function` | none: the page gives it a logger |

   `null` fits every kind but `function`, because flai leaves `null` and `undefined` out of a type when it reads the kind ([ADR-0162](0162-the-component-library-reaches-sub-projects-of-kind-sveltekit-or-svelte-finds.md)), so the kind cannot say whether `null` is allowed.

   The check reads only the files that end in `.examples.json` under `design/components/`, and reads props, with `node`, only of the components that have one. A project with no such file is not checked and gets no finding. A project that has one but no supported sub-project, or a host without `node`, is passed over with one note: a warning `flai check --strict` passes over, saying what was not checked and why.

## Consequences

- An example is a file in the repository: it is versioned, reviewed with the story that changes it, and seen by agents in the folder the convention lists.
- Moving, renaming, or removing a component, or renaming its sub-project, leaves its examples file without a component; the check reports it, and the story that moved the component moves the file. A prop renamed, removed, made required, or given another kind reports every example it breaks, so a story that changes a component's props keeps its examples true before `flai check --strict` passes.
- A save from the page that would break an example is refused, with the finding, as any save whose content the check reports is ([ADR-0147](0147-the-dashboard-reads-and-saves-a-text-file-that-is-not-markdown-through-file-get.md)).
- Every `flai check` in a project with examples runs `node` once per sub-project that has them, and `file.save` runs the check twice. A project without examples files pays nothing.
- A host without `node` cannot check examples; it is told so and is not failed, since it cannot render the components either.
- `file.save` can now create a file, but only where it could already write one: the create widens no reach. Every editor of `file.save` must send the hash it read; one that sends none on a file that exists is refused, not overwritten.
- The baseline documentation convention and [documentation-standard.md](../system/documentation-standard.md) change, in the template first, and `design/components/` gets a `README.md` here and in the template.
- The catalogue carries each component's examples and the path of its file, whether or not the file exists, so the page knows where a first example goes.
- A new key in an example needs a flai release that knows it, since the check refuses keys it does not name.
- [component-library.md § Examples](../system/component-library.md#examples-s-0443) describes the format.

## Alternatives considered

- **`<Component>.examples.json` beside the component**, the planner's recommendation. Easier to find from the code, but it puts files the app does not use in its source folders and beyond `file.save`'s reach, which would have to widen; the operator chose `design/components`.
- **The sub-project's path in place of its name.** A path may be several segments or `.`, and every method and the page name a sub-project by its name; in most projects, as here, the two are the same.
- **One file per sub-project.** Two stories that add examples to different components would change the same file, and a file found from the component would need searching.
- **One file per example.** Many small files, and a name that must also be a valid file name.
- **Markdown with the examples in front matter.** Forces JSON values into YAML for a file no one reads as prose, which the operator declined for the token file on TH-0445.
- **An error in place of a warning.** The strict check that closes a story already fails on a warning, and a save is refused for a finding it introduces at any level. An error already on a file refuses every save of that file until the save fixes it, so a page that only adds examples could not add one to a file that holds a stale one.
- **A note on every project without a Svelte sub-project.** Most projects made from the template have none; a note they could never clear would be noise.
