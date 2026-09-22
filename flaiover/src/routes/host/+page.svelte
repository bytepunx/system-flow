<script lang="ts">
	// /host: the dashboard managed from its own page (S-0081) — what it is running, and, once
	// the operator enables the dashboard host action, restart, upgrade, and stop.
	import HostPanel from '$lib/components/HostPanel.svelte';
</script>

<svelte:head><title>flaiover — host</title></svelte:head>

<h1 class="mb-4 text-lg font-semibold text-ink">Host</h1>
<HostPanel />
