<script lang="ts">
	// The header's word on the host flai (ADR-0029). Nothing is shown when the container was given no
	// agent credential: an older flai started it, and the banner below the header explains. Links to
	// /host, where restart, upgrade, and stop live (S-0081).
	import { resolve } from '$app/paths';
	import { hostFlai } from '$lib/hostflai.svelte';

	const status = $derived(hostFlai.status);
</script>

{#if status?.configured}
	{#if hostFlai.usable}
		<a
			href={resolve('/host')}
			class="rounded border border-line px-2 py-1 text-xs text-muted no-underline hover:text-accent"
			data-host-flai="connected"
			title={`flai ${status.flai ?? ''} on the host, connected since ${status.since ?? ''}`}
			>host flai: connected</a
		>
	{:else}
		<a
			href={resolve('/host')}
			class="rounded border border-warn bg-warn-soft px-2 py-1 text-xs no-underline"
			data-host-flai="gone"
			title={status.error ??
				'Start it on the host with: flai dashboard (in the project), or flai serve start'}
			>host flai: not connected</a
		>
	{/if}
{/if}
