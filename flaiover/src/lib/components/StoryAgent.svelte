<script lang="ts">
	// What the agent flai started for this story is doing (S-0104): the board card's dot, with the
	// harness, the model, when it started, and why it waits or failed. It asks flai again when the
	// project's files change and, while the agent runs, now and then, since an agent ends without
	// changing a file. Nothing shows for a story no agent was started for.
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { projectState } from '$lib/project.svelte';
	import { activityLine, anyRunning, type HostAgent } from '$lib/activity';
	import AgentDot from './AgentDot.svelte';

	let { story }: { story: string } = $props();
	let status = $state<HostAgent | null>(null);

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
		{/if}
	</section>
{/if}
