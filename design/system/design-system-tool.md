---
title: A design system tool for the dashboard and agents
updated: 2026-10-13
status: active
topics: [design, dashboard, cli]
---

# A design system tool for the dashboard and agents

The finding of S-0380, for E-0020. The operator asked for a Design menu in the dashboard where a project's design system is settled before UI is built: its tokens and its guidance, previewed live on sample Svelte and Tailwind components, and changed by the operator and agents together, so that the visual language of a solution can be judged before it is delivered. This document says what the dashboard, flai, and the prime pack offer such a tool today, where and in what format a project could keep its design system so that flai, the dashboard, and agents read one source, how the dashboard could preview it without building the project, and how the operator and agents could work on it together. Nothing here is built. The operator's decision is recorded at the end.

Every statement about this repository comes from its files, read on 2026-10-13, and names the file or ADR it comes from. Design documents are linked; code is named by its path.

## Today

No project has a design system that flai or the dashboard knows of. The dashboard has a design system of its own, a token layer in one stylesheet, built into its image. flai serves the dashboard Markdown under three folders and nothing else. The prime pack selects design documents by their topics. Each subsection names what a design system tool could build on and what stands in its way.

### The dashboard's token layer

flaiover's theme is a token layer in `flaiover/src/routes/layout.css`, imported once by `flaiover/src/routes/+layout.svelte` ([flaiover-dashboard.md § Theme](flaiover-dashboard.md#theme-s-0044)):

| Layer | What it is |
|-------|------------|
| Tokens | CSS custom properties named `--t-<name>`, 39 per theme, every one a six-digit hex colour: surfaces (`ground`, `surface`, `raised`), text (`ink`, `ink-soft`, `muted`), lines, `primary`, `accent`, the four states (`good`, `warn`, `danger`, `info`) with their `-soft` and `-strong` steps, the agent dots, and the board's nature and type colours. |
| Light and dark | The light values are on `:root`, the dark ones on `[data-theme='dark']`. `@custom-variant dark` makes Tailwind's `dark:` follow the attribute, not the system preference. |
| Theme choice | `data-theme` on `<html>`, stamped before first paint by a script in `flaiover/src/app.html` from `localStorage` (`flaiover-theme`) or `prefers-color-scheme`, and kept current by `flaiover/src/lib/theme.svelte.ts`. |
| Tailwind mapping | `@theme inline` maps each token to a colour of Tailwind's: `--color-ground: var(--t-ground)` gives `bg-ground`, `text-ground`, `border-ground`. `inline` writes `var(--t-ground)` into each utility, so a value set on an element's `--t-*` reaches every utility below it. |
| Rendered markdown | `.prose` sets `@tailwindcss/typography`'s `--tw-prose-*` variables from the tokens. |

Only the colour namespace is themed. The other namespaces of a Tailwind v4 theme are left at Tailwind's defaults:

| Namespace | In `layout.css` | What the components use |
|-----------|-----------------|-------------------------|
| Colours (`--color-*`) | 39 tokens, light and dark | `bg-*`, `text-*`, `border-*` of the tokens |
| Fonts (`--font-*`) | None | Tailwind's default sans and `font-mono` |
| Type scale (`--text-*`) | None | Tailwind's default `text-xs` to `text-*` |
| Spacing (`--spacing`) | None | Tailwind's default step |
| Radii (`--radius-*`) | None | Tailwind's default `rounded`, `rounded-lg`, `rounded-full` |
| Shadows (`--shadow-*`) | None | Tailwind's default `shadow-sm`, `shadow-lg` |

Chart colours are a second palette outside the stylesheet: `flaiover/src/lib/viz/palette.ts` holds categorical slots as hex values per theme, which the charts read from TypeScript.

### Its contrast test

`flaiover/src/lib/theme.test.ts` reads `layout.css` and parses the `:root` and `[data-theme='dark']` blocks, so the test cannot drift from the stylesheet. Its parser reads only `--t-<name>: #rrggbb`; a token in any other colour syntax is not seen. For each theme it checks:

