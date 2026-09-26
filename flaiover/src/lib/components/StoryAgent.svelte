<script lang="ts">
	// What the agent flai started for this story is doing (S-0104): the board card's dot, with the
	// harness, the model, when it started, and why it waits or failed. It asks flai again when the
	// project's files change and, while the agent runs, now and then, since an agent ends without
	// changing a file. Nothing shows for a story no agent was started for. For a story in ready or
	// in progress whose agent dropped or failed, Retry, at the top right, has flai start a new one
	// (S-0116), and hides once pressed until flai refuses or that one fails too (S-0118). For a story
	// in ready that has had no agent, Start agent has flai start it now, and the panel says why flai
	// serve has not (S-0115).
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { projectState } from '$lib/project.svelte';
	import { activityLine, anyRunning, type HostAgent } from '$lib/activity';
	import AgentDot from './AgentDot.svelte';

	let {
		story,
		status: storyStatus = '',
		writable = false
	}: { story: string; status?: string; writable?: boolean } = $props();
	let status = $state<HostAgent | null>(null);
	let acting = $state<'start' | 'restart' | null>(null);
	let actError = $state<{ action: 'start' | 'restart'; message: string } | null>(null);
	// The failed run Retry was pressed for: the button stays hidden while flai has it.
	let retried = $state<string | null>(null);

	async function ask() {
		try {
			const r = await api('/api/host-agent');
			if (r.ok) status = await r.json();
		} catch {
			// keep what we had
		}
	}
	onMount(() => {
		void ask();
		const es = new EventSource(projectState.tag('/api/events'));
		es.addEventListener('change', () => void ask());
		return () => es.close();
	});
	$effect(() => {
		if (!anyRunning(status)) return;
		const t = setInterval(() => void ask(), 15000);
		return () => clearInterval(t);
	});

	const activity = $derived(status?.enabled ? status.state?.stories?.[story] : undefined);
	const canRetry = $derived(
		writable &&
			activity?.state === 'failed' &&
			activity.run.started !== retried &&
			(storyStatus === 'ready' || storyStatus === 'in-progress')
	);
	// A ready story no agent was started for; flai says why not when it refuses.
	const canStart = $derived(writable && !!status?.enabled && storyStatus === 'ready' && !activity);
	// Why flai serve has not started it, from the reasons it gives for each story that waits.
	const waitingWhy = $derived(
		status?.state?.waiting?.split('; ').find((w) => new RegExp(`\\b${story}\\b`).test(w))
	);

	async function act(action: 'start' | 'restart') {
		if (acting) return;
		acting = action;
		actError = null;
		if (action === 'restart') retried = activity?.run.started ?? null;
		try {
			const r = await api(`/api/items/${story}/agent`, {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ action })
			});
			if (!r.ok) {
				const body = (await r.json().catch(() => ({}))) as { error?: string };
				actError = { action, message: body.error ?? `the agent could not be ${action}ed` };
			}
		} catch (e) {
			actError = { action, message: e instanceof Error ? e.message : String(e) };
		}
		if (actError) retried = null;
		acting = null;
		await ask();
	}
	const at = (s: string) => s.replace('T', ' ').replace(/:\d\dZ$/, ' UTC');
</script>

{#if activity}
	<section class="rounded border border-line bg-surface p-3" data-testid="story-agent">
		<div class="mb-2 flex items-center gap-2">
			<h2 class="flex items-center gap-2 font-medium"><AgentDot {activity} /> Agent</h2>
			{#if canRetry}
				<button
					class="ml-auto rounded border border-line px-2 py-1 text-xs disabled:opacity-60"
					onclick={() => act('restart')}
					disabled={acting !== null}
					data-testid="story-agent-retry">Retry</button
				>
			{/if}
		</div>
		<p class="text-xs" data-testid="story-agent-line">{activityLine(activity)}</p>
		<p class="mt-1 text-xs text-muted">
			{activity.run.agent}, started {at(activity.run.started)}{activity.run.ended
				? `, ended ${at(activity.run.ended)}`
				: ''}{activity.run.exit !== undefined && activity.run.exit !== null
				? ` (exit ${activity.run.exit})`
				: ''}
		</p>
		{#if activity.state === 'waiting' && activity.thread}
			<p class="mt-1 text-xs">It asked in {activity.thread}: answer it below and it goes on.</p>
		{/if}
		{#if activity.state === 'failed'}
			<p class="mt-1 text-xs text-muted" data-testid="story-agent-failed">
				{#if activity.run.log}Its output is in <code class="break-all">{activity.run.log}</code> on the
					host.{/if}
				Moving the story back to ready starts another, and so does changing its agent while it is in ready.
			</p>
			{#if actError?.action === 'restart'}
				<p class="mt-1 text-xs text-warn" data-testid="story-agent-retry-error">
					{actError.message}
				</p>
			{/if}
		{/if}
	</section>
{:else if canStart || actError?.action === 'start'}
	<section class="rounded border border-line bg-surface p-3" data-testid="story-agent">
		<h2 class="mb-2 font-medium">Agent</h2>
		<p class="text-xs" data-testid="story-agent-line">
			flai serve has started no agent for this story{waitingWhy ? `: ${waitingWhy}` : '.'}
		</p>
		{#if canStart}
			<button
				class="mt-2 rounded border border-line px-2 py-1 text-xs disabled:opacity-60"
				onclick={() => act('start')}
				disabled={acting !== null}
				data-testid="story-agent-start">{acting === 'start' ? 'Starting…' : 'Start agent'}</button
			>
		{/if}
		{#if actError?.action === 'start'}
			<p class="mt-1 text-xs text-warn" data-testid="story-agent-start-error">{actError.message}</p>
		{/if}
	</section>
{/if}
