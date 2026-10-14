---
id: ADR-0149
title: "The Design page previews a project's tokens in a blank same-origin frame it fills with a gallery compiled with the dashboard against the fixed token vocabulary, with system font stacks and Google Fonts loaded only when the design system names them"
status: accepted
date: 2026-10-14
supersedes: []
superseded_by: []
topics: [design, dashboard, sveltekit]
---

# ADR-0149 The Design page previews a project's tokens in a blank same-origin frame it fills with a gallery compiled with the dashboard against the fixed token vocabulary, with system font stacks and Google Fonts loaded only when the design system names them

## Context

The Design page shows a project's tokens on sample components, light and dark side by side, and follows a change to a token at once. The dashboard's container holds no Tailwind or Svelte compiler and nothing of the project ([ADR-0029](0029-the-dashboard-reaches-a-project-only-through-flai-on-the-host.md), [ADR-0031](0031-the-dashboard-s-container-holds-nothing-of-the-project-a-port-and-two-secrets.md)), and the dashboard's own theme must not leak into the preview. The finding, [design-system-tool.md](../system/design-system-tool.md), weighed a same-origin frame, a shadow root, Tailwind's browser build, and the project's own build on the host, and recommended choices 6 and 7 below. On TH-0445, on 2026-10-14, the operator answered only choice 1, so these stand as recommended.

## Decision

The Design page previews a project's tokens in a blank same-origin frame it fills with a gallery compiled with the dashboard against the fixed token vocabulary, with system font stacks and Google Fonts loaded only when the design system names them.

1. **The frame.** An `about:blank` frame, same-origin with the page. The page imports the gallery's stylesheet with Vite's `?inline`, writes it into the frame, and mounts the gallery into the frame's body with Svelte's `mount`. The project's tokens are `--ds-*` values in one `<style>`; light and dark are two panels, each setting its theme's values on its wrapper. A token change rewrites that `<style>`, with no compile.
2. **The fallback.** If the `?inline` stylesheet or Svelte's events across the frame do not hold, the frame loads a route outside the root layout, in a route group with its own stylesheet.
3. **The gallery.** Generic components written once in flaiover against the fixed vocabulary of roles and steps ([ADR-0148](0148-a-design-token-s-colour-is-in-the-srgb-colour-space-with-its-hex-or-in-oklch.md)): buttons, inputs, a form, a card, an alert per state, a table, navigation, and type specimens. Extra tokens are swatches and specimens.
4. **Values.** Token values reach the frame only after the parse `design.tokens` uses, so none carries `url()` or a stray `}` into the stylesheet.
5. **Fonts.** A font token is a stack, shown as the operator's browser resolves it. Google Fonts are loaded into the preview by family name only when the design system names one, since that sends a request from the operator's browser to Google.

## Consequences

- Nothing new runs in the container, and ADR-0029 and ADR-0031 stand.
- The preview shows the tokens, not the project's own components or markup.
- A class outside the gallery's vocabulary cannot be previewed until a later decision adds Tailwind's browser build or the project's own build.
- S-0382 checks the `?inline` import and the frame's events before it relies on them.

## Alternatives considered

- **A shadow root.** Inherited properties cross into it, `@property` rules in it are ignored, and Chromium does not apply `@font-face` in it.
- **Tailwind's browser build in the frame.** Any class works, but its documentation says it is for development only, and it adds a runtime asset.
- **The project's own build on the host.** It needs a host action, a way to carry built assets over the channel, and a build per change.
- **System font stacks only.** No request to Google, but no way to judge a web font.
