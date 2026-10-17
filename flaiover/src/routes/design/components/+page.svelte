<script lang="ts">
	// /design/components: the project's Svelte components (S-0439), flai's components.list and
	// components.status through /api/components. The catalogue is listed by sub-project and folder,
	// with a filter on name and path; a sub-project flai lists no components for is shown with its
	// kind. The component picked, kept in the URL as ?p=<sub-project>&c=<path>, is framed from its
	// sub-project's dev server, through the harness that dev server serves, in light, dark, or both
	// side by side, at 375, 768, or 1280 pixels or the panel's own width. A stopped dev server is
	// started from the page, and a failed start shows flai's message with Vite's last lines. The
	// page tells the operator to turn the components host action on while it is off, and to upgrade
	// flai on the host when it is too old to answer. While open, the page asks again every minute,
	// which keeps the host from stopping the dev server as idle and finds a component file added or
	// removed, and on a change to the manifest, whose components key says which files are components.
	// Beside the frame (S-0440), the props panel sets the component's props, which are sent to the
	// harness whole on each change and kept in the URL as ?props=, base64url JSON, so that a link
	// renders the same; the event log lists the calls of its function props. The frames show the
	// tried tokens (S-0441): the Design page's unsaved edits, kept in this browser as the project's
	// draft and followed from another tab by the storage event, laid over the saved tokens. The
	// source picker (S-0442), kept in the URL as ?story=, reads the catalogue and the dev servers
	// from a story's worktree instead of main's, and lists first, marked, the components its branch
	// changes, from its stream.diff. Beside the panel (S-0443), the component's examples, its
	// examples file's named sets of props: choosing one loads its props, theme, and width into the
	// panel, and Save as example saves the panel as one, on main or on the story the page reads.
	import type {
		ComponentExample,
		ComponentProp,
		ComponentsCatalogue,
		ComponentsProject,
		ComponentsStatus,
		DesignSystem,
		DesignTheme
	} from '$lib/server/repo';
	import type { Called, Props, Theme } from '$lib/explorer/protocol';
	import {
		PROPS_PARAM,
		fromQuery,
		functions,
		missing,
		parse,
		toQuery,
		toTexts,
		toValues
	} from '$lib/explorer/controls';
	import { api } from '$lib/api';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { follow } from '$lib/events';
	import { themeState } from '$lib/theme.svelte';
	import ComponentList from '$lib/components/explorer/ComponentList.svelte';
	import ComponentFrame from '$lib/components/explorer/ComponentFrame.svelte';
	import PropsPanel from '$lib/components/explorer/PropsPanel.svelte';
	import EventLog, { logged, type Logged } from '$lib/components/explorer/EventLog.svelte';
	import { projectState } from '$lib/project.svelte';
	import { clearDraft, loadDraft, onDraft, type Edits } from '$lib/design/draft';
	import { editedThemes } from '$lib/design/tokens';
	import { stylesheetDeclarations } from '$lib/explorer/tokens';
	import SourcePicker, {
		openStories,
		type SourceStory
	} from '$lib/components/explorer/SourcePicker.svelte';
	import { changedComponents, changedFirst, componentFile, diffPaths } from '$lib/explorer/changed';
	import ExampleList, {
		exampleFrom,
		type SaveFailure
	} from '$lib/components/explorer/ExampleList.svelte';

	const ENABLE = 'flai serve enable components';
	// Asking at least this often keeps the host from stopping a dev server as idle.
	const POLL_MS = 60_000;
	const WIDTHS = [375, 768, 1280] as const;
	const button = 'rounded border px-3 py-1 text-sm disabled:opacity-50';
	const toggle = (on: boolean) => [button, on ? 'border-primary bg-primary-soft' : 'border-line'];

	let catalogue = $state<ComponentsCatalogue | null>(null);
	let status = $state<ComponentsStatus | null>(null);
	let error = $state<string | null>(null);
	// The host's flai answered 501: its message, which says what it cannot do.
	let tooOld = $state<string | null>(null);
	// The host action was refused (403): what the operator runs to turn it on.
	let refused = $state<string | null>(null);
	let startError = $state<string | null>(null);
	let busy = $state(false);

	let pickedProject = $state(page.url.searchParams.get('p') ?? '');
	let pickedPath = $state(page.url.searchParams.get('c') ?? '');
	let both = $state(false);
	let single = $state<Theme>(themeState.dark ? 'dark' : 'light');
	let width = $state<number | null>(null);

	// The props the URL gave as the page opened, for the component it names, until another is picked.
	let linked: { project: string; path: string; values: Props } | null = $state.raw({
		project: page.url.searchParams.get('p') ?? '',
		path: page.url.searchParams.get('c') ?? '',
		values: fromQuery(page.url.searchParams.get(PROPS_PARAM))
	});
	// The props panel's texts, for the component they were written for, and the props last sent
	// from them: text that does not parse holds back every prop until it does.
	let edits: {
		project: string;
		path: string;
		texts: Record<string, string>;
		values: Props;
	} | null = $state.raw(null);
	let events = $state.raw<Logged[]>([]);

	// The story whose worktree the page reads, undefined for main's; the stories it may pick; and the
	// paths its branch changes, null until its diff is read.
	let story = $state(page.url.searchParams.get('story') || undefined);
	let stories = $state.raw<SourceStory[]>([]);
	let changedPaths = $state.raw<string[] | null>(null);
	let changedError = $state<string | null>(null);

	const changed = $derived(
		catalogue && changedPaths ? changedComponents(catalogue, changedPaths) : new Set<string>()
	);
	// The catalogue's sub-projects, each with the components the story's branch changes first.
	const projects = $derived(
		(catalogue?.projects ?? []).map((p) => ({
			...p,
			components: changedFirst(p.components, changed)
		}))
	);
	const changedList = $derived(
		story && catalogue && changedPaths
			? projects.flatMap((p) =>
					p.components
						.filter((c) => changed.has(componentFile(c)))
						.map((c) => ({ project: p.name, component: c }))
				)
			: null
	);
	// The sub-project the dev server controls are for: the picked one, else the first supported.
	const project = $derived<ComponentsProject | undefined>(
		projects.find((p) => p.name === pickedProject && p.supported) ??
			projects.find((p) => p.supported)
	);
	const component = $derived(project?.components.find((c) => c.path === pickedPath));
	const server = $derived(
		status?.servers.find((s) => s.project === project?.name && s.story === story)
	);
	const off = $derived(refused ?? (status && !status.enabled ? status.enable || ENABLE : null));
	const themes = $derived<Theme[]>(both ? ['light', 'dark'] : [single]);
	const picked = $derived(
		component && project ? { project: project.name, path: component.path } : null
	);

	const forPicked = (o: { project: string; path: string } | null) =>
		!!o && o.project === project?.name && o.path === component?.path;
	const propsOf = $derived(component?.props ?? []);
	// Without edits, the URL's props for the component it names, else the defaults.
	const start = $derived(toTexts(propsOf, forPicked(linked) ? linked!.values : {}));
	const texts = $derived(forPicked(edits) ? edits!.texts : start);
	const values = $derived(forPicked(edits) ? edits!.values : toValues(propsOf, start));
	const logs = $derived(functions(propsOf));
	const snippets = $derived(propsOf.filter((p) => p.kind === 'snippet').map((p) => p.name));
	const unset = $derived(propsOf.filter((p) => missing(p, texts[p.name])).map((p) => p.name));

	/** Shows the one theme. */
	function one(theme: Theme) {
		both = false;
		single = theme;
	}

	const errorText = (e: unknown) => (e instanceof Error ? e.message : String(e));

	/** Lists the stories in progress or in review, from the board; none when it cannot be read. */
	async function loadStories() {
		try {
			const r = await api('/api/board');
			stories = r.ok ? openStories((await r.json()).columns ?? {}) : [];
		} catch {
			stories = [];
		}
	}

	/** Reads the paths the story's branch changes, from its stream.diff, if it is still picked. */
	async function loadChanged(from: string) {
		try {
			const r = await api(`/api/items/${encodeURIComponent(from)}/diff`);
			const body = await r.json().catch(() => ({}));
			if (from !== story) return;
			if (!r.ok) throw new Error(body.error ?? r.statusText);
			changedPaths = diffPaths(body);
			changedError = null;
		} catch (e) {
			if (from === story) changedError = errorText(e);
		}
	}

	/** Reads the catalogue and the dev servers; a failed read keeps what was last shown. */
	async function load() {
		const from = story;
		void loadStories();
		if (from) void loadChanged(from);
		try {
			const r = await api(
				from ? `/api/components?story=${encodeURIComponent(from)}` : '/api/components'
			);
			const body = await r.json().catch(() => ({}));
			// an answer for the source picked before is not this one's
			if (from !== story) return;
			if (r.status === 501) {
				tooOld = body.error ?? r.statusText;
				return;
			}
			if (r.status === 403) {
				refused = body.enable ?? ENABLE;
				return;
			}
			if (!r.ok) throw new Error(body.error ?? r.statusText);
			catalogue = body.catalogue as ComponentsCatalogue;
			status = body.status as ComponentsStatus;
			tooOld = refused = error = null;
		} catch (e) {
			if (from === story) error = `Could not read the components: ${errorText(e)}`;
		}
	}

	/** Reads the source chosen, main's or a story's worktree, and keeps it in the URL. */
	function choose(next: string | undefined) {
		if (next === story) return;
		story = next;
		catalogue = status = changedPaths = null;
		error = changedError = startError = null;
		const kept = Object.fromEntries(page.url.searchParams);
		delete kept.story;
		const q = new URLSearchParams({ ...kept, ...(next ? { story: next } : {}) }).toString();
		// eslint-disable-next-line svelte/no-navigation-without-resolve -- the path is resolve()d; the rule does not follow the query added to it
		void goto(resolve('/design/components') + (q ? `?${q}` : ''), {
			replaceState: true,
			keepFocus: true,
			noScroll: true
		});
		void load();
	}

	/** Starts or stops the sub-project's dev server, taking the status the host answers. */
	async function control(action: 'start' | 'stop', name: string) {
		busy = true;
		startError = null;
		try {
			const r = await api(`/api/components/${action}`, {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ project: name, ...(story ? { story } : {}) })
			});
			const body = await r.json().catch(() => ({}));
			if (r.ok) status = body as ComponentsStatus;
			else if (r.status === 501) tooOld = body.error ?? r.statusText;
			else if (r.status === 403) refused = body.enable ?? ENABLE;
			else {
				startError = body.error ?? r.statusText;
				// the dev server's last error is in the status once the start has failed
				await load();
			}
		} catch (e) {
			startError = errorText(e);
		} finally {
			busy = false;
		}
	}

	/** Keeps the component picked and its props in the URL, so that a reload or a link frames the
	 * same; no props is no PROPS_PARAM. */
	function keep(props: Props) {
		const kept = Object.fromEntries(page.url.searchParams);
		delete kept[PROPS_PARAM];
		delete kept.story;
		const q = toQuery(props);
		const query = new URLSearchParams({
			...kept,
			p: pickedProject,
			c: pickedPath,
			...(story ? { story } : {}),
			...(q ? { [PROPS_PARAM]: q } : {})
		});
		// eslint-disable-next-line svelte/no-navigation-without-resolve -- the path is resolve()d; the rule does not follow the query added to it
		void goto(resolve('/design/components') + `?${query}`, {
			replaceState: true,
			keepFocus: true,
			noScroll: true
		});
	}

	/** Frames the component with its default props, and keeps it in the URL. */
	function pick(name: string, path: string) {
		if (name !== pickedProject || path !== pickedPath) {
			linked = edits = null;
			events = [];
		}
		pickedProject = name;
		pickedPath = path;
		keep(values);
	}

	/** Takes the panel's texts; once every one parses, its props are sent and kept in the URL. */
	function edit(props: ComponentProp[], next: Record<string, string>) {
		if (!picked) return;
		const valid = props.every((p) => parse(p, next[p.name]).ok);
		edits = { ...picked, texts: next, values: valid ? toValues(props, next) : values };
		if (valid) keep(edits.values);
	}

	/** Loads the example's props into the panel, and its theme and width where it gives them. */
	function chooseExample(example: ComponentExample) {
		if (!picked) return;
		const next = toTexts(propsOf, example.props);
		edits = { ...picked, texts: next, values: toValues(propsOf, next) };
		keep(edits.values);
		if (example.theme) one(example.theme);
		if (example.width) width = example.width;
	}

	/** Saves the panel as the example name in the picked component's examples file, on the story
	 * the page reads or on main, and lists the file's examples as saved; why not when it fails. */
	async function saveExample(name: string): Promise<SaveFailure | null> {
		const target = picked;
		if (!target) return { error: 'Pick a component to save an example of it.' };
		try {
			const r = await api('/api/components/examples', {
				method: 'PUT',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({
					project: target.project,
					component: target.path,
					example: exampleFrom(name, values, themes, width),
					...(story ? { story } : {})
				})
			});
			const body = await r.json().catch(() => ({}));
			if (!r.ok)
				return {
					error: body.error ?? r.statusText,
					findings: body.findings,
					conflict: r.status === 409
				};
			const saved = catalogue?.projects
				.find((p) => p.name === target.project)
				?.components.find((c) => c.path === target.path);
			if (saved) saved.examples = body.examples as ComponentExample[];
			return null;
		} catch (e) {
			return { error: errorText(e) };
		}
	}

	function called(m: Called) {
		events = logged(events, { at: Date.now(), name: m.name, args: m.args });
	}

	// ---- Tried tokens (S-0441) ----
	// The saved tokens are read only while there is a draft to lay over them, again on a document
	// change, which a save of the token file is, and forgotten with the draft, so that a later draft
	// is laid over what is saved by then. The frames are sent the generated stylesheet's custom
	// properties for the tried themes, every one, light and dark.

	// The layout remounts the page on a change of project, so its draft is this mount's.
	const draftProject = projectState.current;
	let draft = $state.raw<Edits>({});
	// The saved themes, and the token file's text, which says which tokens dark overrides.
	let saved: { themes: DesignTheme[]; raw: string | null } | null = $state.raw(null);
	let savedError = $state<string | null>(null);
	let reading = false;

	const held = (o: object) => Object.keys(o).length > 0;
	const tried = $derived.by(() => {
		if (!saved || !held(draft)) return null;
		const { light, dark } = stylesheetDeclarations(editedThemes(saved.themes, draft, saved.raw));
		return held(light) || held(dark) ? { light, dark } : null;
	});

	/** The token file's text, for a resolver document's dark context; null when it is not read,
	 * and then dark takes only its own edits, as on the Design page before the file is read. */
	async function readRaw(): Promise<string | null> {
		const r = await api('/api/design/tokens');
		if (!r.ok) return null;
		const body = await r.json().catch(() => ({}));
		return typeof body.content === 'string' ? body.content : null;
	}

	/** Reads the saved tokens the draft is laid over; a failed read keeps what was last read. */
	async function readSaved() {
		if (reading) return;
		reading = true;
		try {
			const r = await api('/api/design');
			const body = await r.json().catch(() => ({}));
			if (!r.ok) throw new Error(body.error ?? r.statusText);
			const design = body as DesignSystem;
			// a plain token file's one value serves both themes, so a light edit reaches dark
			const raw = design.resolver ? await readRaw() : '{}';
			if (held(draft)) saved = { themes: design.themes, raw };
			savedError = null;
		} catch (e) {
			savedError = `Could not read the saved tokens to lay the tried ones over: ${errorText(e)}`;
		} finally {
			reading = false;
		}
	}

	/** Takes the project's draft, reading the saved tokens when there are none to lay it over. */
	function takeDraft(edits: Edits) {
		draft = edits;
		if (!held(edits)) saved = savedError = null;
		else if (!saved) void readSaved();
	}

	/** Drops the draft, here and on the Design page's next reload, and the frames' tried tokens. */
	function dropDraft() {
		clearDraft(draftProject);
		takeDraft({});
	}

	onMount(() => {
		takeDraft(loadDraft(draftProject));
		return onDraft(draftProject, takeDraft);
	});

	onMount(() => {
		void load();
		const timer = setInterval(() => void load(), POLL_MS);
		// flai serve watches the manifest but not a sub-project's sources: a new component shows at the next ask
		const stop = follow(['document'], () => {
			void load();
			if (held(draft)) void readSaved();
		});
		return () => {
			clearInterval(timer);
			stop();
		};
	});
