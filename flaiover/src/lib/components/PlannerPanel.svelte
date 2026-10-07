<script lang="ts">
	// The planner's page (S-0259): whether the plan host action is on and how to turn it on, the run
	// under way with its stream, the planner's activity document with its totals and its log newest
	// first, each past run with what its log entries say it cost and how it went, and a form that asks
	// flai serve to plan an epic or a story. What flai answered the form is said in its status line,
	// as PlanAction words it, and the page is asked to load again. What it shares with the
	// Orchestrator and Analyzer pages is StrategicAgentPanel's (S-0228); the item a run plans, and
	// the form, are its own.
	import { resolve } from '$app/paths';
	import { api } from '$lib/api';
	import type { PlanRun } from '$lib/activity';
	import type { PlannerView } from '$lib/planner';
	import StrategicAgentPanel, { enableCommand } from './StrategicAgentPanel.svelte';

	let { view, onplanned }: { view: PlannerView; onplanned?: () => void } = $props();

	let id = $state('');
	let said = $state<string | null>(null);
	let planning = $state(false);

	const ITEM_ID = /^[ES]-\d+$/;

	async function plan(e: SubmitEvent) {
		e.preventDefault();
		if (planning) return;
		const of = id.trim().toUpperCase();
		if (!ITEM_ID.test(of)) {
			said = `${of || 'nothing'} is not an epic's or a story's ID; give one such as E-0001 or S-0001`;
			return;
		}
		planning = true;
		try {
			const r = await api(`/api/items/${of}/plan`, { method: 'POST' });
			const data = (await r.json().catch(() => ({}))) as Partial<PlanRun> & { error?: string };
			said = r.ok
				? `planner started for ${of} (pid ${data.pid})${data.log ? `; its output is in ${data.log} on the host` : ''}`
				: `refused: ${data.error ?? r.statusText}`;
		} catch (err) {
			said = `refused: ${err instanceof Error ? err.message : String(err)}`;
		}
		planning = false;
		onplanned?.();
	}
</script>

{#snippet item(run: PlanRun)}
	<dt class="text-muted">Item</dt>
	<dd data-testid="planner-current-item">{@render itemLink(run.item)}</dd>
{/snippet}

{#snippet itemHead()}
	<th class="py-1 pr-4 font-normal">Item</th>
{/snippet}

{#snippet itemCell(run: PlanRun)}
	<td class="py-1 pr-4">{@render itemLink(run.item)}</td>
{/snippet}

{#snippet itemLink(id: string)}
	<a class="underline" href={resolve('/items/[id]', { id })}>{id}</a>
{/snippet}

<StrategicAgentPanel
	kind="planner"
	action="plan"
	enabled={view.plan_enabled}
	on="flai serve starts the planner when you ask."
	activity={view.activity}
	runs={view.runs}
	stream={(run) => ({ story: run.item, plan: true })}
	currentFields={item}
	runHead={itemHead}
	runCells={itemCell}
>
	<section class="rounded border border-line bg-surface p-4" data-testid="planner-ask">
		<h2 class="font-semibold text-ink">Plan an epic or a story</h2>
		<form class="mt-2 flex flex-wrap items-center gap-2" onsubmit={plan} data-testid="planner-form">
			<label class="sr-only" for="planner-id">The epic's or the story's ID</label>
			<input
				id="planner-id"
				class="w-40 rounded border border-line-strong px-2 py-1 font-mono"
				placeholder="E-0001 or S-0001"
				bind:value={id}
				disabled={!view.plan_enabled}
				data-testid="planner-id"
			/>
			<button
				type="submit"
				class="rounded border border-line-strong px-2 py-1 hover:bg-raised disabled:opacity-60"
				disabled={!view.plan_enabled || planning}
				data-testid="planner-plan">{planning ? 'Planning…' : 'Plan'}</button
			>
		</form>
		{#if !view.plan_enabled}
			<p class="mt-2 text-muted" data-testid="planner-form-off">
				The plan host action is off for this project, so flai serve cannot plan here. {@render enableCommand(
					'plan'
				)}
			</p>
		{/if}
		{#if said}
			<p class="mt-2" role="status" data-testid="planner-said">{said}</p>
		{/if}
	</section>
</StrategicAgentPanel>
