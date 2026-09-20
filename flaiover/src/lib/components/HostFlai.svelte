<script lang="ts">
	// Whether a flai on the host has this dashboard connected (ADR-0029). Nothing is shown when the
	// container was given no agent credential: an older flai started it, and there is nothing to say.
	import { onMount } from 'svelte';
	import { api } from '$lib/api';

	type Status = {
		configured: boolean;
		connected: boolean;
		since?: string;
		flai?: string;
		error?: string;
	};

	let { everyMs = 10000 }: { everyMs?: number } = $props();
	let status = $state<Status | null>(null);

	async function load() {
		try {
			const r = await api('/api/agent');
			if (r.ok) status = await r.json();
		} catch {
			// The page says nothing new when the dashboard itself cannot be reached.
		}
	}
	onMount(() => {
		void load();
		const t = setInterval(load, everyMs);
		return () => clearInterval(t);
	});
</script>

{#if status?.configured}
	{#if status.connected && !status.error}
		<span
			class="rounded border border-line px-2 py-1 text-xs text-muted"
			data-host-flai="connected"
			title={`flai ${status.flai ?? ''} on the host, connected since ${status.since ?? ''}`}
			>host flai: connected</span
		>
	{:else}
		<span
			class="rounded border border-warn bg-warn-soft px-2 py-1 text-xs"
			data-host-flai="gone"
			role="status"
			title={status.error ??
				'Start it on the host with: flai dashboard (in the project), or flai serve start'}
			>host flai: not connected</span
		>
	{/if}
{/if}
