<script lang="ts">
	// Plan on an epic's or a story's page (S-0208): flai starts the planner for the item now, through
	// the plan host action. It shows only while the operator has that action on for the project, asked
	// again for each item and as the project changes, and says when the item's newest planner run has
	// not ended. What flai answered, the run or why it refused, goes to the page's notice, as the
	// page's other actions say theirs.
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { follow, listen } from '$lib/events';
	import { localTime } from '$lib/localtime';
	import type { PlanRun } from '$lib/activity';

	let { id, onresult }: { id: string; onresult: (text: string) => void } = $props();
	let enabled = $state(false);
	let run = $state<PlanRun | null>(null);
	let starting = $state(false);

	async function ask(of: string) {
		try {
			const r = await api(`/api/items/${of}/plan`);
			if (!r.ok || of !== id) return;
			const body = (await r.json()) as { plan_enabled?: boolean; run?: PlanRun | null };
			enabled = body.plan_enabled === true;
			run = body.run ?? null;
		} catch {
			// keep what we had
		}
	}
	$effect(() => {
		run = null;
		void ask(id);
	});
	// The planner writes the items it plans, and flai serve says when its run starts or ends.
	onMount(() => {
		const unfollow = follow(['item'], () => void ask(id));
		const unlisten = listen({ agent: (of) => of === id && void ask(id) });
		return () => {
			unfollow();
			unlisten();
		};
	});

	const going = $derived(!!run && !run.ended && !run.error);

	async function plan() {
		if (starting) return;
		starting = true;
		try {
			const r = await api(`/api/items/${id}/plan`, { method: 'POST' });
			const data = (await r.json().catch(() => ({}))) as Partial<PlanRun> & { error?: string };
			if (!r.ok) onresult(`refused: ${data.error ?? r.statusText}`);
			else
				onresult(
					`planner started for ${id} (pid ${data.pid})${data.log ? `; its output is in ${data.log} on the host` : ''}`
				);
		} catch (e) {
			onresult(`refused: ${e instanceof Error ? e.message : String(e)}`);
		}
		starting = false;
		await ask(id);
	}
</script>

{#if enabled}
	<button
		type="button"
		class="rounded border border-line-strong px-2 py-1 text-xs hover:bg-raised disabled:opacity-60"
		onclick={plan}
		disabled={starting}
		data-testid="item-plan">{starting ? 'Planning…' : 'Plan'}</button
	>
	{#if going && run}
		<span class="self-center text-xs text-muted" data-testid="item-plan-running"
			>planner running since {localTime(run.started)}</span
		>
	{/if}
{/if}
