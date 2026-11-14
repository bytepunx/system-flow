<script module lang="ts">
	import type { RoleState } from '$lib/roles';
	import { currentRun, type RoleRun } from '$lib/strategic';

	/**
	 * A role's state from a page that has every run of it in one list, as the Planner,
	 * Orchestrator, and Analyzer pages' views do: those under way, the newest of them its run, and
	 * those that ended.
	 */
	export function roleWithRuns<R extends RoleRun>(
		role: Omit<RoleState<R>, 'run' | 'running' | 'runs'>,
		runs: R[]
	): RoleState<R> {
		const going = (r: R) => !r.ended && !r.error;
		return {
			...role,
			run: currentRun(runs),
			running: runs.filter(going),
			runs: runs.filter((r) => !going(r))
		};
	}
</script>

<script lang="ts" generics="R extends RoleRun">
	// A role's page (S-0420), for any role agent.status reports, the built-in ones and those the
	// project defines, on what the Planner, Orchestrator, and Analyzer pages share
	// (StrategicAgentPanel): its host action's state, held while the operator has stopped it; the run
	// under way with its stream; its activity document, wip/agents/<name>.md; and its runs with their
	// cost and outcome, each with its item for a role whose run takes one. Its controls are those
	// roleControls offers: Run for a role its file lets the operator ask for, with an item when its
	// run takes one, and Stop and Start for the orchestrator, the one role with a hold. Each goes
	// through /api/roles/<name>, says what flai answered, and asks the page to load again; each is
	// disabled, saying why, while the host action is off. A page with its own controls gives them as
	// actions in their place, or none there when it shows them elsewhere, and its settings, the
	// fields of its run, its runs' columns, what an entry shows besides its own, and sections before
	// the log and after the runs as the panel's snippets. Before the log, what the role spent and the
	// hours it worked each day of the charts' window, read from flai stats (ADR-0249), drawn in the
	// colour the Strategic Cost and Strategic Use charts give it, and read again when its activity
	// document's totals change.
	import { untrack, type Snippet } from 'svelte';
	import { resolve } from '$app/paths';
	import { api } from '$lib/api';
	import { roleControls, roleRuns, roleStream, triggerLine, type RoleActionName } from '$lib/roles';
	import type { ActivityDocument, ActivityEntry } from '$lib/strategic';
	import { dollars } from '$lib/usage';
	import { chartWindow } from '$lib/chartwindow.svelte';
	import { themeState } from '$lib/theme.svelte';
	import { theme } from '$lib/viz/palette';
	import { hours, normalise, roleByDay, type Report, type RoleMeasure } from '$lib/viz/charts';
	import Chart from './Chart.svelte';
	import LogText from './LogText.svelte';
	import StrategicAgentPanel, { enableCommand } from './StrategicAgentPanel.svelte';

	let {
		role,
		activity,
		onchanged,
		on,
		actions,
		settings,
		currentFields,
		runHead,
		runCells,
		entryExtra,
		beforeActivity,
		children
	}: {
		/** The role as agent.status reports it under roles. */
		role: RoleState<R>;
		/** Its activity document. */
		activity: ActivityDocument;
		/** Called once flai has answered a Run, a Stop, or a Start. */
		onchanged?: () => void;
		/** What the action's state says after "On:", in place of what its triggers say. */
		on?: string;
		/** The page's own controls, in place of Run, Stop, and Start; null for none there. */
		actions?: Snippet | null;
		/** The role's settings, between the host action and the run under way. */
		settings?: Snippet;
		/** dt/dd pairs of the run under way, after its item and before its agent. */
		currentFields?: Snippet<[R]>;
		/** th cells of the runs table, after its item and before when each run started. */
		runHead?: Snippet;
		/** td cells of a run, matching runHead. */
		runCells?: Snippet<[R]>;
		/** What an entry of the log shows after its items, duration, and cost. */
		entryExtra?: Snippet<[ActivityEntry]>;
		/** Sections between the run under way and the activity log. */
		beforeActivity?: Snippet;
		/** Sections after the runs. */
		children?: Snippet;
	} = $props();

	const offers = $derived(roleControls(role));
	const runs = $derived(roleRuns(role));
	const takesItem = $derived(role.item.length > 0);
	/** A role that runs for the project runs one at a time; one per item runs one on each. */
	const busy = $derived(!takesItem && role.run !== null);
	const enabled = $derived(role.action.enabled);
	const placeholder = $derived(
		role.item.map((k) => (k === 'epic' ? 'E-0001' : k === 'story' ? 'S-0001' : k)).join(' or ')
	);

	let item = $state('');
	let asking = $state<RoleActionName | null>(null);
	let said = $state<string | null>(null);

	/** The charts' window's flai stats report, null until it is read, and why it could not be. */
	let report = $state<Report | null>(null);
	let statsError = $state<string | null>(null);
	const since = $derived(chartWindow.since);
	const t = $derived(theme(themeState.dark));
	/** The window's days: flai sends one for every day of it, so none is a flai older than S-0205. */
	const days = $derived(report?.strategic_days ?? []);
	/** The role's own days in the window, those it has an entry that ended on. */
	const own = $derived(days.flatMap((d) => (d.agents[role.name] ? [d.agents[role.name]] : [])));
	const ownCost = $derived(Math.round(own.reduce((n, u) => n + u.cost, 0) * 1e4) / 1e4);
	const ownSeconds = $derived(own.reduce((n, u) => n + u.seconds, 0));
	/** What the chart draws: what the role spent each day, or the hours it worked. */
	let measure = $state<RoleMeasure>('cost');
	const measures: { measure: RoleMeasure; label: string }[] = [
		{ measure: 'cost', label: 'Cost' },
		{ measure: 'use', label: 'Use' }
	];

	let asked = 0;
	async function loadStats(window: string) {
		const mine = ++asked;
		try {
			const r = await api(`/api/stats?since=${window}&type=story&bucket=day`);
			const body = await r.json().catch(() => ({}));
			if (mine !== asked) return;
			if (!r.ok) throw new Error(body.error ?? r.statusText);
			report = normalise(body as Report);
			statsError = null;
		} catch (e) {
			if (mine === asked) statsError = e instanceof Error ? e.message : String(e);
		}
	}
	// asked again for another role or window, and when the role's totals change: a run that ended
	$effect(() => {
		void [role.name, activity.accrued_cost, activity.accrued_seconds];
		const window = since;
		untrack(() => void loadStats(window));
	});

	type Answer = {
		item?: string;
		agent?: string;
		command?: string;
		pid?: number;
		log?: string;
		orchestrator?: { agent?: string; command?: string; pid?: number };
		error?: string;
	};

	/** What flai did, in the words its journal uses (hostapi describeRole, describeOrchestrate). */
	function done(action: RoleActionName, a: Answer): string {
		const run = a.orchestrator ?? a;
		const agent = run.agent ?? role.name;
		const pid = run.pid ? ` (pid ${run.pid})` : '';
		if (action === 'stop') return `stopped the ${role.name} ${agent}${pid} and held it stopped`;
		if (action === 'start')
			return `lifted the hold and started ${run.command ?? `the ${role.name}`} as ${agent}${pid}`;
		const onItem = a.item ? ` on ${a.item}` : '';
		const log = a.log ? `; its output: ${a.log}` : '';
		return `started ${run.command ?? `the ${role.name}`} to run the ${role.name}${onItem} as ${agent}${pid}${log}`;
	}

	async function ask(action: RoleActionName) {
		if (asking) return;
		asking = action;
		try {
			const r = await api(`/api/roles/${encodeURIComponent(role.name)}`, {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify(
					action === 'run' && takesItem ? { action, item: item.trim() } : { action }
				)
			});
			const data = (await r.json().catch(() => ({}))) as Answer;
			said = r.ok ? done(action, data) : `refused: ${data.error ?? r.statusText}`;
		} catch (err) {
			said = `refused: ${err instanceof Error ? err.message : String(err)}`;
		}
		asking = null;
		onchanged?.();
	}
