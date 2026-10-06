<script lang="ts">
	// /workflow/analyzer (S-0229): the analyzer's settings, read from settings.get on arrival, after a
	// save, and when system-flow.yaml changes. Its status, activity, and runs come with S-0228.
	import { api } from '$lib/api';
	import { onMount } from 'svelte';
	import { follow } from '$lib/events';
	import StrategicSettings from '$lib/components/StrategicSettings.svelte';
	import type { SettingsView, StrategicSettings as Strategic } from '$lib/settings';

	let strategic = $state<Strategic | null>(null);
	let read = $state(false);
	let error = $state<string | null>(null);

	async function loadSettings() {
		try {
			const r = await api('/api/settings');
			const body = await r.json().catch(() => ({}));
			if (!r.ok) throw new Error(body.error ?? r.statusText);
			strategic = (body as SettingsView).host?.strategic ?? null;
			read = true;
			error = null;
		} catch (e) {
			error = e instanceof Error ? e.message : String(e);
		}
	}
	onMount(() => {
		void loadSettings();
		return follow(['project'], (changes) => {
			if (changes.some((c) => c.kind === 'project')) void loadSettings();
		});
	});
</script>

<svelte:head><title>Analyzer · flaiover</title></svelte:head>

<h1 class="mb-1 text-2xl font-semibold">Analyzer</h1>
<p class="mb-4 text-sm text-muted">Its status, activity log, and runs will be shown here.</p>
{#if error}
	<p class="mb-4 rounded border border-danger bg-danger-soft p-3 text-sm text-danger" role="alert">
		{error}
	</p>
{/if}
{#if strategic}
	<StrategicSettings block="analysis" {strategic} onsaved={() => void loadSettings()} />
{:else if read}
	<p class="text-sm text-muted" data-testid="strategic-absent">
		flai on the host gave no settings for the analyzer: the project may not be open, or flai may be
		older than this dashboard.
	</p>
{:else if !error}
	<p class="text-sm text-muted">Loading…</p>
{/if}
