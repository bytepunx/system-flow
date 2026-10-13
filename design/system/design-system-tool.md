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

## Working on it together

The operator decides the design system; agents propose changes and apply them. Each way of working is weighed for how a decision is reviewed, recorded, and undone, what is new to build, and how it fits the board.

### Options

| Option | Reviewed | Recorded | Undone | New to build | Fit with the board |
|--------|----------|----------|--------|--------------|--------------------|
| 1. The operator edits tokens on the Design page, the preview following each change, and saves through `doc.save` on main | By the operator as they type, with contrast beside each token; no second reader | A commit of the document on main, as every dashboard save is ([ADR-0023](../adrs/0023-documents-are-saved-through-flai.md)) | Another save, or a revert of the commit | A token editor and a front matter round trip on the page (S-0383); the stylesheet written with the save (below) | Outside it: no story, no WIP slot, no close-out |
| 2. A proposal is a story whose branch changes the document; the Design page previews the branch's tokens beside main's, and the operator tweaks them on the branch | On the Design page beside main's, then on the review page's diff (`stream.diff`) | The story, its narrative, its branch's commits, its acceptance | Before acceptance, send it back or cancel it; after, another story or a revert | Reading the document from a story's worktree and listing open proposals (S-0384); saving on the branch is S-0378's | It is the board: claims, WIP limit, close-out, acceptance |
| 3. A thread on a heading of the document | In the thread | Under `wip/threads`, anchored to the document and heading | Resolve it | Nothing in flai; the Design page reuses the document page's threads | Beside it: no agent is started for it (below) |
| 4. Draft variants outside git, under `.flai-cache` or `wip/`, which the operator switches between | On the page only | Under `.flai-cache`, on one host and in no history; under `wip/`, a new kind of file with its own schema and check | Delete the variant | A store, host methods to read and write it, a picker, and a way to promote a variant to the document | None: a variant is no item, and prime does not give it to an agent |

- **S-0378**, in progress on `story/S-0378`, gives option 2 its reads and writes. `flai doc show` and `flai doc save` take `--story`: the file resolves in the story's worktree, the check runs there, and the save commits that one path on the story's branch; a wip path, or a story with no worktree, is refused (T-1514). The dashboard's `/api/docs/edit` and `PUT /api/docs/file` pass `story` (T-1517), the editor opens on a story's branch (T-1518), and T-1521 records it as an ADR refining ADR-0023. Option 2 depends on S-0378's acceptance. A story's worktree exists while it is open, so the operator tweaks a proposal while it is in progress or in review.
- **The stylesheet goes with every token save.** `doc.save` refuses any finding the new content introduces, on any path (`flai/internal/docedit/docedit.go`). S-0385 as drafted adds `design.stale`, a finding when the stylesheet differs from the tokens. A token save that leaves the stylesheet as it was introduces it, so every save from the Design page, on main or on a branch, would be refused. So whenever flai writes the document's tokens, it writes the stylesheet the manifest names in the same commit, and `design.stale` is left to hand edits. On main that is a file of a sub-project changed by a dashboard save, outside any story: a new kind of write, for the ADR that refines ADR-0023.
- **Threads.** `thread.new` takes a document and a heading (`flai/internal/hostapi/writes.go`), and `/docs/<path>` shows a document's threads with a composer per heading (`Threads.svelte`, [ADR-0020](../adrs/0020-files-plus-mcp.md)). Every agent's inbox counts such a thread as awaiting it while the last entry is not its own (`flai/internal/inbox/inbox.go`). But flai serve starts agents for stories, and restarts a story's agent only for a thread on that story (`asking` in `flai/internal/serve/agents.go`). The orchestrator answers threads while its permission is on, and edits no file. A thread on the document discusses and records; it does not get the document changed. A thread on a proposal story reaches that story's agent.
- **Variants.** Option 1 already lets the operator try values without saving. Keeping the unsaved values in the browser across a reload gives most of what a store would, with nothing on the host.

### Asking for a change from the page

A form on the Design page makes a story as the board's new-item form does: `item.new` with a title, a body, `topics`, and `touches`, made in backlog and moved to ready (`flaiover/src/lib/components/NewItemForm.svelte`, `movesTo` in `flaiover/src/lib/lanes.ts`). The form writes the goal from the operator's words and fixed criteria: the tokens as asked, the stylesheet in step, the UI's tests passing. Its touches are the document and the stylesheet, and it may carry the values the operator tried on the page as a YAML block in the body, for the agent to start from.

The story meets the definition of ready as made: a goal, criteria as checkboxes, an open epic if it has one, and not a draft (`work-management.md`). It needs no finalizing. `draft` marks a story an agent wrote ([ADR-0074](../adrs/0074-work-items-carry-planning-data-a-story-s-draft-flag-an-epic-s-or-story-s-cost.md)), and `finalized` records who took a draft on ([ADR-0077](../adrs/0077-a-story-that-was-a-draft-records-who-finalized-it-and-when-and-the-dashboard.md)); the operator wrote this one. A design story the planner drafts waits for the operator's finalize, as any draft does.