</script>

{#snippet offNotice(command: string)}
	<div
		class="rounded border border-warn bg-warn-soft p-3 text-sm text-warn"
		role="status"
		data-testid="components-off"
	>
		<p>The components host action is off.</p>
		<p>Run <code class="font-mono">{command}</code> on the host, in the project.</p>
	</div>
{/snippet}

<svelte:head><title>Components · flaiover</title></svelte:head>

<h1 class="text-2xl font-semibold">Components</h1>
<p class="mt-1 text-sm text-muted">
	The project's Svelte components, each framed from its sub-project's dev server. The dev server
	listens on 127.0.0.1, so the frame shows only in a browser on the host.
</p>
<div class="mt-4">
	<SourcePicker
		{stories}
		{story}
		onchoose={choose}
		changed={changedList}
		{changedError}
		{picked}
		onpick={pick}
	/>
</div>
{#if error}
	<p class="mt-4 rounded border border-danger bg-danger-soft p-3 text-sm text-danger" role="alert">
		{error}
	</p>
{/if}
{#if tooOld}
	<div
		class="mt-4 rounded border border-warn bg-warn-soft p-3 text-sm text-warn"
		role="status"
		data-testid="components-too-old"
	>
		<p>This flai cannot show the components yet; upgrade flai on the host.</p>
		<p class="mt-1 text-xs">{tooOld}</p>
	</div>
{:else if catalogue === null}
	{#if refused}
		<div class="mt-4">{@render offNotice(refused)}</div>
	{:else if !error}
		<p class="mt-4 text-sm text-muted">Loading…</p>
	{/if}
{:else}
	<div class="mt-6 grid grid-cols-1 gap-6 lg:grid-cols-[16rem_minmax(0,1fr)]">
		<nav aria-label="Components">
			<ComponentList {projects} {picked} onpick={pick} {changed} />
		</nav>
		<section class="min-w-0" aria-labelledby="components-frame-heading">
			<h2 id="components-frame-heading" class="text-lg font-semibold">
				{component?.name ?? 'Preview'}
			</h2>
			{#if component}
				<p class="font-mono text-xs text-muted">{component.path}</p>
			{/if}
			<div class="mt-2 flex flex-wrap items-center gap-4">
				<div class="flex gap-1" role="group" aria-label="Theme">
					<button
						type="button"
						class={toggle(!both && single === 'light')}
						aria-pressed={!both && single === 'light'}
						onclick={() => one('light')}>Light</button
					>
					<button
						type="button"
						class={toggle(!both && single === 'dark')}
						aria-pressed={!both && single === 'dark'}
						onclick={() => one('dark')}>Dark</button
					>
					<button
						type="button"
						class={toggle(both)}
						aria-pressed={both}
						onclick={() => (both = true)}>Both</button
					>
				</div>
				<div class="flex gap-1" role="group" aria-label="Width">
					{#each WIDTHS as w (w)}
						<button
							type="button"
							class={toggle(width === w)}
							aria-pressed={width === w}
							aria-label="{w} pixels wide"
							onclick={() => (width = w)}>{w}</button
						>
					{/each}
					<button
						type="button"
						class={toggle(width === null)}
						aria-pressed={width === null}
						aria-label="The panel's width"
						onclick={() => (width = null)}>Fill</button
					>
				</div>
			</div>
			{#if tried}
				<div
					class="mt-2 flex flex-wrap items-center gap-3 text-sm"
					role="status"
					data-testid="components-tried"
				>
					<span class="text-muted">
						Showing tried tokens: the Design page's unsaved edits, kept in this browser.
					</span>
					<button type="button" class={[button, 'border-line']} onclick={dropDraft}
						>Drop tried tokens</button
					>
				</div>
			{/if}
			{#if savedError}
				<p class="mt-2 text-xs text-warn" data-testid="components-tried-error">{savedError}</p>
			{/if}
			{#if component?.error}
				<p class="mt-2 text-xs text-warn" data-testid="component-props-error">
					Its props could not be read: {component.error}
				</p>
			{/if}
			<div class="mt-4" data-testid="components-server">
				{#if !project}
					<p class="text-sm text-muted">No sub-project is Svelte, so there is nothing to frame.</p>
				{:else if off}
					{@render offNotice(off)}
				{:else if server?.running && server.url}
					<div class="flex flex-wrap items-center gap-3 text-sm">
						<span class="text-muted">{project.name}'s dev server runs at {server.url}.</span>
						<button
							type="button"
							class={[button, 'border-line']}
							disabled={busy}
							onclick={() => control('stop', project.name)}>Stop</button
						>
					</div>
				{:else if busy || server?.starting}
					<p class="text-sm text-muted" role="status" data-testid="components-starting">
						Starting {project.name}'s dev server…
					</p>
				{:else}
					<div class="flex flex-wrap items-center gap-3 text-sm">
						<span class="text-muted">{project.name}'s dev server is not running.</span>
						<button
							type="button"
							class="rounded bg-primary px-3 py-1 text-sm text-on-primary"
							onclick={() => control('start', project.name)}>Start</button
						>
					</div>
					{#if !startError && server?.last_error}
						<p class="mt-2 text-xs text-muted">Its last dev server failed:</p>
						<pre
							class="mt-1 max-h-64 overflow-auto rounded border border-line bg-ground p-2 font-mono text-xs whitespace-pre-wrap text-ink"
							data-testid="components-last-error">{server.last_error}</pre>
					{/if}
				{/if}
				{#if startError}
					<pre
						class="mt-2 max-h-64 overflow-auto rounded border border-danger bg-danger-soft p-2 font-mono text-xs whitespace-pre-wrap text-danger"
						role="alert"
						data-testid="components-start-error">{startError}</pre>
				{/if}
			</div>
			<div class="mt-4 flex flex-col gap-6 xl:flex-row">
				<div class="min-w-0 flex-1">
					{#if component && server?.running && server.url}
						<div class="flex gap-4 overflow-x-auto" data-testid="components-frames">
							{#each themes as theme, i (i)}
								<div class={[!width && 'min-w-0 flex-1']}>
									<ComponentFrame
										server={server.url}
										component={component.path}
										{theme}
										{width}
										props={values}
										functions={logs}
										{snippets}
										missing={unset}
										tokens={tried}
										onevent={called}
									/>
								</div>
							{/each}
						</div>
					{:else if project && !component}
						<p class="text-sm text-muted">Pick a component to frame it.</p>
					{/if}
				</div>
				{#if component}
					<div class="flex flex-col gap-6 xl:w-80 xl:shrink-0">
						{#key `${project?.name}:${component.path}`}
							<ExampleList
								examples={component.examples ?? []}
								path={component.examples_path}
								{story}
								onchoose={chooseExample}
								onsave={saveExample}
							/>
						{/key}
						<PropsPanel
							props={component.props}
							{texts}
							onchange={(next) => edit(component.props, next)}
						/>
						<EventLog {events} onclear={() => (events = [])} />
					</div>
				{/if}
			</div>
		</section>
	</div>
{/if}
