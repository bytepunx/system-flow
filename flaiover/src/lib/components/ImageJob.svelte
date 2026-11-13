<script module lang="ts">
	/** The element id of an image job's card, for a link to it from another page. */
	export const imageJobAnchor = (id: string) => `image-job-${id}`;
</script>

<script lang="ts">
	// One image job of the project (S-0757): the images it makes and the way it left each, its state
	// and outcome, its current step, when it was asked for, started, and ended, and how long it ran,
	// with its output in a LogFeed that reads it through the page's readLog, open at first for a job
	// that runs or failed with output: followed about every second while the job runs and read once
	// more when it ends (ADR-0263 §4). A failed job shows its cause and its failed step in the error
	// style, and Retry, which calls the page's onretry and shows its refusal word for word. The
	// card's element id is imageJobAnchor(job.id), image-job-<id>, so a story's agent can link to it.
	import { untrack } from 'svelte';
	import LogFeed, { type LogPart } from '$lib/components/LogFeed.svelte';
	import type { ImageJob, ImageLog } from '$lib/images';
	import { localTime } from '$lib/localtime';

	let {
		job,
		readLog,
		onretry,
		disabled = false
	}: {
		/** The job, as images.get answers it. */
		job: ImageJob;
		/** Reads the job's log from an offset, as images.log answers it; throws a refusal's words. */
		readLog: (id: string, offset: number) => Promise<ImageLog>;
		/** Starts a new job for the same images; throws a refusal's words. Without it, no Retry. */
		onretry?: () => Promise<void>;
		/** Whether the agent host action is off: Retry is disabled. */
		disabled?: boolean;
	} = $props();

	const id = $derived(job.id);
	const live = $derived(job.state === 'running');
	const failed = $derived(job.outcome === 'failed');
	// the output opens for a job that runs or failed with output, and stays as the reader leaves it
	const opens = $derived(live || (failed && job.log_size > 0));

	let open = $state(false);
	let retrying = $state(false);
	let retryError = $state<string | null>(null);
	let openedFor = '';

	$effect.pre(() => {
		const now = opens;
		untrack(() => {
			if (id !== openedFor) {
				openedFor = id;
				open = now;
				retryError = null;
			} else if (now) {
				open = true;
			}
		});
	});

	// the job's log as LogFeed reads a step's, running while the job runs
	async function read(offset: number): Promise<LogPart> {
		const got = await readLog(id, offset);
		if (!got) throw new Error('The log read answered nothing.');
		return { ...got, outcome: got.state === 'running' ? 'running' : got.outcome };
	}

	async function retry() {
		if (!onretry) return;
		retrying = true;
		retryError = null;
		try {
			await onretry();
		} catch (e) {
			retryError = e instanceof Error ? e.message : String(e);
		} finally {
			retrying = false;
		}
	}

	const shown = $derived(job.state === 'ended' ? (job.outcome ?? 'ended') : job.state);
	// the job's state in LogFeed's words: a queued job's log is a step's that has not started
	const feedOutcome = $derived(
		job.state === 'queued' ? 'pending' : job.state === 'running' ? 'running' : (job.outcome ?? '')
	);
	const badge = $derived(
		shown === 'failed'
			? 'bg-danger-soft text-danger'
			: shown === 'passed'
				? 'bg-good-soft text-good-strong'
				: shown === 'running'
					? 'bg-info-soft text-info'
					: 'bg-surface text-muted'
	);
	const ways = $derived(new Map(job.ways.map((w) => [w.image, w.way])));

	function took(ms: number): string {
		const s = Math.round(ms / 1000);
		const m = Math.floor(s / 60);
		return m ? `${m}m ${String(s % 60).padStart(2, '0')}s` : `${s}s`;
	}
	const duration = $derived(
		job.duration_ms
			? took(job.duration_ms)
			: live && job.started
				? `since ${localTime(job.started)}`
				: ''
	);
</script>

<article
	id={imageJobAnchor(job.id)}
	class="rounded border border-line p-3"
	data-testid="image-job"
	data-job={job.id}
	data-state={job.state}
	data-outcome={job.outcome ?? ''}
>
	<div class="flex flex-wrap items-center gap-2 text-sm">
		<span class="font-mono font-medium">{job.id}</span>
		<span class="rounded px-2 py-0.5 text-xs {badge}" data-testid="image-job-status">{shown}</span>
		{#if duration}<span class="text-xs text-muted" data-testid="image-job-duration">{duration}</span
			>{/if}
	</div>
	<ul class="mt-2 space-y-1 text-sm" data-testid="image-job-images">
		{#each job.images as image (image)}
			<li data-testid="image-job-image">
				<code class="font-mono">{image}</code>
				{#if ways.get(image)}<span class="ml-1 text-muted" data-testid="image-job-way"
						>{ways.get(image)}</span
					>{/if}
			</li>
		{/each}
	</ul>
	<dl class="mt-2 grid grid-cols-[max-content_1fr] gap-x-4 gap-y-1 text-sm">
		{#if job.step && !failed}
			<dt class="text-muted">Step</dt>
			<dd data-testid="image-job-step">{job.step}</dd>
		{/if}
		<dt class="text-muted">Asked for</dt>
		<dd data-testid="image-job-requested">{localTime(job.requested)}</dd>
		{#if job.started}
			<dt class="text-muted">Started</dt>
			<dd data-testid="image-job-started">{localTime(job.started)}</dd>
		{/if}
		{#if job.ended}
			<dt class="text-muted">Ended</dt>
			<dd data-testid="image-job-ended">{localTime(job.ended)}</dd>
		{/if}
	</dl>
	{#if failed || (job.outcome === 'stopped' && job.cause)}
		<div
			class="mt-3 rounded border border-danger bg-danger-soft p-2 text-sm text-danger"
			role="alert"
			data-testid="image-job-failure"
		>
			{#if job.cause}<p data-testid="image-job-cause">{job.cause}</p>{/if}
			{#if failed && job.step}
				<p class="mt-1">
					At its step <code class="font-mono" data-testid="image-job-failed-step">{job.step}</code>
				</p>
			{/if}
		</div>
	{/if}
	{#if failed && onretry}
		<div class="mt-2 flex flex-wrap items-center gap-2">
			<button
				type="button"
				class="rounded bg-primary px-3 py-1 text-sm text-on-primary disabled:opacity-50"
				disabled={disabled || retrying}
				onclick={retry}
				data-testid="image-job-retry">Retry</button
			>
			{#if retryError}
				<p class="text-xs text-danger" role="alert" data-testid="image-job-retry-error">
					{retryError}
				</p>
			{/if}
		</div>
	{/if}
	<div class="mt-2">
		<LogFeed
			step={`the image job ${job.id}`}
			label="Output"
			outcome={feedOutcome}
			{read}
			bind:open
		/>
	</div>
</article>
