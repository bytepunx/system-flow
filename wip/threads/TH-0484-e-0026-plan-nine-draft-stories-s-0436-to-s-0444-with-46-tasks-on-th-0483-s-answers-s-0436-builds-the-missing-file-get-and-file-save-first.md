---
id: TH-0484
title: "E-0026 plan: nine draft stories, S-0436 to S-0444, with 46 tasks, on TH-0483's answers; S-0436 builds the missing file.get and file.save first"
anchor:
  path: wip/kanban/epics/E-0026-component-library.md
  item: E-0026
status: resolved
participants: [planner-E-0026, alex]
created: 2026-10-16T21:51:26Z
updated: 2026-10-16T22:09:08Z
---

# TH-0484 E-0026 plan: nine draft stories, S-0436 to S-0444, with 46 tasks, on TH-0483's answers; S-0436 builds the missing file.get and file.save first

On wip/kanban/epics/E-0026-component-library.md.

## Entries

### 2026-10-16T21:51:26Z planner-E-0026
I drafted nine stories for E-0026, all in the backlog as drafts. Each has its tasks, file-level touches, a forecast, and a cost of delay value from the epic's 600 USD a week. They follow TH-0483's answers, including examples under `design/components/`. Each story carries its own docs, and each one's `### Planning` notes say where every touch came from and why each figure stands. `flai check --strict` reports nothing.

## Stories and order

| Story | What | After | Forecast | flai's | Value a week |
|-------|------|-------|----------|--------|--------------|
| S-0436 | `file.get` and `file.save`, `flai file show` and `save`; a token save writes the stylesheet with it | none | 37m | 37m | 50.92 |
| S-0437 | The catalogue: ADR, `design/system/component-library.md`, the manifest's `components`, props read by the project's own `svelte/compiler`, `flai components list` and `config` | none | 43m | 43m | 59.17 |
| S-0438 | The project's own Vite dev server on 127.0.0.1, the harness, the protocol, the `components` host action, the idle stop | S-0437 | 1h20m | 55m | 110.09 |
| S-0439 | The Components page in the Design group: the list, the frame, light and dark, widths, start and errors | S-0438 | 37m | 37m | 50.92 |
| S-0440 | Props controls from each prop's kind, in-place updates, function props logged, props in the URL | S-0439 | 26m | 26m | 35.78 |
| S-0441 | Token changes reach the components: tried values from the Design page, saved ones by hot update | S-0436, S-0439 | 31m | 31m | 42.66 |
| S-0442 | A story's worktree as the source, with deps installed by lock file, and a link from the review page | S-0439 | 58m | 58m | 79.82 |
| S-0443 | Examples under `design/components/`, `components.examples` in `flai check`, saved from the page | S-0436, S-0440, S-0442 | 1h4m | 1h4m | 88.07 |
| S-0444 | Every flaiover component renders, with its setup, exclusions, and examples | S-0443 | 1h | 16m | 82.57 |

The total is 7h16m of agent time. flai plays the last one out to 2026-10-18T03:04Z, behind E-0021 to E-0025's stories in the pull order. Outcome 1 lands with S-0439, outcome 2 with S-0440, and outcome 3 with S-0439 for source and S-0441 for the design system. S-0442 to S-0444 complete them for a story's work, for components that need data, and for this repository.

## Tasks and layers

- **S-0436:** T-1849; then T-1850; then T-1851; then T-1852.
- **S-0437:** T-1853; then T-1854, T-1855; then T-1856; then T-1857.
- **S-0438:** T-1858; then T-1859, T-1860; then T-1861, T-1863; then T-1862; then T-1864.
- **S-0439:** T-1865, T-1866, T-1867; then T-1868; then T-1869.
- **S-0440:** T-1870, T-1871; then T-1872; then T-1873.
- **S-0441:** T-1874, T-1875, T-1876; then T-1877; then T-1878.
- **S-0442:** T-1879, T-1881, T-1882; then T-1880; then T-1883.
- **S-0443:** T-1884; then T-1885, T-1886, T-1894; then T-1887; then T-1888; then T-1889.
- **S-0444:** T-1890; then T-1891, T-1892; then T-1893.

No two tasks of one layer share a path.

## Assumptions

1. **The flai side of `file.get` and `file.save` does not exist.** S-0383's notes left it to a story E-0020's planner would draft, and none was. The Design page's token save answers 501 today, so S-0436 is worth pulling early on its own. It waits for nothing in this epic.
2. **ADR-0151's unsaved values are not kept.** It says the Design page keeps them across a reload, but the page keeps them in page state only. S-0441 builds it, since the Components page needs them.
3. **The operator's browser runs on the host**, as yours does at 127.0.0.1:4242, and `node` is installed there. Vite sends no header that refuses framing, and flaiover has no Content-Security-Policy, so the cross-origin frame needs nothing new; S-0438's manual run confirms it.
4. **`flai manifest set` reaches no sub-project key.** `CLAUDE.md` gives flai the manifest, so S-0437 adds `flai components config`, which S-0444 uses for flaiover.
5. **A story's worktree has no `node_modules`.** This repository installs them with `scripts/flaiover-install.sh` (I-0080). S-0442 has flai install by the lock file before a story's dev server starts.
6. **The two raised forecasts.** S-0438's merge of the project's Vite config with SvelteKit's plugin has no precedent. S-0444's one folder touch stands for some 20 to 40 files. The values are flai's shares by duration, so the raises move value to those two.
7. **The stamps.** The forecasts of S-0436 to S-0442 carry `by: system-flow`. The CLI wrote them under the host's `FLAI_AGENT`, and flai refused to restamp an unchanged block. The planner set them.

## Proposals

- **Pull S-0436 first (not applied).** It fixes the Design page's save today and is independent of the rest; the order is yours.
- **Not drafted:** React and Vue adapters (TH-0483 answer 2, later); a proxy through the channel for a dashboard off the host (TH-0483 answer 1's alternative); accessibility and screenshot checks per example, which the epic's "test UI components" points at; and the catalogue given to a story's agent by prime or MCP, so that it reuses components rather than writes new ones. Say which you want and I will draft them.

## Overlap with other epics

These stories share files with drafts of E-0021 to E-0025, so the holds run them one after the other:

- `flai/internal/hostapi/writes.go` and `hostapi.go`: S-0436, S-0438, and S-0442 with S-0393, S-0395, S-0397, S-0399, S-0417, S-0418, S-0420, S-0422, S-0431, and S-0433.
- `flai/cmd/serve_actions.go`: S-0438 and S-0442 with S-0392, S-0399, S-0417, S-0422, and S-0431.
- `flai/internal/manifest/manifest.go` and `settings.go`: S-0437 with S-0392, S-0399, S-0405, S-0417, S-0420, S-0424, S-0431, and S-0432.
- `flai/cmd/root.go`: S-0436 and S-0437 with S-0392, S-0398, S-0414, S-0421, S-0428, and S-0431.
- `flai/internal/check/check.go`: S-0443 with S-0406, S-0414, S-0417 to S-0419, S-0421, and S-0428.
- The site menu: S-0439 with S-0420. The review page: S-0442 with S-0394 and S-0395.

### 2026-10-16T22:09:08Z alex
Resolved.
