<script module lang="ts">
	import type { AgentRun } from '$lib/activity';
	import type { RunSpan } from '$lib/strategic';

	/** What the panel reads of a strategic agent's run, whichever agent's it is. */
	export type StrategicRun = RunSpan &
		Pick<AgentRun, 'agent' | 'command' | 'harness' | 'model' | 'why'> & { trigger?: string };

	/** The words and the command that turn a host action on, for a line that says it is off. */
	export { enableCommand };
</script>

<script lang="ts" generics="R extends StrategicRun">
	// What the Planner, Orchestrator, and Analyzer pages share (S-0259, S-0228): whether the agent's
	// host action is on, with the command that turns it on, the run under way with its stream, the
	// agent's activity document with its totals and its log newest first, and each past run with what
	// its log entries say it cost and how it went. Each page adds what is its own through snippets:
	// the state and the actions of its host action, the fields of its runs, what an entry shows besides
	// its items, duration, and cost, and sections before the log and after the runs.
	import type { Snippet } from 'svelte';
	import { resolve } from '$app/paths';
	import { localTime } from '$lib/localtime';
	import { dollars, duration } from '$lib/usage';
	import {
		currentRun,
		newestFirst,
		outcomeWord,
		pastRuns,
		type ActivityDocument,
		type ActivityEntry,
		type StrategicKind,
		type StreamRole
	} from '$lib/strategic';
	import AgentStream from './AgentStream.svelte';

	let {
		kind,
		action,
		enabled,
		on,
		activity,
		runs,
		stream,
		actionState,
		actions,
		currentFields,
		runHead,
		runCells,
		entryExtra,
		beforeActivity,
		children
	}: {
		/** The agent, in the page's words and as the prefix of its test IDs. */
		kind: StrategicKind;
		/** The host action that has flai serve start it: plan, orchestrate, or analyze. */
		action: string;
		/** Whether that host action is on for the project. */
		enabled: boolean;
		/** What the action's state says after "On:" while it is on. */
		on: string;
		activity: ActivityDocument;
		/** Every run flai knows of, in any order. */
		runs: R[];
		/** Where the run under way's stream is read: AgentStream's story, and its plan or role. */
		stream: (run: R) => { story: string; plan?: boolean; role?: StreamRole };
		/** Said in place of the "On:" line while the action is on, such as the orchestrator held. */
		actionState?: Snippet;
		/** What the operator can have flai do, under the action's state. */
		actions?: Snippet;
		/** dt/dd pairs of the run under way, before its agent and when it started. */
		currentFields?: Snippet<[R]>;
		/** th cells of the runs table, before when each run started. */
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

	const current = $derived(currentRun(runs));
	const entries = $derived(newestFirst(activity.entries));
	const past = $derived(pastRuns(runs, activity.entries));

	const who = (r: R) => [r.harness, r.model].filter(Boolean).join(', ') || r.command;
	const cost = (c: number, estimated: boolean) => `${dollars(c)}${estimated ? ' (estimated)' : ''}`;
</script>

{#snippet enableCommand(name: string)}
	On the host, in the project, run
	<code class="rounded bg-ground px-1 text-ink">flai serve enable {name}</code>.
{/snippet}

{#snippet itemLink(item: string)}
	<a class="underline" href={resolve('/items/[id]', { id: item })}>{item}</a>
{/snippet}

<div class="space-y-4 text-sm">
	<section class="rounded border border-line bg-surface p-4" data-testid="{kind}-action">
		<h2 class="font-semibold text-ink">The {action} host action</h2>
		{#if !enabled}
			<p class="mt-1 text-warn" data-testid="{kind}-action-off">
				Off: flai serve starts no {kind} for this project. {@render enableCommand(action)}
			</p>
		{:else if actionState}
			{@render actionState()}
		{:else}
			<p class="mt-1" data-testid="{kind}-action-on">On: {on}</p>
		{/if}
		{@render actions?.()}
	</section>

	<section class="rounded border border-line bg-surface p-4" data-testid="{kind}-current">
		<h2 class="font-semibold text-ink">Running now</h2>
		{#if current}
			{@const read = stream(current)}
			<dl class="mt-1 grid grid-cols-[max-content_1fr] gap-x-4 gap-y-0.5">
				{@render currentFields?.(current)}
				<dt class="text-muted">Agent</dt>
				<dd>{current.agent} ({who(current)})</dd>
				<dt class="text-muted">Started</dt>
				<dd>{localTime(current.started)}</dd>
				{#if current.trigger}
					<dt class="text-muted">Trigger</dt>
					<dd data-testid="{kind}-current-trigger">{current.trigger}</dd>
				{/if}
			</dl>
			<AgentStream story={read.story} plan={read.plan} role={read.role} started={current.started} />
		{:else}
			<p class="mt-1 text-muted" data-testid="{kind}-current-none">No {kind} is running.</p>
		{/if}
	</section>

	{@render beforeActivity?.()}

	<section class="rounded border border-line bg-surface p-4" data-testid="{kind}-activity">
		<h2 class="font-semibold text-ink">
			Activity <span class="font-normal text-muted">in <code>{activity.path}</code></span>
		</h2>
		<dl class="mt-1 grid grid-cols-[max-content_1fr] gap-x-4 gap-y-0.5" data-testid="{kind}-totals">
			<dt class="text-muted">Accrued cost</dt>
			<dd data-testid="{kind}-accrued-cost">{dollars(activity.accrued_cost)}</dd>
			<dt class="text-muted">Accrued time</dt>
			<dd data-testid="{kind}-accrued-time">{duration(activity.accrued_seconds)}</dd>
			<dt class="text-muted">Activities</dt>
			<dd data-testid="{kind}-activities">{activity.tasks_completed}</dd>
			<dt class="text-muted">Last run</dt>
			<dd data-testid="{kind}-last-run">
				{activity.last_run ? localTime(activity.last_run) : 'never'}
			</dd>
		</dl>
		{#if entries.length}
			<ol class="mt-3 space-y-2" data-testid="{kind}-entries">
				{#each entries as e (e.at + e.summary)}
					<li class="rounded border border-line p-2" data-testid="{kind}-entry">
						<p>
							<span class="font-mono text-xs text-muted">{localTime(e.at)}</span>
							<span data-testid="{kind}-entry-summary">{e.summary}</span>
						</p>
						<p class="text-xs text-muted" data-testid="{kind}-entry-details">
							{#if e.trigger}<span data-testid="{kind}-entry-trigger">trigger {e.trigger}</span>;
								items{:else}items{/if}
							<span data-testid="{kind}-entry-items"
								>{#each e.items as item, i (item)}{i ? ', ' : ''}{@render itemLink(
										item
									)}{:else}none{/each}</span
							>; took <span data-testid="{kind}-entry-duration">{duration(e.seconds)}</span>; cost
							<span data-testid="{kind}-entry-cost">{cost(e.cost, e.estimated)}</span>
						</p>
						{@render entryExtra?.(e)}
					</li>
				{/each}
			</ol>
		{:else}
			<p class="mt-3 text-muted" data-testid="{kind}-entries-none">
				The {kind} has not logged an activity yet.
			</p>
		{/if}
	</section>

	<section class="rounded border border-line bg-surface p-4" data-testid="{kind}-runs">
		<h2 class="font-semibold text-ink">Runs</h2>
		{#if past.length}
			<table class="mt-2 w-full text-left" data-testid="{kind}-runs-table">
				<thead class="text-muted">
					<tr>
						{@render runHead?.()}
						<th class="py-1 pr-4 font-normal">Started</th>
						<th class="py-1 pr-4 font-normal">Ended</th>
						<th class="py-1 pr-4 font-normal">Outcome</th>
						<th class="py-1 pr-4 font-normal">Cost</th>
						<th class="py-1 font-normal">Why</th>
					</tr>
				</thead>
				<tbody>
					<!-- two runs may start in the same second, so a run is keyed by its place -->
					{#each past as { run, cost: c }, i (i)}
						<tr class="border-t border-line" data-testid="{kind}-run">
							{@render runCells?.(run)}
							<td class="py-1 pr-4 whitespace-nowrap">{localTime(run.started)}</td>
							<td class="py-1 pr-4 whitespace-nowrap">{run.ended ? localTime(run.ended) : '—'}</td>
							<td class="py-1 pr-4" data-testid="{kind}-run-outcome">{outcomeWord(run)}</td>
							<td class="py-1 pr-4" data-testid="{kind}-run-cost"
								>{c.count ? cost(c.cost, c.estimated) : '—'}</td
							>
							<td class="py-1 text-danger" data-testid="{kind}-run-why">{run.why || run.error}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		{:else}
			<p class="mt-1 text-muted" data-testid="{kind}-runs-none">
				flai serve has started no {kind} for this project.
			</p>
		{/if}
	</section>

	{@render children?.()}
</div>