| Check | Threshold |
|-------|-----------|
| The core, nature, and type tokens are present as six-digit hex | Present |
| Text pairs: `ink`, `ink-soft`, `muted`, and `accent` on `ground` and `surface`; `on-primary` on `primary`; each state on its `-soft` and on `ground` | 4.5:1 |
| `ink`, `muted`, and `danger` on each nature tint, and `ink` and `muted` on each lit nature tag | 4.5:1 |
| Control pairs: `primary`, `accent`, and `line-strong` against `ground`, and `line-strong` against `surface` | 3:1 |
| Every nature and every type has a colour of its own, and each lit tag differs from its tint | Distinct |

A last group pins the five brand colours to their roles and records why the brand green and cyan are fills, not text, in the light theme. The contrast arithmetic is the WCAG relative luminance formula, written in the test.

### How flaiover builds its CSS

| | |
|-|-|
| Tailwind | 4.3 ([tailwind.md](../tech/tailwind.md)); `tailwindcss` and `@tailwindcss/vite` 4.3.3 in `flaiover/pnpm-lock.yaml`, with the native `@tailwindcss/oxide` |
| Svelte | 5.57.0, SvelteKit 2 with `@sveltejs/adapter-node`, Vite 8 |
| Integration | The `@tailwindcss/vite` plugin in `flaiover/vite.config.ts`; no PostCSS and no `tailwind.config.js` |
| When CSS is made | At `vite build`, from the class names Tailwind finds in flaiover's own sources; `layout.css` names no `@source` |
| Image | `flaiover/Dockerfile` runs `pnpm build` then `pnpm prune --prod` in a build stage and copies `build/` and the production `node_modules` to the runtime stage |
| At run time | `tailwindcss`, `@tailwindcss/vite`, `svelte`, and `vite` are devDependencies in `flaiover/package.json`, so the running image holds no Tailwind compiler, no Svelte compiler, and no Vite. The lock file has no browser build of Tailwind. |

[tailwind.md](../tech/tailwind.md) says the tokens are kept in one `app.css`; they are in `src/routes/layout.css`.

### The site menu

`flaiover/src/lib/sitemenu.ts` (S-0172) defines the header's two tiers as data, `SITE_MENU`: three groups, each with its pages in a row below.

| Group | Pages | Also claims |
|-------|-------|-------------|
| Workflow | Overview, Board, Inbox, Threads, Messages, Activity, Planner, Orchestrator, Analyzer | `/items`, `/review`, `/new` |
| Status | Charts, ADRs, Docs, Search | `/edit` |
| Host | Updates, Settings, License | |

A new top-level group, such as Design, takes four changes: its key in the `GroupKey` union and its pages' keys in `PageKey`, an entry in `SITE_MENU` in the order it is shown, each page's link in the `HREF` map of `flaiover/src/lib/components/SiteMenu.svelte`, and the routes under `flaiover/src/routes`. `flaiover/src/lib/sitemenu.test.ts` asserts the groups and their pages in order, and which group a path belongs to, so it changes with them. `SiteMenu.svelte` renders whatever `SITE_MENU` holds.

### The document view and editor

[ADR-0023](../adrs/0023-documents-are-saved-through-flai.md) has documents edited in the dashboard saved through flai. The host methods are in `flai/internal/hostapi`, and the rules in `flai/internal/docedit/docedit.go`:

| Method | What it does |
|--------|--------------|
| `docs.tree` | Every file under the manifest's design, docs, and wip folders, recursively, with its title (`flai/internal/hostapi/docs.go`) |
| `doc.get` | One Markdown file under those folders: its path, its front matter parsed from YAML into a map, its body, and its raw text. Any other path, or a file that is not `.md`, is refused (`docPath`). |
| `doc.show` | `flai doc show`: the raw content, its SHA-256, and its edit mode, `full`, `body`, or `none`, with the reason (`flai/internal/hostapi/writes.go`) |
| `doc.save` | `flai doc save` with the hash the editor loaded; exit 3 is a conflict, exit 4 a refusal |

What a save does to a design or docs file, which is in mode `full`:

