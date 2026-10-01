<script lang="ts">
	// /threads: every open thread of the project, shown as a story's page shows its own (S-0173,
	// TH-0041). An inbox thread on no item leads here, opened on that thread, to be answered.
	import Threads from '$lib/components/Threads.svelte';
	import { api } from '$lib/api';
	import { page } from '$app/state';
	import { onMount } from 'svelte';

	let writable = $state(false);
	onMount(() => {
		api('/api/board')
			.then(async (r) => (writable = r.ok ? (await r.json()).writable : false))
			.catch(() => (writable = false));
	});
</script>

<svelte:head><title>Threads · flaiover</title></svelte:head>

<h1 class="text-2xl font-semibold">Threads</h1>
<p class="mt-1 text-sm text-muted">
	Every open thread of the project, on whatever it is anchored to. Threads are started on the item
	or document they are about.
</p>
<Threads {writable} select={page.url.searchParams.get('thread') ?? undefined} />