### Two changes at once

A design story claims the document's path, and a second ready story that claims it is held until the first is accepted ([ADR-0046](../adrs/0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md)). That is right for one shared document, where two branches editing the same YAML conflict line by line. Two gaps remain:

- A story that widens its touches to the document while open holds only ready stories, never an open one. Two open branches can then change the document, and `flai stream sync`'s trial merge reports the conflict to both in their conversation.
- The stylesheet lies inside the UI sub-project. A story that claims the sub-project by name overlaps every design story, since a sub-project's path is a prefix of the stylesheet's, so design stories and UI stories wait on each other. Narrow touches on both sides keep them apart.

### Tokens and guidance

| | A token change | A guidance change |
|-|----------------|-------------------|
| What changes | Front matter values, and the stylesheet with them | The body's prose |
| Seen before saving | On the gallery | In the editor's preview |
| Reviewed | Beside main's on the Design page, and in the diff | In the diff |
| Reaches code | At once, through the stylesheet | Through the agents that read it |
| Saved | On main or on a proposal's branch (§ Recommendation) | On main through the editor, as any document, or on a story's branch |

### Recommendation

Options 1, 2, and 3; not 4.

- **1, for the operator's own decisions.** Try freely; save to main, flai writing the stylesheet with the tokens. It is the shortest path from a decision to the code, and the operator reviews every change on this board anyway. Its cost is a code file changed on main without a story's close-out. If every code change must pass through a story, token saves go on a proposal's branch instead, and option 1 is for trying, not saving: a choice for the operator.
- **2, for what an agent proposes or applies**, and for a change the operator asks for from the form. Review, history, and undo stay the board's.
- **3, for discussion.** The Design page lists the document's threads by heading and opens one on a heading. A request that needs an agent becomes a story, which the thread names.
- **Not 4.** Option 1's unsaved edits give what it would, and a store outside git is lost to history and unseen by agents.

## How agents build from it

An agent building UI needs the token names and the guidance on when to use each, more than the values. It should get them without being told, and its code should not drift from the document.

### Options

| Option | What the agent gets | What it needs | Cost | What it leaves open |
|--------|---------------------|---------------|------|---------------------|
| 1. Prime: the document's topics name the UI sub-project | A brief, title, first paragraph, and outline, with no tokens (§ Today); the whole file, front matter included, when the story names it | For the whole file on every UI story, a link in each story, or a rule in prime | The tokens are a few KB; with the guidance, likely under an eighth of the budget, 10 KB at 80 KB, which [ADR-0050](../adrs/0050-a-document-a-story-names-only-by-its-path-written-out-is-briefed-when-it-is.md) loads whole (an estimate) | A link depends on whoever writes the story |
| 2. A baseline convention, `design-system.md` | The rules: token utilities only, no raw colour or size, change the tokens on the branch rather than hard-code | A template file, and a rule that gives it only where a design system exists | One short convention in each UI story's pack | Advice; only the verifier holds an agent to it |
| 3. A stylesheet flai writes from the tokens, which the manifest names, and `design.stale` when it falls behind | Utilities for the tokens and no others of the project's: `bg-primary` exists, a colour outside the tokens has none | S-0385 | Small | Arbitrary values such as `bg-[#3b82f6]` still compile |
| 4. A lint that flags raw values in components | A finding on each arbitrary value or hex in `.svelte` and `.css` outside the generated file | A scanner of sub-project code, which `flai check` does not read today, or a rule in the project's own linter declared in the manifest's `tests` | Medium; false positives on layout values such as `grid-cols-[...]`, SVG, and chart palettes such as flaiover's `palette.ts` | Nothing 2 and 3 leave, at a price |

