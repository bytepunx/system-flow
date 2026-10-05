<script lang="ts">
	// The planner's page (S-0259): whether the plan host action is on and how to turn it on, the run
	// under way with its stream, the planner's activity document with its totals and its log newest
	// first, each past run with what its log entries say it cost and how it went, and a form that asks
	// flai serve to plan an epic or a story. What flai answered the form is said in its status line,
	// as PlanAction words it, and the page is asked to load again.
	import { resolve } from '$app/paths';
	import { api } from '$lib/api';
	import { dollars, duration } from '$lib/usage';
	import type { PlanRun } from '$lib/activity';
	import { currentRun, newestFirst, outcomeWord, pastRuns, type PlannerView } from '$lib/planner';
	import AgentStream from './AgentStream.svelte';

	let { view, onplanned }: { view: PlannerView; onplanned?: () => void } = $props();

	const current = $derived(currentRun(view.runs));
	const entries = $derived(newestFirst(view.activity.entries));
	const past = $derived(pastRuns(view.runs, view.activity.entries));

	let id = $state('');
	let said = $state<string | null>(null);
	let planning = $state(false);

	const ITEM_ID = /^[ES]-\d+$/;
	const at = (s: string) => s.replace('T', ' ').replace(/:\d\dZ$/, ' UTC');
	const who = (r: PlanRun) => [r.harness, r.model].filter(Boolean).join(', ') || r.command;
	const cost = (c: number, estimated: boolean) => `${dollars(c)}${estimated ? ' (estimated)' : ''}`;

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

{#snippet itemLink(item: string)}
	<a class="underline" href={resolve('/items/[id]', { id: item })}>{item}</a>
{/snippet}

{#snippet enableCommand()}
	On the host, in the project, run
	<code class="rounded bg-ground px-1 text-ink">flai serve enable plan</code>.
{/snippet}

<div class="space-y-4 text-sm">
	<section class="rounded border border-line bg-surface p-4" data-testid="planner-action">
		<h2 class="font-semibold text-ink">The plan host action</h2>
		{#if view.plan_enabled}
			<p class="mt-1" data-testid="planner-action-on">
				On: flai serve starts the planner when you ask.
			</p>
		{:else}
			<p class="mt-1 text-warn" data-testid="planner-action-off">
				Off: flai serve starts no planner for this project. {@render enableCommand()}
			</p>
		{/if}
	</section>

	<section class="rounded border border-line bg-surface p-4" data-testid="planner-current">
		<h2 class="font-semibold text-ink">Running now</h2>
		{#if current}
			<dl class="mt-1 grid grid-cols-[max-content_1fr] gap-x-4 gap-y-0.5">
				<dt class="text-muted">Item</dt>
				<dd data-testid="planner-current-item">{@render itemLink(current.item)}</dd>
				<dt class="text-muted">Agent</dt>
				<dd>{current.agent} ({who(current)})</dd>
				<dt class="text-muted">Started</dt>
				<dd>{at(current.started)}</dd>
				{#if current.trigger}
					<dt class="text-muted">Trigger</dt>
					<dd data-testid="planner-current-trigger">{current.trigger}</dd>
				{/if}
			</dl>
			<AgentStream plan story={current.item} started={current.started} />
		{:else}
			<p class="mt-1 text-muted" data-testid="planner-current-none">No planner is running.</p>
		{/if}
	</section>

	<section class="rounded border border-line bg-surface p-4" data-testid="planner-activity">
		<h2 class="font-semibold text-ink">
			Activity <span class="font-normal text-muted">in <code>{view.activity.path}</code></span>
		</h2>
		<dl
			class="mt-1 grid grid-cols-[max-content_1fr] gap-x-4 gap-y-0.5"
			data-testid="planner-totals"
		>
			<dt class="text-muted">Accrued cost</dt>
			<dd data-testid="planner-accrued-cost">{dollars(view.activity.accrued_cost)}</dd>
			<dt class="text-muted">Accrued time</dt>
			<dd data-testid="planner-accrued-time">{duration(view.activity.accrued_seconds)}</dd>
			<dt class="text-muted">Activities</dt>
			<dd data-testid="planner-activities">{view.activity.tasks_completed}</dd>
			<dt class="text-muted">Last run</dt>
			<dd data-testid="planner-last-run">
				{view.activity.last_run ? at(view.activity.last_run) : 'never'}
			</dd>
		</dl>
		{#if entries.length}
			<ol class="mt-3 space-y-2" data-testid="planner-entries">
				{#each entries as e (e.at + e.summary)}
					<li class="rounded border border-line p-2" data-testid="planner-entry">
						<p>
							<span class="font-mono text-xs text-muted">{at(e.at)}</span>
							<span data-testid="planner-entry-summary">{e.summary}</span>
						</p>
						<p class="text-xs text-muted" data-testid="planner-entry-details">
							{#if e.trigger}<span data-testid="planner-entry-trigger">trigger {e.trigger}</span>;
								items{:else}items{/if}
							<span data-testid="planner-entry-items"
								>{#each e.items as item, i (item)}{i ? ', ' : ''}{@render itemLink(
										item
									)}{:else}none{/each}</span
							>; took <span data-testid="planner-entry-duration">{duration(e.seconds)}</span>; cost
							<span data-testid="planner-entry-cost">{cost(e.cost, e.estimated)}</span>
						</p>
					</li>
				{/each}
			</ol>
		{:else}
			<p class="mt-3 text-muted" data-testid="planner-entries-none">
				The planner has not logged an activity yet.
			</p>
		{/if}
	</section>

	<section class="rounded border border-line bg-surface p-4" data-testid="planner-runs">
		<h2 class="font-semibold text-ink">Runs</h2>
		{#if past.length}
			<table class="mt-2 w-full text-left" data-testid="planner-runs-table">
				<thead class="text-muted">
					<tr>
						<th class="py-1 pr-4 font-normal">Item</th>
						<th class="py-1 pr-4 font-normal">Started</th>
						<th class="py-1 pr-4 font-normal">Ended</th>
						<th class="py-1 pr-4 font-normal">Outcome</th>
						<th class="py-1 pr-4 font-normal">Cost</th>
						<th class="py-1 font-normal">Why</th>
					</tr>
				</thead>
				<tbody>
					{#each past as { run, cost: c } (run.item + run.started)}
						<tr class="border-t border-line" data-testid="planner-run">
							<td class="py-1 pr-4">{@render itemLink(run.item)}</td>
							<td class="py-1 pr-4 whitespace-nowrap">{at(run.started)}</td>
							<td class="py-1 pr-4 whitespace-nowrap">{run.ended ? at(run.ended) : '—'}</td>
							<td class="py-1 pr-4" data-testid="planner-run-outcome">{outcomeWord(run)}</td>
							<td class="py-1 pr-4" data-testid="planner-run-cost"
								>{c.count ? cost(c.cost, c.estimated) : '—'}</td
							>
							<td class="py-1 text-danger" data-testid="planner-run-why">{run.why || run.error}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		{:else}
			<p class="mt-1 text-muted" data-testid="planner-runs-none">
				flai serve has started no planner for this project.
			</p>
		{/if}
	</section>

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
				The plan host action is off for this project, so flai serve cannot plan here. {@render enableCommand()}
			</p>
		{/if}
		{#if said}
			<p class="mt-2" role="status" data-testid="planner-said">{said}</p>
		{/if}
	</section>
</div>
