<script lang="ts">
	// /ops/log?path=<path>: one log on the host in a LogViewer (S-0785), the page a LogLink in a
	// message opens, so the operator reads the log here rather than in a shell on the host. It is in
	// the rail's Ops group with no page of its own current. With no path, it says how to open a log.
	import { page } from '$app/state';
	import LogViewer from '$lib/components/LogViewer.svelte';

	const path = $derived(page.url.searchParams.get('path') ?? '');
</script>

<svelte:head><title>{path ? `${path.split('/').pop()} · ` : ''}Log · flaiover</title></svelte:head>

<h1 class="text-2xl font-semibold">Log</h1>
{#if path}
	<div class="mt-4">
		<LogViewer {path} />
	</div>
{:else}
	<p class="mt-2 text-sm text-muted" data-testid="log-page-none">
		No log given. Open a log from a link to its path in a message, a status, or an error.
	</p>
{/if}