- **Which stories.** The document's `topics` already say which stories it reaches (§ Today). S-0386 as drafted adds a record of the sub-projects the document governs, written by `flai design init --project`. That is a second list beside `topics`, and the two can disagree. Use the topics.
- **The whole document.** Prime loads the design system document whole for every story its topics reach, as if the story named it, and briefs it above ADR-0050's threshold, where the convention tells the agent to read it with `doc_get`. A rule the agent must remember to follow is the failure the operator reports.
- **The convention, where it applies.** A baseline convention ships with `topics: [all]` ([conventions.md § Topics](conventions.md#topics)) and cannot name a project's UI. Prime gives `design-system.md` to exactly the stories the document reaches, and leaves it out where the project has no design system document. Its `roles` are `[story, verify]` ([ADR-0068](../adrs/0068-a-convention-s-roles-list-every-agent-that-reads-it-the-story-s-agent-included.md)), so that the verifier checks a UI diff against it.
- **What the convention says.** Build UI from the token utilities. No arbitrary value or hex for colour, type, spacing, radius, or shadow, and no colour in a `style` attribute. When the UI needs a value the tokens lack, change the document on the story's branch, widen the touches to it and the stylesheet, run `flai design css --write`, and say why in the guidance; the operator sees it on the Design page before accepting. Ask on a thread first when the change alters what a role means.
- **On a proposal's branch.** The agent edits the tokens, runs `flai design css --write`, and commits both. A forgotten write stops its close-out, provided `design.stale` names the document as well as the stylesheet: the close-out's check, scoped to the story, prints a finding on a file the story did not change as a note and goes on (`work-management.md`), so a branch that changed only the document would pass.
- **On main.** flai writes the stylesheet with every token save, and branches stop at close-out, so `design.stale` appears on main only after a hand edit there. The generated file should be whole, imported by the UI's own stylesheet, so that no hand-written CSS sits in it. How it selects the dark values, by a media query, a class, or an attribute as the project's `dark:` variant does, is S-0385's to settle.
- **The lint, later.** Option 3 removes most reasons for a raw value, and option 2 with the verifier catches the rest at review. Add a lint when `design/issues` shows raw values getting through, as a rule in the project's own linter rather than in flai, which reads no code today.

### Recommendation

Option 3 first, then 1 and 2 together; 4 deferred. The stylesheet makes the tokens the easy path and catches drift; prime then gives each UI story the whole document and the convention, chosen by one list, the document's topics. Build order: the stylesheet with S-0385, before the Design page saves anything; prime and the convention with S-0386, without its record of governed sub-projects.

## Recommendation

Build in this order:

1. **The document and its check** ([Where a project keeps its design system](#where-a-project-keeps-its-design-system)): `design/system/design-system.md`, tokens in its front matter, `design.tokens`, `flai design init`, `show`, and `css`, and `design.get`. S-0381.
2. **The stylesheet** ([How agents build from it](#how-agents-build-from-it)): `design.stylesheet`, `flai design css --write`, `design.stale` naming the document and the stylesheet, and flai writing the stylesheet whenever it writes the tokens. S-0385, before S-0383, whose save depends on it.
3. **The Design group and its preview** ([Previewing it live](#previewing-it-live-without-building-the-project)): a blank same-origin frame the page fills, the gallery against the fixed vocabulary, the tokens as `--ds-*` values. S-0382.
4. **Editing on the page** ([Working on it together](#working-on-it-together), option 1): the preview follows before saving, and the save writes the document and the stylesheet. S-0383.
5. **Proposals and threads** (options 2 and 3): the branch's tokens beside main's, the open proposals listed, the request form, the document's threads on the page. S-0384, after S-0378 is accepted.
6. **Agents' use** ([How agents build from it](#how-agents-build-from-it)): prime loads the document whole and gives the convention to the stories its topics reach; the baseline convention in the template and here. S-0386, without the governed sub-projects.
7. **Not now:** a store of variants, the project's own components, the browser build of Tailwind, and a lint of raw values.

What it does to the current ADRs: [ADR-0023](../adrs/0023-documents-are-saved-through-flai.md) is refined again, beyond S-0378's refinement, for a save that writes the stylesheet with the document. [ADR-0049](../adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md) and ADR-0050 are refined for the design system document loaded whole by topics. [ADR-0029](../adrs/0029-the-dashboard-reaches-a-project-only-through-flai-on-the-host.md), ADR-0046, and ADR-0074 stand.

Questions only the operator can answer, each with the recommended answer first:

1. **Format and name.** Tokens in the front matter of `design/system/design-system.md`, guidance in its body; or a fenced block in the body, a DTCG file, or the UI's stylesheet as the source.
2. **Colour syntax.** Six-digit hex and `oklch()`; or hex only; or any CSS colour, which flai cannot always compute contrast for.
3. **Roles.** The fixed set of colour roles the gallery is written against, extras allowed; or free names, which the gallery cannot style.
4. **Contrast.** Reported on the Design page and by `flai design`, not a check finding; or a finding, which refuses a save while the operator tries a colour.
5. **Scope.** One design system per project, a document per UI sub-project later; or one per UI sub-project now.
6. **Fonts.** System stacks, and Google Fonts loaded into the preview by family name only when the document asks; or system stacks only.
7. **Preview isolation.** A blank same-origin frame the page fills, a route outside the root layout as the fallback; or the browser build of Tailwind; or the project's own build on the host.
8. **Where the operator's token save lands.** On main, flai writing the stylesheet in the same commit; or only on a proposal's branch, so that every code change passes a story.
9. **Working together.** Edits on the page, proposals as stories, and threads on the document; or those and a store of variants.
10. **Agents' use.** The generated stylesheet with `design.stale`, the document whole and the convention for the stories its topics reach, the lint deferred; or briefs and the convention only; or the lint now.
