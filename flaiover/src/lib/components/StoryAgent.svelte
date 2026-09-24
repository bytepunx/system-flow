<script lang="ts">
	// What the agent flai started for this story is doing (S-0104): the board card's dot, with the
	// harness, the model, when it started, and why it waits or failed. It asks flai again when the
	// project's files change and, while the agent runs, now and then, since an agent ends without
	// changing a file. Nothing shows for a story no agent was started for. For a story in ready or
	// in progress whose agent dropped or failed, Restart agent has flai start a new one (S-0116).
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
	let restarting = $state(false);
	let restartError = $state<string | null>(null);

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
	const canRestart = $derived(
		writable &&
			activity?.state === 'failed' &&
			(storyStatus === 'ready' || storyStatus === 'in-progress')
	);

	async function restart() {
		if (restarting) return;
		restarting = true;
		restartError = null;
		try {
			const r = await api(`/api/items/${story}/agent`, {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ action: 'restart' })
			});
			if (!r.ok) {
				const body = (await r.json().catch(() => ({}))) as { error?: string };
				restartError = body.error ?? 'the agent could not be restarted';
			}
		} catch (e) {
			restartError = e instanceof Error ? e.message : String(e);
		}
		restarting = false;
		await ask();
	}
	const at = (s: string) => s.replace('T', ' ').replace(/:\d\dZ$/, ' UTC');
</script>

{#if activity}
	<section class="rounded border border-line bg-surface p-3" data-testid="story-agent">
		<h2 class="mb-2 flex items-center gap-2 font-medium">
			<AgentDot {activity} /> Agent
		</h2>
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
			{#if canRestart}
				<button
					class="mt-2 rounded border border-line px-2 py-1 text-xs disabled:opacity-60"
					onclick={restart}
					disabled={restarting}
					data-testid="story-agent-restart">{restarting ? 'Restarting…' : 'Restart agent'}</button
				>
			{/if}
			{#if restartError}
				<p class="mt-1 text-xs text-warn" data-testid="story-agent-restart-error">{restartError}</p>
			{/if}
		{/if}
	</section>
{/if}
