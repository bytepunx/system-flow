<script lang="ts">
	// /workflow/orchestrator (S-0228): what the orchestrator is doing and has decided, read from its
	// activity document and from what flai serve knows of its runs. It is asked again when the
	// orchestrator's document changes, and when flai serve says an orchestrator run started or ended,
	// or that the operator held it or lifted the hold, which change no file. Its settings (S-0229) are
	// read apart, on arrival, after a save, and when system-flow.yaml changes, so that a reload of the
	// orchestrator does not refill the form while the operator types.
	import { api } from '$lib/api';
	import { onMount } from 'svelte';
	import { debounced, follow, listen } from '$lib/events';
	import OrchestratorPanel from '$lib/components/OrchestratorPanel.svelte';
	import StrategicSettings from '$lib/components/StrategicSettings.svelte';
	import type { OrchestratorView } from '$lib/strategic';
	import type { SettingsView, StrategicSettings as Strategic } from '$lib/settings';

	/** What flai serve's agent event names for the orchestrator (ORCHESTRATOR in $lib/server/agent). */
	const ORCHESTRATOR = 'orchestrator';

	let view = $state<OrchestratorView | null>(null);
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
			const r = await api('/api/orchestrator');
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
			// of the narratives, only the orchestrator's own document is this page's
			follow(['narrative'], (changes) => {
				if (!view || changes.some((c) => c.kind !== 'narrative' || c.path === view!.activity.path))
					void load();
			}),
			// a story's agent or a planner says nothing of the orchestrator
			listen({
				agent: (id) => {
					if (id === ORCHESTRATOR) again();
				}
			})
		];
		return () => {
			for (const stop of stops) stop();
			again.stop();
		};
	});
</script>

<svelte:head><title>Orchestrator · flaiover</title></svelte:head>

<h1 class="mb-1 text-2xl font-semibold">Orchestrator</h1>
<p class="mb-4 text-sm text-muted">
	Read from <code>wip/agents/orchestrator.md</code> and from what flai serve knows of the orchestrator's
	runs.
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
	<OrchestratorPanel {view} onchanged={() => void load()} />
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
		<StrategicSettings block="orchestration" {strategic} onsaved={() => void loadSettings()} />
	{:else if settingsRead}
		<p class="text-sm text-muted" data-testid="strategic-absent">
			flai on the host gave no settings for the orchestrator: the project may not be open, or flai
			may be older than this dashboard.
		</p>
	{/if}
</div>
