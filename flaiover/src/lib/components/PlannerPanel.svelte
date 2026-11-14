<script lang="ts">
	// The planner's page (S-0259), the role page for the planner (S-0420): whether the plan host
	// action is on and how to turn it on, the run under way with its item and its stream, the
	// planner's activity document with its totals and its log newest first, each past run with its
	// item and what its log entries say it cost and how it went, and a form after them that asks flai
	// serve to plan an epic or a story. What flai answered the form is said in its status line, as
	// PlanAction words it, and the page is asked to load again. The form is its own and stands for
	// the role page's Run, which it does not show; the rest is the role page's.
	import { api } from '$lib/api';
	import type { PlanRun } from '$lib/activity';
	import type { PlannerView } from '$lib/planner';
	import type { Snippet } from 'svelte';
	import LogText from './LogText.svelte';
	import RolePage, { roleWithRuns } from './RolePage.svelte';
	import { enableCommand } from './StrategicAgentPanel.svelte';

	let {
		view,
		onplanned,
		settings
	}: {
		view: PlannerView;
		onplanned?: () => void;
		/** The planner's settings, shown under the plan host action (S-0533). */
		settings?: Snippet;
	} = $props();

	/** The planner as the role page reads a role, from its file and from what the view says. */
	const role = $derived(
		roleWithRuns(
			{
				name: 'planner',
				role: 'plan',
				path: 'design/workflows/roles/planner.md',
				item: ['epic', 'story'],
				trigger: ['asked', 'event', 'schedule'],
				action: { name: 'plan', enabled: view.plan_enabled },
				held: false
			},
			view.runs
		)
	);

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
				? `planner started for ${of} (pid ${data.pid})${data.log ? `; its output: ${data.log}` : ''}`
				: `refused: ${data.error ?? r.statusText}`;
		} catch (err) {
			said = `refused: ${err instanceof Error ? err.message : String(err)}`;
		}
		planning = false;
		onplanned?.();
	}
</script>

<RolePage
	{role}
	activity={view.activity}
	on="flai serve starts the planner when you ask."
	actions={null}
	{settings}
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
			<p class="mt-2" role="status" data-testid="planner-said"><LogText text={said} /></p>
		{/if}
	</section>
</RolePage>