- **Front matter is preserved as written.** It is saved byte for byte as the editor sent it. The one change flai makes is to set `updated` to today when the file has the key and the editor left its value as it was.
- **Front matter is validated only by `flai check`.** flai runs the whole check before and after writing. Any finding the content introduces, at any level and on any path, restores the old file and refuses the save with the findings; an error already on the file refuses too. There is no schema for a document's front matter beyond what the check reads (below).
- **Concurrent edits** are caught by the hash; a conflict carries the current content and a diff.
- **Commit.** Unless `dashboard.autocommit` is false, the save commits that one path on the main checkout's current branch, as `docs:`.
- Work items, narratives, and `board.md` are mode `body`: a save that changes their front matter is refused. Threads, issues, generated indexes, the archive, and accepted ADRs (but for their topics) are mode `none`.

The editor is `/edit/<path>`, `flaiover/src/routes/edit/[...path]/+page.svelte` with `flaiover/src/lib/components/DocEditor.svelte`. It splits the document into raw front matter and body (`splitRaw` in `flaiover/src/lib/edit.ts`) so an unchanged document round-trips exactly. The front matter is a plain YAML textarea; the body is a Markdown textarea beside a preview. The view `/docs/<path>` and the preview render with `flaiover/src/lib/markdown.ts`: markdown-it with raw HTML allowed, a fence rule that turns `mermaid` blocks into diagrams, and shiki for other fences, all in the browser. The rendered HTML goes into the page with `{@html}`; flaiover has no iframe, shadow root, sanitizer, or Content-Security-Policy anywhere in `flaiover/src`.

### What the dashboard can reach of a project

[ADR-0029](../adrs/0029-the-dashboard-reaches-a-project-only-through-flai-on-the-host.md) has the dashboard reach a project only through flai on the host, over a WebSocket flai opens, by named methods with typed arguments. There is no method that reads or writes an arbitrary file, and none that takes a command line. [ADR-0031](../adrs/0031-the-dashboard-s-container-holds-nothing-of-the-project-a-port-and-two-secrets.md) gives the container a port and two secrets and nothing of the project: no volume, no git, no flai.

For previewing a project's design system this means:

- The dashboard cannot read a project's stylesheet, `tailwind.config`, components, or `package.json`. It can read Markdown under design, docs, and wip, and the structured answers of the other methods.
- The dashboard cannot build a project. Its image holds no compiler for Tailwind or Svelte (above), and the project's sources are not in the container to build.
- Anything the dashboard shows of a project's tokens must reach it as Markdown it reads with `doc.get`, or as the answer of a new host method flai implements.
- Host actions that run on the host, such as `checks.run`, `test.run`, and `verify.run`, are off until the operator enables each by name in the host's flai configuration (ADR-0029).

### Another branch's version of a file

The dashboard cannot read a whole file as a story's branch has it:

- `doc.get` and `doc.show` read the project's checkout on the host, where flai serve runs them. Neither takes a story or a ref.
- `stream.diff` (`flai/internal/hostapi/reads.go`, `flai/internal/storygit/diff.go`) gives a story branch against its merge base with main: each changed file's status, counts, and unified hunks, for any file type, cut at 64 KB a file and 768 KB in all. The review page draws it (`DiffView.svelte`). It carries hunks, not whole files.
- Two reads already use a story's worktree: `issue.list` with `story` reads the issues the story recorded from its worktree (S-0198), and `test.run` and `verify.run` run in it. None of them returns a document.

### How flai prime selects design documents

`flai prime --story` builds the story's pack in `flai/internal/context` ([agent-context.md § As built](agent-context.md#as-built), [ADR-0047](../adrs/0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md), [ADR-0049](../adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md)):

| Step | What it selects | How it is given |
|------|-----------------|-----------------|
| Named | A document the story, its epic, or its tasks link, or name by `ADR-nnnn` | Whole, but for a large file named only by its path written out, which is briefed |
| Briefed | A design or tech file with a section whose topics match the story's | Its title, first paragraph, and outline, unless the budget lets it load whole |
| Ranked | Sections that match the story's words | What the budget leaves |

