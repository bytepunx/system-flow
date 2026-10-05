<script lang="ts">
	// /activity: who is working on what (S-0042), and what each agent flai started is saying and
	// doing (S-0142), the project's orchestrator among them (S-0218). What flai knows of agents is asked
	// again when a work item or a thread changes, when flai serve says an agent started or ended, and
	// while one runs now and then, since an agent ends without changing a file.
	import { api } from '$lib/api';
	import { onMount } from 'svelte';
	import { debounced, follow, listen } from '$lib/events';
	import ActivityView from '$lib/components/ActivityView.svelte';
	import AgentStream from '$lib/components/AgentStream.svelte';
	import { anyRunning, orchestratorLine, orchestratorRunning, type HostAgent } from '$lib/activity';

	let streams = $state<never[] | null>(null);
	let error = $state<string | null>(null);
	let host = $state<HostAgent | null>(null);
	const orchestrator = $derived(host?.state?.orchestrator ?? null);

	async function loadAgents() {
		try {
			const r = await api('/api/host-agent');
			if (!r.ok) return;
			const next: HostAgent = await r.json();
			// with no flai connected the answer has no state: keep the agents last known, so their
			// windows stay put rather than vanish and come back from the log's tail when flai does
			// (S-0178); whether one can be stopped follows what flai says now
			host = next.state || !host?.state ? next : { ...next, state: host.state };
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
		// the threads; flai serve says when an agent starts or ends (S-0161), the orchestrator's too
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
		if (!anyRunning(host) && !orchestratorRunning(orchestrator)) return;
		const t = setInterval(() => void loadAgents(), 15000);
		return () => clearInterval(t);
	});
</script>

<svelte:head><title>Activity · flaiover</title></svelte:head>

<h1 class="mb-1 text-2xl font-semibold">Activity</h1>
<p class="mb-4 text-sm text-muted">
	Read from the narratives in <code>wip/agents</code>. An agent that stops writing simply grows old
	here. Each agent flai serve started shows its stream, read from the log flai gave it, and one that
	runs can be stopped. The project's orchestrator, which flai serve runs while the orchestrate host
	action is on, shows its newest run and its stream above the stories.
</p>
<!-- a reload that fails says so above the streams last shown rather than in their place, so the
	page keeps its length and the operator their place in it (S-0178) -->
{#if error}
	<p class="mb-4 rounded border border-danger bg-danger-soft p-3 text-sm text-danger" role="alert">
		{error}
	</p>
{/if}
{#if streams === null}
	{#if !error}<p class="text-sm text-muted">Loading…</p>{/if}
{:else}
	{#if orchestrator}
		<section
			class="mb-3 rounded border border-line bg-surface p-3 text-sm"
			data-testid="orchestrator-run"
		>
			<div class="flex flex-wrap items-baseline gap-x-3 gap-y-1">
				<span class="font-medium">Orchestrator</span>
				<span class="min-w-0 flex-1 text-xs text-muted" data-testid="orchestrator-line"
					>{orchestratorLine(orchestrator)}</span
				>
			</div>
			<dl class="mt-2 grid grid-cols-[max-content_1fr] gap-x-4 gap-y-1 text-xs">
				<dt class="text-muted">agent</dt>
				<dd>
					{orchestrator.agent}{#if orchestrator.session}<span class="text-muted"
							>&nbsp;· session {orchestrator.session}</span
						>{/if}
				</dd>
				<dt class="text-muted">started</dt>
				<dd>{orchestrator.started}</dd>
				{#if orchestrator.ended}
					<dt class="text-muted">ended</dt>
					<dd>{orchestrator.ended}</dd>
				{/if}
			</dl>
			<!-- a run that could not be started has no log to read -->
			{#if !orchestrator.error}
				<AgentStream
					orchestrator
					story="orchestrator"
					started={orchestrator.started}
					open={orchestratorRunning(orchestrator)}
				/>
			{/if}
		</section>
	{/if}
	<ActivityView
		{streams}
		agents={host?.state?.stories ?? {}}
		canStop={!!host?.enabled}
		onstopped={() => void loadAgents()}
	/>
{/if}
