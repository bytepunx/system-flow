<script lang="ts">
	// /workflow/planner (S-0259): what the planner is doing and has done, read from its activity
	// document and from what flai serve knows of its runs. It is asked again when the planner's
	// document or a work item changes, when flai serve says a run started or ended, and every 15
	// seconds while a run is under way, since a run's stream and its end change no file.
	import { api } from '$lib/api';
	import { onMount } from 'svelte';
	import { debounced, follow, listen } from '$lib/events';
	import PlannerPanel from '$lib/components/PlannerPanel.svelte';
	import { currentRun, type PlannerView } from '$lib/planner';

	let view = $state<PlannerView | null>(null);
	let error = $state<string | null>(null);

	async function load() {
		try {
			const r = await api('/api/planner');
			if (!r.ok) throw new Error((await r.json().catch(() => ({}))).error ?? r.statusText);
			view = await r.json();
			error = null;
		} catch (e) {
			error = e instanceof Error ? e.message : String(e);
		}
	}
	onMount(() => {
		void load();
		const again = debounced(() => void load());
		const stops = [
			// of the narratives, only the planner's own document is this page's
			follow(['narrative'], (changes) => {
				if (!view || changes.some((c) => c.kind !== 'narrative' || c.path === view!.activity.path))
					void load();
			}),
			// the planner writes the items it plans
			follow(['item'], () => void load()),
			listen({ agent: () => again() })
		];
		return () => {
			for (const stop of stops) stop();
			again.stop();
		};
	});
	$effect(() => {
		if (!view || !currentRun(view.runs)) return;
		const t = setInterval(() => void load(), 15000);
		return () => clearInterval(t);
	});
</script>

<svelte:head><title>Planner · flaiover</title></svelte:head>

<h1 class="mb-1 text-2xl font-semibold">Planner</h1>
<p class="mb-4 text-sm text-muted">
	Read from <code>wip/agents/planner.md</code> and from what flai serve knows of the planner's runs.
</p>
<!-- a reload that fails says so above what was last shown rather than in its place (S-0178) -->
{#if error}
	<p class="mb-4 rounded border border-danger bg-danger-soft p-3 text-sm text-danger" role="alert">
		{error}
	</p>
{/if}
{#if view === null}
	{#if !error}<p class="text-sm text-muted">Loading…</p>{/if}
{:else}
	<PlannerPanel {view} onplanned={() => void load()} />
{/if}
