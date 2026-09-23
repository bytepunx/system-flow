<script lang="ts">
	import { api } from '$lib/api';
	import { onMount } from 'svelte';
	import { projectState, type Project } from '$lib/project.svelte';

	type Manifest = { name: string; description?: string; template?: { version?: string } };
	type Item = { id: string; type: string; status: string; title: string; archived: boolean };
	type ProjectGlance = Project & {
		review?: number;
		threadsAwaiting?: number;
		agentAttending?: boolean;
	};

	let manifest = $state<Manifest | null>(null);
	let items = $state<Item[]>([]);
	let error = $state<string | null>(null);
	let glances = $state<ProjectGlance[]>([]);
	let filter = $state('');
	const filteredGlances = $derived(
		filter.trim()
			? glances.filter((p) =>
					`${p.name} ${p.key}`.toLowerCase().includes(filter.trim().toLowerCase())
				)
			: glances
	);

	const statuses = ['backlog', 'ready', 'in-progress', 'review', 'done'];
	const active = $derived(items.filter((i) => !i.archived));
	const count = (type: string, status: string) =>
		active.filter((i) => i.type === type && i.status === status).length;

	// The project list: shown while more than one project is known and none is chosen yet, in place
	// of the summary below, which needs one project named to ask anything of (S-0080).
	const showList = $derived(projectState.needsChoice && !projectState.current);

	async function loadOverview() {
		try {
			const [m, i] = await Promise.all([api('/api/manifest'), api('/api/items')]);
			if (!m.ok) throw new Error((await m.json()).error ?? m.statusText);
			manifest = await m.json();
			items = await i.json();
		} catch (e) {
			error = e instanceof Error ? e.message : String(e);
		}
	}
	async function loadGlances() {
		const r = await fetch('/api/projects');
		if (r.ok) glances = ((await r.json()).projects ?? []) as ProjectGlance[];
	}
	onMount(() => {
		void projectState.refresh(); // sets projectState.ready, which decides what this page shows
		void loadGlances();
	});
	// Whether to show the list or the single-project summary is not known until the project list has
	// been asked for at least once (S-0080): deciding at mount, before that answer arrives, would
	// call the single-project routes even when there turn out to be several.
	let decided = false;
	$effect(() => {
		if (decided || !projectState.ready) return;
		decided = true;
		if (!showList) void loadOverview();
	});
</script>

{#if showList}
	<h1 class="text-2xl font-semibold">Projects</h1>
	<p class="mt-1 text-ink-soft">flai on the host serves more than one project; choose one.</p>
	{#if glances.length > 1}
		<input
			type="search"
			placeholder="Filter projects…"
			class="mt-4 w-full max-w-sm rounded border border-line-strong bg-surface px-2 py-1 text-sm"
			data-testid="project-filter"
			bind:value={filter}
		/>
	{/if}
	<ul class="mt-4 space-y-2" data-testid="project-list">
		{#each filteredGlances as p (p.key)}
			<li class="rounded border border-line bg-surface p-3">
				<button
					type="button"
					class="text-left font-medium text-accent hover:underline"
					onclick={() => projectState.pick(p.key)}
				>
					{p.name}
				</button>
				<span class="ml-2 text-xs text-muted">{p.connected ? 'connected' : 'not connected'}</span>
				{#if p.connected}
					<p class="mt-1 text-xs text-muted">
						{p.review ?? '?'} in review · {p.threadsAwaiting ?? '?'} thread{p.threadsAwaiting === 1
							? ''
							: 's'} awaiting you · agent {p.agentAttending ? 'attending' : 'not attending'}
					</p>
				{/if}
			</li>
		{/each}
		{#if glances.length === 0}
			<li class="text-sm text-muted">No project has connected here yet.</li>
		{:else if filteredGlances.length === 0}
			<li class="text-sm text-muted">No project matches "{filter}".</li>
		{/if}
	</ul>
{:else if error}
	<p class="rounded border border-danger bg-danger-soft p-3 text-sm text-danger">
		Cannot read the project: {error}
	</p>
{:else if manifest}
	<h1 class="text-2xl font-semibold">{manifest.name}</h1>
	{#if manifest.description}<p class="mt-1 text-ink-soft">
			{manifest.description}
		</p>{/if}
	<p class="mt-1 text-xs text-muted">
		template {manifest.template?.version ?? '?'} · {active.length} active items · {items.length -
			active.length} archived
	</p>

	<div class="mt-6 overflow-x-auto">
		<table class="min-w-full text-sm">
			<thead>
				<tr class="text-left text-muted">
					<th class="py-2 pr-4"></th>
					{#each statuses as s (s)}<th class="py-2 pr-4 font-medium">{s}</th>{/each}
				</tr>
			</thead>
			<tbody>
				{#each ['epic', 'story', 'task'] as t (t)}
					<tr class="border-t border-line">
						<td class="py-2 pr-4 font-medium">{t}s</td>
						{#each statuses as s (s)}<td class="py-2 pr-4 tabular-nums">{count(t, s)}</td>{/each}
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{:else if !showList}
	<p class="text-sm text-muted">Loading…</p>
{/if}
