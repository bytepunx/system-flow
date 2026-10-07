<script lang="ts">
	// The analyzer's page (S-0228), on what it shares with the Planner and Orchestrator pages
	// (StrategicAgentPanel): whether the analyze host action is on, or off with the command that turns
	// it on; the run under way with its focus and its stream; its activity log, each entry with the
	// report it wrote linked on the documents page; and its runs, each with its focus and report. Run
	// starts it, looking for the focus chosen or for all of them, through /api/analyzer, saying what
	// flai answered, and the page is asked to load again. Run is disabled, saying why, while the
	// action is off or a run goes.
	import { resolve } from '$app/paths';
	import { api } from '$lib/api';
	import {
		ANALYZER_FOCUSES,
		type ActivityEntry,
		type AnalyzerRun,
		type AnalyzerStarted,
		type AnalyzerView
	} from '$lib/strategic';
	import StrategicAgentPanel, { enableCommand } from './StrategicAgentPanel.svelte';

	let { view, onchanged }: { view: AnalyzerView; onchanged?: () => void } = $props();

	/** The focus asked for; '' asks for all of them. */
	let focus = $state('');
	let running = $state(false);
	let said = $state<string | null>(null);

	/**
	 * A report's path in a summary: flai names it there (ADR-0099), as `<design>/analysis/<UTC day>-<focus>.md`,
	 * whether the agent's last line named it or flai added `(report …)` to it.
	 */
	const REPORT = /(?:[\w.-]+\/)*analysis\/\d{4}-\d{2}-\d{2}-[\w-]+\.md/;

	/** The report an entry names: a run's report its summary holds, or a report's path in it. */
	function reportOf(e: ActivityEntry): string | null {
		const run = view.runs.find((r) => r.report && e.summary.includes(r.report));
		return run?.report ?? e.summary.match(REPORT)?.[0] ?? null;
	}

	/** A run's focus in the page's words: flai's `all` is what a run asked for no focus has. */
	const focusWord = (f: string) => (f === 'all' || !f ? 'all three' : f);

	async function run() {
		if (running) return;
		running = true;
		try {
			const r = await api('/api/analyzer', {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify(focus ? { focus } : {})
			});
			const data = (await r.json().catch(() => ({}))) as Partial<AnalyzerStarted> & {
				error?: string;
			};
			// as flai's journal words it (hostapi describeAnalyze)
			said = r.ok
				? `started ${data.command} to analyze, focus ${data.focus}, as ${data.agent} (pid ${data.pid})${data.log ? `; its output is in ${data.log} on the host` : ''}`
				: `refused: ${data.error ?? r.statusText}`;
		} catch (err) {
			said = `refused: ${err instanceof Error ? err.message : String(err)}`;
		}
		running = false;
		onchanged?.();
	}
</script>

{#snippet docLink(path: string)}
	<a class="underline" href={resolve('/docs/[...path]', { path })}>{path}</a>
{/snippet}

{#snippet controls()}
	<div class="mt-2 flex flex-wrap items-center gap-2">
		<label class="text-muted" for="analyzer-focus">Focus</label>
		<select
			id="analyzer-focus"
			class="rounded border border-line-strong bg-surface px-2 py-1"
			bind:value={focus}
			disabled={!view.enabled || !!view.run}
			data-testid="analyzer-focus"
		>
			<option value="">all three</option>
			{#each ANALYZER_FOCUSES as f (f)}
				<option value={f}>{f}</option>
			{/each}
		</select>
		<button
			type="button"
			class="rounded border border-line-strong px-2 py-1 hover:bg-raised disabled:opacity-60"
			disabled={!view.enabled || !!view.run || running}
			onclick={() => void run()}
			data-testid="analyzer-run">{running ? 'Running…' : 'Run'}</button
		>
	</div>
	{#if !view.enabled}
		<p class="mt-2 text-muted" data-testid="analyzer-run-off">
			The analyze host action is off for this project, so flai serve cannot run the analyzer here. {@render enableCommand(
				'analyze'
			)}
		</p>
	{:else if view.run}
		<p class="mt-2 text-muted" data-testid="analyzer-run-busy">
			An analyzer is running, focus {focusWord(view.run.focus)}; one runs at a time, so Run waits
			for it to end.
		</p>
	{/if}
	{#if said}
		<p class="mt-2" role="status" data-testid="analyzer-said">{said}</p>
	{/if}
{/snippet}

{#snippet currentFocus(r: AnalyzerRun)}
	<dt class="text-muted">Focus</dt>
	<dd data-testid="analyzer-current-focus">{focusWord(r.focus)}</dd>
{/snippet}

{#snippet runHead()}
	<th class="py-1 pr-4 font-normal">Focus</th>
	<th class="py-1 pr-4 font-normal">Report</th>
{/snippet}

{#snippet runCells(r: AnalyzerRun)}
	<td class="py-1 pr-4" data-testid="analyzer-run-focus">{focusWord(r.focus)}</td>
	<td class="py-1 pr-4" data-testid="analyzer-run-report"
		>{#if r.report}{@render docLink(r.report)}{:else}—{/if}</td
	>
{/snippet}

{#snippet entryReport(e: ActivityEntry)}
	{@const report = reportOf(e)}
	<p class="text-xs" data-testid="analyzer-entry-report">
		{#if report}report {@render docLink(report)}{:else}<span class="text-muted">no report</span
			>{/if}
	</p>
{/snippet}

<StrategicAgentPanel
	kind="analyzer"
	action="analyze"
	enabled={view.enabled}
	on="flai serve starts the analyzer when you Run it, and when analysis.schedule comes round."
	activity={view.activity}
	runs={view.runs}
	stream={() => ({ story: 'analyzer', role: 'analyze' })}
	actions={controls}
	currentFields={currentFocus}
	{runHead}
	{runCells}
	entryExtra={entryReport}
/>
