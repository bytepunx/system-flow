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
	// the log and after the runs as the panel's snippets.
	import type { Snippet } from 'svelte';
	import { resolve } from '$app/paths';
	import { api } from '$lib/api';
	import { roleControls, roleRuns, roleStream, triggerLine, type RoleActionName } from '$lib/roles';
	import type { ActivityDocument, ActivityEntry } from '$lib/strategic';
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
		const log = a.log ? `; its output is in ${a.log} on the host` : '';
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
		<p class="mt-2" role="status" data-testid="{role.name}-said">{said}</p>
	{/if}
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
	{beforeActivity}
	{children}
/>
