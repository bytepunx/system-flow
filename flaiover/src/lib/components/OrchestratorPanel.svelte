<script lang="ts">
	// The orchestrator's page (S-0228), the role page for the orchestrator (S-0420): whether the
	// orchestrate host action is on, off with the command that turns it on, or held after the
	// operator's Stop; the run under way with its stream; its last decisions, each with its reason,
	// above its whole log; and its runs. Its Stop and Start stand for the role page's: Stop ends the
	// run and holds it stopped, and Start lifts the hold, each through /api/orchestrator, saying what
	// flai answered, and the page is asked to load again. Both are disabled, saying why, while the
	// action is off.
	import { api } from '$lib/api';
	import { localTime } from '$lib/localtime';
	import {
		newestFirst,
		type OrchestratorAction,
		type OrchestratorRun,
		type OrchestratorView
	} from '$lib/strategic';
	import type { Snippet } from 'svelte';
	import LogText from './LogText.svelte';
	import RolePage, { roleWithRuns } from './RolePage.svelte';
	import { enableCommand } from './StrategicAgentPanel.svelte';

	let {
		view,
		onchanged,
		settings
	}: {
		view: OrchestratorView;
		onchanged?: () => void;
		/** The orchestrator's settings, shown under Stop and Start (S-0533). */
		settings?: Snippet;
	} = $props();

	/** The orchestrator as the role page reads a role, from its file and from what the view says. */
	const role = $derived(
		roleWithRuns(
			{
				name: 'orchestrator',
				role: 'orchestrate',
				path: 'design/workflows/roles/orchestrator.md',
				item: [],
				trigger: ['asked', 'event', 'schedule'],
				action: { name: 'orchestrate', enabled: view.enabled },
				held: view.held
			},
			view.runs
		)
	);

	/** How many of the newest entries are shown as its last decisions. */
	const DECISIONS = 5;

	const decisions = $derived(newestFirst(view.activity.entries).slice(0, DECISIONS));

	let asking = $state<OrchestratorAction | null>(null);
	let said = $state<string | null>(null);

	/** What flai did, as its journal words it (hostapi describeOrchestrate). */
	function done(action: OrchestratorAction, run: Partial<OrchestratorRun> | undefined): string {
		const agent = run?.agent ?? 'orchestrator';
		if (action === 'start')
			return `lifted the hold and started ${run?.command ?? 'the orchestrator'} as ${agent} (pid ${run?.pid})`;
		return run?.pid
			? `stopped the orchestrator ${agent} (pid ${run.pid}) and held it stopped`
			: `held the orchestrator ${agent} stopped`;
	}

	async function ask(action: OrchestratorAction) {
		if (asking) return;
		asking = action;
		try {
			const r = await api('/api/orchestrator', {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ action })
			});
			const data = (await r.json().catch(() => ({}))) as {
				orchestrator?: Partial<OrchestratorRun>;
				error?: string;
			};
			said = r.ok ? done(action, data.orchestrator) : `refused: ${data.error ?? r.statusText}`;
		} catch (err) {
			said = `refused: ${err instanceof Error ? err.message : String(err)}`;
		}
		asking = null;
		onchanged?.();
	}
</script>

{#snippet button(action: OrchestratorAction, label: string, disabled: boolean)}
	<button
		type="button"
		class="rounded border border-line-strong px-2 py-1 hover:bg-raised disabled:opacity-60"
		{disabled}
		onclick={() => void ask(action)}
		data-testid="orchestrator-{action}">{asking === action ? `${label}…` : label}</button
	>
{/snippet}

{#snippet controls()}
	<div class="mt-2 flex flex-wrap items-center gap-2">
		{#if !view.enabled}
			{@render button('stop', 'Stop', true)}
			{@render button('start', 'Start', true)}
		{:else if view.held}
			{@render button('start', 'Start', asking !== null)}
		{:else if view.runs.length}
			{@render button('stop', 'Stop', asking !== null)}
		{/if}
	</div>
	{#if !view.enabled}
		<p class="mt-2 text-muted" data-testid="orchestrator-controls-off">
			The orchestrate host action is off for this project, so flai serve runs no orchestrator to
			stop or start. {@render enableCommand('orchestrate')}
		</p>
	{/if}
	{#if said}
		<p class="mt-2" role="status" data-testid="orchestrator-said"><LogText text={said} /></p>
	{/if}
{/snippet}

{#snippet lastDecisions()}
	<section class="rounded border border-line bg-surface p-4" data-testid="orchestrator-decisions">
		<h2 class="font-semibold text-ink">Last decisions</h2>
		{#if decisions.length}
			<ol class="mt-1 space-y-1">
				{#each decisions as d (d.at + d.summary)}
					<li data-testid="orchestrator-decision">
						<span class="font-mono text-xs text-muted">{localTime(d.at)}</span>
						<span data-testid="orchestrator-decision-summary">{d.summary}</span>
					</li>
				{/each}
			</ol>
		{:else}
			<p class="mt-1 text-muted" data-testid="orchestrator-decisions-none">
				The orchestrator has decided nothing yet.
			</p>
		{/if}
	</section>
{/snippet}

<RolePage
	{role}
	activity={view.activity}
	actions={controls}
	{settings}
	beforeActivity={lastDecisions}
/>