</script>

{#snippet held()}
	<p class="mt-1 text-warn" data-testid="{role.name}-action-held">
		Held: you stopped the {role.name}, and flai serve does not start it again until you Start it.
	</p>
{/snippet}

<!-- Run is run-now: the panel's runs table names each row {kind}-run -->
{#snippet button(action: RoleActionName, label: string, disabled: boolean)}
	<button
		type="button"
		class="rounded border border-line-strong px-2 py-1 hover:bg-raised disabled:opacity-60"
		{disabled}
		onclick={() => void ask(action)}
		data-testid="{role.name}-{action === 'run' ? 'run-now' : action}"
		>{asking === action ? `${label}…` : label}</button
	>
{/snippet}

{#snippet controls()}
	{#if offers.run || offers.stop || offers.start}
		<div class="mt-2 flex flex-wrap items-center gap-2">
			{#if offers.item}
				<label class="text-muted" for="{role.name}-item">Item</label>
				<input
					id="{role.name}-item"
					class="w-40 rounded border border-line-strong bg-surface px-2 py-1 font-mono"
					{placeholder}
					bind:value={item}
					disabled={!enabled}
					data-testid="{role.name}-item"
				/>
			{/if}
			{#if offers.run}
				{@render button(
					'run',
					'Run',
					!enabled || busy || asking !== null || (offers.item && !item.trim())
				)}
			{/if}
			{#if offers.stop}{@render button('stop', 'Stop', !enabled || asking !== null)}{/if}
			{#if offers.start}{@render button('start', 'Start', !enabled || asking !== null)}{/if}
		</div>
		{#if !enabled}
			<p class="mt-2 text-muted" data-testid="{role.name}-controls-off">
				The {role.action.name} host action is off for this project, so flai serve cannot run the {role.name}
				here. {@render enableCommand(role.action.name)}
			</p>
		{:else if offers.run && busy}
			<p class="mt-2 text-muted" data-testid="{role.name}-run-busy">
				A {role.name} is running; one runs at a time, so Run waits for it to end.
			</p>
		{/if}
	{/if}
	{#if said}
		<p class="mt-2" role="status" data-testid="{role.name}-said"><LogText text={said} /></p>
	{/if}
{/snippet}

<!-- what the role spent and worked per day of the charts' window, before the log -->
{#snippet usage()}
	<section class="rounded border border-line bg-surface p-4" data-testid="{role.name}-usage">
		<h2 class="font-semibold text-ink">Cost and use per day</h2>
		{#if statsError}
			<p class="mt-2 text-danger" role="alert" data-testid="{role.name}-usage-error">
				flai stats could not be read: {statsError}
			</p>
		{:else if !report}
			<p class="mt-2 text-muted">Loading…</p>
		{:else if days.length === 0}
			<p class="mt-2 text-muted" data-testid="{role.name}-usage-older">
				The flai on the host sends no strategic use per day: it is older than the dashboard. Upgrade
				it with <code>flai self-upgrade</code>.
			</p>
		{:else if own.length === 0}
			<p class="mt-2 text-muted" data-testid="{role.name}-usage-none">
				No activity of the {role.name} ended in the last {since}.
			</p>
		{:else}
			<p class="mt-2 text-muted" data-testid="{role.name}-usage-summary">
				Over the last {since} the {role.name} spent {dollars(ownCost)} and worked {hours(
					ownSeconds
				)}h, on {own.length}
				{own.length === 1 ? 'day' : 'days'}.
			</p>
			<div class="mt-2 flex gap-1" role="group" aria-label="What the chart draws">
				{#each measures as m (m.measure)}
					<button
						type="button"
						class="rounded border px-3 py-1 text-sm {measure === m.measure
							? 'border-primary bg-primary-soft'
							: 'border-line'}"
						aria-pressed={measure === m.measure}
						onclick={() => (measure = m.measure)}
						data-testid="{role.name}-usage-{m.measure}">{m.label}</button
					>
				{/each}
			</div>
			<div class="mt-2" data-testid="{role.name}-usage-chart">
				<Chart option={roleByDay(report, t, role.name, measure)} theme={t} height={240} />
			</div>
		{/if}
		<p class="mt-2 text-sm text-muted">
			Beside the other roles on <a
				class="underline"
				href={resolve('/charts/[kind]', { kind: 'strategic-cost' })}>Strategic Cost</a
			>
			and
			<a class="underline" href={resolve('/charts/[kind]', { kind: 'strategic-use' })}
				>Strategic Use</a
			>.
		</p>
	</section>
	{@render beforeActivity?.()}
{/snippet}

{#snippet currentItem(r: R)}
	{#if takesItem}
		<dt class="text-muted">Item</dt>
		<dd data-testid="{role.name}-current-item">
			{#if r.item}<a class="underline" href={resolve('/items/[id]', { id: r.item })}>{r.item}</a
				>{:else}none{/if}
		</dd>
	{/if}
	{@render currentFields?.(r)}
{/snippet}

{#snippet itemHead()}
	{#if takesItem}<th class="py-1 pr-4 font-normal">Item</th>{/if}
	{@render runHead?.()}
{/snippet}

{#snippet itemCells(r: R)}
	{#if takesItem}
		<td class="py-1 pr-4" data-testid="{role.name}-run-item"
			>{#if r.item}<a class="underline" href={resolve('/items/[id]', { id: r.item })}>{r.item}</a
				>{:else}—{/if}</td
		>
	{/if}
	{@render runCells?.(r)}
{/snippet}

<StrategicAgentPanel
	kind={role.name}
	action={role.action.name}
	{enabled}
	on={on ?? triggerLine(role)}
	{activity}
	{runs}
	stream={(r) => roleStream(role, r)}
	actionState={role.held ? held : undefined}
	actions={actions === undefined ? controls : (actions ?? undefined)}
	{settings}
	currentFields={currentItem}
	runHead={itemHead}
	runCells={itemCells}
	{entryExtra}
	beforeActivity={usage}
	{children}
/>
