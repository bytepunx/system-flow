<script lang="ts">
	// /activity: who is working on what (S-0042), and what each agent flai started is saying and
	// doing (S-0142). What flai knows of agents is asked again when a work item or a thread changes, and while one runs
	// now and then, since an agent ends without changing a file.
	import { api } from '$lib/api';
	import { onMount } from 'svelte';
	import { debounced, follow, listen } from '$lib/events';
	import ActivityView from '$lib/components/ActivityView.svelte';
	import { anyRunning, type HostAgent } from '$lib/activity';

	let streams = $state<never[] | null>(null);
	let error = $state<string | null>(null);
	let host = $state<HostAgent | null>(null);

	async function loadAgents() {
		try {
			const r = await api('/api/host-agent');
			if (r.ok) host = await r.json();
		} catch {
			// keep what we had: the streams are the page
		}
	}

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
		void loadAgents();
		// the streams are read from the narratives and the work items, the agents from the items and
		// the threads; flai serve says when an agent starts or ends (S-0161)
		const agents = debounced(() => void loadAgents());
		const stops = [
			follow(['item', 'narrative'], () => void load()),
			follow(['item', 'thread'], () => void loadAgents()),
			listen({ agent: () => agents() })
		];
		return () => {
			for (const stop of stops) stop();
			agents.stop();
		};
	});
	$effect(() => {
		if (!anyRunning(host)) return;
		const t = setInterval(() => void loadAgents(), 15000);
		return () => clearInterval(t);
	});
</script>

<svelte:head><title>Activity · flaiover</title></svelte:head>

<h1 class="mb-1 text-2xl font-semibold">Activity</h1>
<p class="mb-4 text-sm text-muted">
	Read from the narratives in <code>wip/agents</code>. An agent that stops writing simply grows old
	here. Each agent flai serve started shows its stream, read from the log flai gave it, and one that
	runs can be stopped.
</p>
{#if error}
	<p class="rounded border border-danger bg-danger-soft p-3 text-sm text-danger" role="alert">
		{error}
	</p>
{:else if streams === null}
	<p class="text-sm text-muted">Loading…</p>
{:else}
	<ActivityView
		{streams}
		agents={host?.state?.stories ?? {}}
		canStop={!!host?.enabled}
		onstopped={() => void loadAgents()}
	/>
{/if}