- `LoadDocs` in `flai/internal/context/design.go` reads the Markdown files directly in `design/system`, `design/tech`, and `design/adrs`. It does not read subfolders, nor any other folder under `design`.
- A document's topics are its front matter's `topics:`; a heading may carry its own with `<!-- topics: ... -->`, and a section without its own takes its parent's (`flai/internal/topics/topics.go`). A design document without topics is selected only by links and ranking.
- A story's topics are its own and its epic's, the name, kind, and tags of every sub-project its tags name or its touches reach, `code` when one of those is not the template, and `all` (`flai/internal/topics`). In a project whose UI is a sub-project of kind `sveltekit`, a document with `topics: [sveltekit]` reaches every story that touches it.
- A brief is not the document: the agent reads what bears on its story with `doc_get` and a heading ([conventions.md § Priming](conventions.md#priming)).

### How flai check validates design documents

`flai check` (`flai/internal/check/check.go`, `documentation`) walks the design and docs folders and, for every Markdown file but a README:

| Rule | Level | What it checks |
|------|-------|----------------|
| `doc.front-matter` | Warning | The file has front matter and it parses as YAML |
| `doc.title`, `doc.updated` | Warning | `title` and `updated` are set |
| `doc.topics` | Warning | A file in `design/system` or `design/tech` declares `topics` |
| `doc.topic` | Warning | Each topic, on the front matter or a heading, is one a story can have: `all`, `code`, a sub-project's name, kind, or tag, or a story's or epic's topic |
| `adr.*` | Warning | An ADR's `id`, `title`, `status`, `date`, index row, and decision sentence |

No rule reads a document's other front matter keys. The walk skips the conventions, issues, experiments, and analysis folders, which have rules of their own that validate each file against a schema: `experiments` calls `experiment.Validate` (`flai/internal/experiment`) and `analysis` calls `analysis.Validate` (`flai/internal/analysis`), each reporting a problem against the key or heading it is on. That is the pattern a rule for a design system document would follow: a function of the checker, called from `Run`, with a package that validates the file. Every rule's finding refuses a dashboard save that introduces it (above).

### What describes a project's design system

Nothing. A search on 2026-10-13 of `design/`, `docs/`, `template/`, `flai/`, and `flaiover/src` for "design system", "design-system", and design token files found no document, convention, command, host method, or manifest key for a project's design system. The template carries no sub-project scaffold (`template/projects` is empty) and no stylesheet or Tailwind configuration, and its `design/system` holds only an overview and an index. The only mentions are E-0020's own backlog stories and tasks under `wip/kanban`.

## Where a project keeps its design system

A design system here is tokens, the values code uses, and guidance, the words that say when to use which. Both need one source that `flai check` validates, the dashboard reads through a host method, and prime gives an agent.

### Options

| Option | Documentation standard | Host methods | `flai check` | What an agent reads | What code uses |
|--------|------------------------|--------------|--------------|---------------------|----------------|
| 1. Tokens in the front matter of `design/system/design-system.md`, guidance in its body | Fits: Markdown with YAML front matter in the builders' folder | `doc.get` returns the tokens parsed into a map; `doc.save` saves them byte for byte | One `design.tokens` rule reads the front matter; a bad token then refuses a dashboard save for free | Named or selected by topics, it loads whole, front matter included; a brief gives the guidance's outline, not the tokens | A rendered `@theme` block (S-0385's subject) |
| 2. A folder, `design/design-system/`, tokens and guidance in several files | Bends it: a folder the standard does not list, needing a README and a convention line | `doc.get` reads any `.md` under design, file by file | The documentation walk reads it; a rule must join several files | `LoadDocs` does not read subfolders: prime never selects it | As 1 |
| 3. Tokens in a fenced `yaml` block in one document's body | Fits | `doc.get` returns the block as body text; the dashboard parses it itself | As 1, after finding the block | As 1; the tokens also fill a brief's ranked sections | As 1 |
| 4. A DTCG `design.tokens.json` beside a guidance document | Breaks it: not Markdown | `doc.get` refuses a file that is not `.md`, and `doc.save` cannot write it: new methods to read and save it | A new walk; no refused save | Prime reads no JSON | Figma Variables, Style Dictionary, and Terrazzo read DTCG; Tailwind does not |
| 5. A Tailwind v4 stylesheet in the UI sub-project, its `@theme` block the source | Breaks it: CSS outside design | The dashboard cannot read it (ADR-0029); a new method would read a sub-project's file | flai would parse CSS; `@theme` holds one value per variable, so dark values and roles are the file's own convention | Not primed; the guidance lives elsewhere | Exactly what code uses |

The DTCG format's first stable version is 2025.10 ([designtokens.org](https://www.designtokens.org/), read 2026-10-13). It recommends the extension `.tokens` or `.tokens.json`, writes a colour as an object with `colorSpace` and `components`, and expresses light and dark through its Resolver Module, a second file. Interop does not need DTCG as the source: flai can write a `.tokens.json` from option 1 for a project that wants one.

### The token schema

Tokens sit under one key, `tokens`, in namespaces named after Tailwind v4's, so rendering them as `@theme` is a rename:

| Namespace | Tailwind variable | Value | Required |
|-----------|-------------------|-------|----------|
| `color.<name>` | `--color-<name>` | `{light, dark}`, or one value for both | The roles below |
| `font.<name>` | `--font-<name>` | A CSS `font-family` stack | `sans`, `mono` |
| `text.<step>` | `--text-<step>`, `--text-<step>--line-height` | `{size, line-height}`: a length, and a length or a unitless number | `xs`, `sm`, `base`, `lg`, `xl`, `2xl`, `3xl`, `4xl` |
| `spacing` | `--spacing` | One length, the step every spacing utility multiplies | Yes |
| `radius.<name>` | `--radius-<name>` | A length | `sm`, `md`, `lg` |
| `shadow.<name>` | `--shadow-<name>` | A CSS `box-shadow` whose colours follow the colour rule | `sm`, `md`, `lg` |

Lengths are `rem` or `px`. Tailwind's other namespaces, such as breakpoints, weights, and easing, stay at its defaults until a project asks for them.

**Colour syntaxes.** Two: six-digit hex, `#rrggbb`, and `oklch(L C H)`, with `/ alpha` only in shadows and extra tokens. flai must compute contrast, which is defined on sRGB. Hex is sRGB; oklch converts to it by fixed arithmetic, with no new dependency. oklch is what Tailwind v4's own palette is written in, and it steps evenly in lightness, which makes a dark value easier to derive from a light one. An oklch colour outside sRGB is clipped for contrast and reported. `rgb()`, `hsl()`, and names add nothing hex lacks; `var()`, `color-mix()`, and `light-dark()` must be resolved before contrast can be computed; `color(display-p3 ...)` has no sRGB contrast. A hex value is quoted, since `#` starts a YAML comment; an unquoted one parses as empty and is a finding.

**Semantic roles.** A fixed set of colour roles is required, because the preview's components are written against names they can rely on ([The fixed vocabulary](#the-fixed-vocabulary)): `background`, `surface`, `text`, `text-muted`, `border`, `primary`, `on-primary`, `secondary`, `on-secondary`, `focus`, and the states `success`, `warning`, `danger`, `info`, each with a `-soft` tint for its alert. Any other name is an extra token: validated, rendered, and shown as a swatch. The roles give the pairs flai computes, at the thresholds flaiover's own test uses: text, muted text, and each state on `background`, `surface`, or its tint at 4.5:1; `on-primary` and `on-secondary` on their fill at 4.5:1; `border` and `focus` on `background` at 3:1.

Contrast should not be a `flai check` finding. Every finding a save introduces refuses the save (§ Today), so a contrast rule would stop the operator saving a colour while trying it. Parse errors are findings; contrast is reported by flai's design command and on the Design page.

**One per project.** One design system per project now: the operator's projects have one UI each, and two documents would need a rule for which one an agent follows. A project that grows a second UI with its own visual language can have a document per sub-project later, `design/system/design-system-<sub-project>.md`, each whole rather than inheriting, with topics naming its sub-project so prime gives each UI story its own.

**The name.** `design/system/design-system.md`. No file, link, or code in this repository or the template uses it (searched 2026-10-13), and it is the operator's own term. It sits beside this finding, `design-system-tool.md`, which describes the tool, not a project's design system. Its `topics` name the UI sub-project or its kind, so every story that touches the UI is briefed the guidance.

```yaml
---
title: Design system
updated: 2026-10-13
status: active
topics: [web]                  # the UI sub-project, or its kind
tokens:
  color:
    background: { light: "#ffffff", dark: "#0b0f14" }
    text:       { light: "#111827", dark: "#e5e7eb" }
    primary:    { light: "oklch(0.52 0.18 255)", dark: "oklch(0.72 0.14 255)" }
    on-primary: { light: "#ffffff", dark: "#0b0f14" }
    # ... the other roles
    brand-sun:  "#f5b700"      # an extra, one value for both themes
  font:
    sans: "Inter, ui-sans-serif, system-ui, sans-serif"
    mono: "ui-monospace, Menlo, monospace"
  text:
    base: { size: 1rem, line-height: 1.5rem }
    # ... xs to 4xl
  spacing: 0.25rem
  radius: { sm: 0.25rem, md: 0.375rem, lg: 0.5rem }
  shadow:
    sm: "0 1px 2px 0 oklch(0 0 0 / 0.05)"
---
```

### Recommendation

Option 1: tokens in the front matter of `design/system/design-system.md`, guidance in its body. It fits the documentation standard, `doc.get` and `doc.save` read and save it as they are, prime selects it by its topics, and one new check rule makes a bad token refuse a save. Option 3 is close; it loses because `doc.get` does not parse the block and the tokens clutter the rendered guidance. Code gets the tokens as a rendered `@theme` block, and Figma or Style Dictionary as a DTCG export, both generated from the document.

## Previewing it live without building the project

The preview shows the project's tokens on sample components, light and dark side by side, and follows a change to a token at once. The components are compiled with the dashboard, because its container holds no compiler and nothing of the project (§ Today).

### The fixed vocabulary

Tailwind generated the dashboard's utilities at its build, from the class names in its sources, and a component compiled with the dashboard can use only those. So the gallery is written against the fixed vocabulary of [The token schema](#the-token-schema), its roles and steps, and a stylesheet of its own maps each to a variable the preview sets from the project's tokens, as `layout.css` maps `--t-*` today:

```css
@import 'tailwindcss' source(none);
@source './gallery';
@theme inline {
  --color-primary: var(--ds-color-primary);
  --text-base: var(--ds-text-base);
  --text-base--line-height: var(--ds-text-base--line-height);
  --spacing: var(--ds-spacing);
  --radius-md: var(--ds-radius-md);
  --font-sans: var(--ds-font-sans);
}
```

`source(none)` and `@source` keep that stylesheet to the gallery's own classes. That is general knowledge of Tailwind v4's source detection; how `--spacing` and `--shadow-*` behave through `inline` is to check in the build story. A token change is then a change of `--ds-*` values, with no compile.

- **Tokens outside the vocabulary** are not on components. Extra colours are swatches with their contrast on `background`; extra steps, radii, and shadows are specimens drawn to size.
- **Fonts.** A font token is a stack, shown as the operator's browser resolves it; a family the browser lacks falls back, and the page says so. A project's own font files are in its sub-project, which the dashboard cannot read. Recommended: system stacks, and Google Fonts loaded into the preview by family name only when the document asks for it, since that sends a request from the operator's browser to Google.

### The gallery

Generic components, written once in flaiover: buttons (primary, secondary, ghost, danger, disabled), inputs (text, select, checkbox, radio, textarea, an error state), a form, a card, an alert per state, a table, navigation (a bar, tabs, breadcrumbs), and type specimens with headings over body text. They are not the project's components; they show its tokens, not its markup.

### Options

| Option | Isolation from the dashboard's theme | A token change reaches it | What the container runs | Light and dark side by side | Cost to build |
|--------|--------------------------------------|---------------------------|-------------------------|-----------------------------|---------------|
| 1. A same-origin iframe; the page writes the gallery's stylesheet and a `<style>` of `--ds-*` values into it | Complete: a new document inherits nothing, provided it does not load `layout.css` (below) | By rewriting one `<style>`: a restyle, no compile | Nothing new | Two panels in one frame, each setting its theme's `--ds-*` on its wrapper | Small to medium |
| 2. A shadow root with the same stylesheet | Partial: inherited properties (`color`, `font-family`, `line-height`) and the dashboard's custom properties cross into it and are reset on `:host`; ancestor selectors such as `[data-theme='dark'] *` do not cross | As 1 | Nothing new | Two wrappers | Small, then pitfalls: `@property` rules in a shadow root are ignored, so Tailwind utilities that rely on registered `--tw-*` initial values (shadows, rings, transforms) misbehave unless the document registers them; `@font-face` in a shadow root is not applied in Chromium |
| 3. `@tailwindcss/browser` in a frame, the tokens as `@theme` in a `<style type="text/tailwindcss">` | As 1 | By a recompile, in the browser, of the classes in the frame | A new runtime asset, bundled in the image rather than fetched from a CDN, recorded in `design/tech` | As 1 | Small; any class works, so markup needs no fixed vocabulary |
| 4. flai runs the project's own Vite and Tailwind build on the host and serves the result | Complete | A build per change | Nothing in the container; on the host, the project's toolchain and dependencies, behind a host action the operator enables (ADR-0029) | The project's own way | Large: a host action, a way to carry built assets over the channel, a build per change |

Tailwind's documentation says the browser build "is designed for development purposes only, and is not intended for production" ([Play CDN](https://tailwindcss.com/docs/installation/play-cdn), read 2026-10-13). The shadow DOM facts of option 2 are general knowledge of current browsers and Tailwind v4, not checked here.

A frame at a dashboard route boots the whole app under the root `flaiover/src/routes/+layout.svelte`, which imports `layout.css`, and `flaiover/src/app.html` stamps `data-theme` on the frame's `<html>`. Two ways keep the dashboard's theme out of option 1:

- **A blank frame the page fills.** An `about:blank` frame is same-origin with its page. The page imports the gallery's stylesheet with Vite's `?inline`, so it reaches the dashboard's document as a string and never as a stylesheet, writes it into the frame, and mounts the gallery into the frame's body with Svelte's `mount`. No route and no second app. Whether `@tailwindcss/vite` compiles an `?inline` import, and whether Svelte's event delegation works across the frame, is to check.
- **A route outside the root layout.** Move the root layout's header and `layout.css` import into a route group, `(app)/+layout.svelte`, and put the frame's page in a group of its own with its own stylesheet. Groups do not change URLs, but every route folder moves, and the frame boots a second copy of the app. Its request is authenticated like any page (`decide` in `flaiover/src/lib/server/auth.ts`), with the `SameSite=Lax` session cookie a same-origin frame sends.

Either way the page sizes the frame to its content, which a same-origin page can read. Token values reach the frame only after the parse `flai check` uses, so none carries `url()` or a stray `}` into the stylesheet. The frame follows the page's own state, so a change shows as the operator types; the document is saved through `doc.save` when the operator saves.

### The project's own components

Not now. They need the project's Svelte compiler and Tailwind build over its sources, which only the host has: option 4, or a component workshop such as Storybook run on the host and reached through flai. Either needs a host action, a way to carry built assets over a channel of named methods (ADR-0029), and a build per change. A cheaper step could come later: HTML snippets with Tailwind classes kept in the guidance, rendered in a frame with option 3.

### Recommendation

Option 1, a same-origin blank frame the page fills: the gallery compiled with the dashboard against the fixed vocabulary, its stylesheet built apart and written into the frame, the project's tokens as `--ds-*` values in one `<style>`, light and dark as two panels. It is isolated completely, adds nothing to the container, and follows a change with no compile. If the `?inline` stylesheet does not hold, use a route outside the root layout instead. Option 3 is the later way to preview markup the vocabulary cannot express, and option 4 the way to preview the project's own components, only if the operator asks for that.
