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
	import type { ComponentsCatalogue, ComponentsProject, ComponentsStatus } from '$lib/server/repo';
	import type { Theme } from '$lib/explorer/protocol';
	import { api } from '$lib/api';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { follow } from '$lib/events';
	import { themeState } from '$lib/theme.svelte';
	import ComponentList from '$lib/components/explorer/ComponentList.svelte';
	import ComponentFrame from '$lib/components/explorer/ComponentFrame.svelte';

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

	const projects = $derived(catalogue?.projects ?? []);
	// The sub-project the dev server controls are for: the picked one, else the first supported.
	const project = $derived<ComponentsProject | undefined>(
		projects.find((p) => p.name === pickedProject && p.supported) ??
			projects.find((p) => p.supported)
	);
	const component = $derived(project?.components.find((c) => c.path === pickedPath));
	const server = $derived(status?.servers.find((s) => s.project === project?.name));
	const off = $derived(refused ?? (status && !status.enabled ? status.enable || ENABLE : null));
	const themes = $derived<Theme[]>(both ? ['light', 'dark'] : [single]);
	const picked = $derived(
		component && project ? { project: project.name, path: component.path } : null
	);

	/** Shows the one theme. */
	function one(theme: Theme) {
		both = false;
		single = theme;
	}

	const errorText = (e: unknown) => (e instanceof Error ? e.message : String(e));

	/** Reads the catalogue and the dev servers; a failed read keeps what was last shown. */
	async function load() {
		try {
			const r = await api('/api/components');
			const body = await r.json().catch(() => ({}));
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
			error = `Could not read the components: ${errorText(e)}`;
		}
	}

	/** Starts or stops the sub-project's dev server, taking the status the host answers. */
	async function control(action: 'start' | 'stop', name: string) {
		busy = true;
		startError = null;
		try {
			const r = await api(`/api/components/${action}`, {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ project: name })
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

	/** Frames the component, and keeps it in the URL so that a reload frames it again. */
	function pick(name: string, path: string) {
		pickedProject = name;
		pickedPath = path;
		const kept = Object.fromEntries(page.url.searchParams);
		const query = new URLSearchParams({ ...kept, p: name, c: path });
		// eslint-disable-next-line svelte/no-navigation-without-resolve -- the path is resolve()d; the rule does not follow the query added to it
		void goto(resolve('/design/components') + `?${query}`, {
			replaceState: true,
			keepFocus: true,
			noScroll: true
		});
	}

	onMount(() => {
		void load();
		const timer = setInterval(() => void load(), POLL_MS);
		// flai serve watches the manifest but not a sub-project's sources: a new component shows at the next ask
		const stop = follow(['document'], () => void load());
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
			<ComponentList {projects} {picked} onpick={pick} />
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
			{#if component && server?.running && server.url}
				<div class="mt-4 flex gap-4 overflow-x-auto" data-testid="components-frames">
					{#each themes as theme, i (i)}
						<div class={[!width && 'min-w-0 flex-1']}>
							<ComponentFrame server={server.url} component={component.path} {theme} {width} />
						</div>
					{/each}
				</div>
			{:else if project && !component}
				<p class="mt-4 text-sm text-muted">Pick a component to frame it.</p>
			{/if}
		</section>
	</div>
{/if}
