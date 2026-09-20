<script lang="ts">
	// /activity: who is working on what (S-0042).
	import { api } from '$lib/api';
	import { projectState } from '$lib/project.svelte';
	import { onMount } from 'svelte';
	import ActivityView from '$lib/components/ActivityView.svelte';

	let streams = $state<never[] | null>(null);
	let error = $state<string | null>(null);

	async function load() {
		try {
			const r = await api('/api/activity');
			if (!r.ok) throw new Error((await r.json().catch(() => ({}))).error ?? r.statusText);
			streams = (await r.json()).streams;
			error = null;
		} catch (e) {
			error = e instanceof Error ? e.message : String(e);
		}
	}
	onMount(() => {
		void load();
		const es = new EventSource(projectState.tag('/api/events'));
		let timer: ReturnType<typeof setTimeout> | undefined;
		es.addEventListener('change', () => {
			clearTimeout(timer);
			timer = setTimeout(load, 300);
		});
		return () => es.close();
	});
</script>

<svelte:head><title>Activity · flaiover</title></svelte:head>

<h1 class="mb-1 text-2xl font-semibold">Activity</h1>
<p class="mb-4 text-sm text-muted">
	Read from the narratives in <code>wip/agents</code>. An agent that stops writing simply grows old
	here.
</p>
{#if error}
	<p class="rounded border border-danger bg-danger-soft p-3 text-sm text-danger" role="alert">
		{error}
	</p>
{:else if streams === null}
	<p class="text-sm text-muted">Loading…</p>
{:else}
	<ActivityView {streams} />
{/if}
