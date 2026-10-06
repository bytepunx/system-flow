<script lang="ts">
	// /workflow/planner (S-0259): what the planner is doing and has done, read from its activity
	// document and from what flai serve knows of its runs. It is asked again when the planner's
	// document or a work item changes, when flai serve says a run started or ended, and every 15
	// seconds while a run is under way, since a run's stream and its end change no file. Its settings
	// (S-0229) are read apart, on arrival, after a save, and when system-flow.yaml changes, never on
	// the timer, so that a read does not refill the form while the operator types.
	import { api } from '$lib/api';
	import { onMount } from 'svelte';
	import { debounced, follow, listen } from '$lib/events';
	import PlannerPanel from '$lib/components/PlannerPanel.svelte';
	import StrategicSettings from '$lib/components/StrategicSettings.svelte';
	import { currentRun, type PlannerView } from '$lib/planner';
	import type { SettingsView, StrategicSettings as Strategic } from '$lib/settings';

	let view = $state<PlannerView | null>(null);
	let error = $state<string | null>(null);
	let strategic = $state<Strategic | null>(null);
	let settingsRead = $state(false);
	let settingsError = $state<string | null>(null);

	async function loadSettings() {
		try {
			const r = await api('/api/settings');
			const body = await r.json().catch(() => ({}));
			if (!r.ok) throw new Error(body.error ?? r.statusText);
			strategic = (body as SettingsView).host?.strategic ?? null;
			settingsRead = true;
			settingsError = null;
		} catch (e) {
			settingsError = e instanceof Error ? e.message : String(e);
		}
	}

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
		void loadSettings();
		const again = debounced(() => void load());
		const stops = [
			follow(['project'], (changes) => {
				if (changes.some((c) => c.kind === 'project')) void loadSettings();
			}),
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

<div class="mt-8 border-t border-line pt-4">
	{#if settingsError}
		<p
			class="mb-4 rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
			role="alert"
		>
			{settingsError}
		</p>
	{/if}
	{#if strategic}
		<StrategicSettings block="planning" {strategic} onsaved={() => void loadSettings()} />
	{:else if settingsRead}
		<p class="text-sm text-muted" data-testid="strategic-absent">
			flai on the host gave no settings for the planner: the project may not be open, or flai may be
			older than this dashboard.
		</p>
	{/if}
</div>
