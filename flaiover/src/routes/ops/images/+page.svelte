<script lang="ts">
	// /ops/images: the images the project's agents run in and flai serve's image jobs that make them
	// (S-0757, ADR-0295 §5), flai's images.get through /api/images, asked afresh on load and every 3 s
	// while a job runs or waits. Each image shows its kind and inner harness, its ref, where its name
	// came from, whether it is on the host, and its state with why it is stale; when flai could not
	// judge them, it says why. Each job, the running one, those that wait, and the last ended of each
	// set of images, is an ImageJob card that follows its output live through
	// /api/images?job=&offset=; a failed one has Retry, which posts /api/images/retry, flai's
	// images.ensure under the agent host action, and the page takes the images and jobs it answers.
	// With no job running, waiting, or failed and an image missing or stale, Make the images starts
	// one the same way, so the page always offers a way to make them. A refusal is shown word for
	// word; a 403 says the agent host action is off with the command that turns it on, and disables
	// Retry; a 501 says to upgrade flai on the host in place of the page. A link to a job's card,
	// /ops/images#image-job-<id> as a story's agent gives it, scrolls to it once the jobs are drawn.
	import type { AgentImage, ImageJob as ImageJobRecord, ImageLog, ImagesState } from '$lib/images';
	import { api } from '$lib/api';
	import { onMount, tick } from 'svelte';
	import ImageJob from '$lib/components/ImageJob.svelte';

	// How often the page asks again while a job runs or waits; a job takes minutes.
	const POLL_MS = 3000;

	let data = $state<ImagesState | null>(null);
	let error = $state<string | null>(null);
	// The host's flai answered 501: its message, which says what it cannot do.
	let tooOld = $state<string | null>(null);
	// The agent host action was refused (403) on a retry: the action and the command that enables it.
	let refused = $state<{ action: string; enable: string } | null>(null);
	// flai's refusal or error for Make the images, word for word.
	let makeError = $state<string | null>(null);
	let making = $state(false);
	// The page has scrolled to the job its hash names, or found none to scroll to.
	let landed = false;

	const runs = $derived(!!data?.running || (data?.queued.length ?? 0) > 0);
	const failed = $derived(!!data?.ended.some((j) => j.outcome === 'failed'));
	const needed = $derived(!!data?.images.some((i) => i.state === 'missing' || i.state === 'stale'));
	const nothing = $derived(!!data && !data.running && !data.queued.length && !data.ended.length);

	const errorText = (e: unknown) => (e instanceof Error ? e.message : String(e));
	const kindText = (i: AgentImage) => (i.kind === 'layer' ? 'project layer' : 'harness');
	const sourceText = (s: AgentImage['source']) =>
		s === 'layer' ? "the project's layer" : s === 'host' ? "the host's own" : "flai's default";
	const stateClass = (s: AgentImage['state']) =>
		s === 'current'
			? 'bg-good-soft text-good-strong'
			: s === 'stale' || s === 'missing'
				? 'bg-warn-soft text-warn'
				: 'bg-surface text-muted';

	/** Reads the images and jobs; a failed read keeps what was last shown. */
	async function load() {
		try {
			const r = await api('/api/images');
			const body = await r.json().catch(() => ({}));
			if (r.status === 501) {
				tooOld = body.error ?? r.statusText;
				return;
			}
			if (!r.ok) throw new Error(body.error ?? r.statusText);
			data = body as ImagesState;
			tooOld = error = null;
			await land();
		} catch (e) {
			error = `Could not read the agent images: ${errorText(e)}`;
		}
	}

	/** Scrolls once to the job card the URL's hash names, as a story's agent links to it. */
	async function land() {
		if (landed) return;
		landed = true;
		const at = decodeURIComponent(location.hash.slice(1));
		if (!at.startsWith('image-job-')) return;
		await tick();
		document.getElementById(at)?.scrollIntoView?.({ block: 'start' });
	}

	/** Reads a job's log from an offset; throws flai's refusal word for word. */
	async function readLog(id: string, offset: number): Promise<ImageLog> {
		const r = await api(`/api/images?job=${encodeURIComponent(id)}&offset=${offset}`);
		const body = await r.json().catch(() => ({}));
		if (!r.ok) throw new Error(body.error ?? r.statusText);
		return body as ImageLog;
	}

	/**
	 * Starts a new image job through flai, taking the images and jobs the host answers; throws flai's
	 * refusal word for word. While a job already makes the images, the page reads them again so that
	 * job shows.
	 */
	async function ensure() {
		const r = await api('/api/images/retry', { method: 'POST' });
		const body = await r.json().catch(() => ({}));
		if (r.ok) {
			data = body as ImagesState;
			return;
		}
		if (r.status === 501) {
			tooOld = body.error ?? r.statusText;
			return;
		}
		if (r.status === 403) refused = { action: body.action ?? 'agent', enable: body.enable ?? '' };
		if (body.job) void load();
		throw new Error(body.error ?? r.statusText);
	}

	async function make() {
		making = true;
		makeError = null;
		try {
			await ensure();
		} catch (e) {
			makeError = errorText(e);
		} finally {
			making = false;
		}
	}

	onMount(() => {
		void load();
		// asked again only while a job runs or waits: nothing else here changes without the page asking
		const timer = setInterval(() => {
			if (runs) void load();
		}, POLL_MS);
		return () => clearInterval(timer);
	});
</script>

{#snippet jobs(label: string, list: ImageJobRecord[], testid: string)}
	{#if list.length}
		<section class="mt-6" aria-label={label} data-testid={testid}>
			<h2 class="text-lg font-semibold">{label}</h2>
			<ul class="mt-2 space-y-3">
				{#each list as job (job.id)}
					<li>
						<ImageJob
							{job}
							{readLog}
							onretry={job.outcome === 'failed' ? ensure : undefined}
							disabled={!!refused}
						/>
					</li>
				{/each}
			</ul>
		</section>
	{/if}
{/snippet}

<svelte:head><title>Agent images · flaiover</title></svelte:head>

<h1 class="text-2xl font-semibold">Agent images</h1>
<p class="mt-1 text-sm text-muted">
	The images the project's agents run in, and the jobs flai serve runs on the host to pull, build,
	and layer them.
</p>

{#if error}
	<p class="mt-4 rounded border border-danger bg-danger-soft p-3 text-sm text-danger" role="alert">
		{error}
	</p>
{/if}
{#if tooOld}
	<div
		class="mt-4 rounded border border-warn bg-warn-soft p-3 text-sm text-warn"
		role="status"
		data-testid="images-too-old"
	>
		<p>This flai cannot show the agent images yet; upgrade flai on the host.</p>
		<p class="mt-1 text-xs">{tooOld}</p>
	</div>
{:else if data === null}
	{#if !error}
		<p class="mt-4 text-sm text-muted">Loading…</p>
	{/if}
{:else}
	{#if refused}
		<div
			class="mt-4 rounded border border-warn bg-warn-soft p-3 text-sm text-warn"
			role="status"
			data-testid="images-off"
		>
			<p>The {refused.action} host action is off, so the dashboard cannot start an image job.</p>
			{#if refused.enable}
				<p class="mt-1">
					On the host, in the project, run <code class="font-mono">{refused.enable}</code>, then
					reload this page.
				</p>
			{/if}
		</div>
	{/if}
	{#if data.unjudged}
		<p
			class="mt-4 rounded border border-warn bg-warn-soft p-3 text-sm text-warn"
			role="status"
			data-testid="images-unjudged"
		>
			{data.unjudged}
		</p>
	{/if}

	<section class="mt-6" aria-label="Images" data-testid="images">
		<h2 class="text-lg font-semibold">Images</h2>
		{#if data.flai}
			<p class="mt-1 text-sm text-muted" data-testid="images-flai">
				Each must hold flai <code class="font-mono">{data.flai}</code>.
			</p>
		{/if}
		{#if data.images.length}
			<table class="mt-2 text-sm">
				<thead>
					<tr class="text-left text-muted">
						<th class="pr-4">Image</th>
						<th class="pr-4">Ref</th>
						<th class="pr-4">Named by</th>
						<th class="pr-4">On the host</th>
						<th>State</th>
					</tr>
				</thead>
				<tbody>
					{#each data.images as image (image.ref)}
						<tr class="align-top" data-testid="image-row">
							<td class="pr-4">{kindText(image)} for {image.inner}</td>
							<td class="pr-4 font-mono" data-testid="image-ref">{image.ref}</td>
							<td class="pr-4">{sourceText(image.source)}</td>
							<td class="pr-4" data-testid="image-present">{image.present ? 'yes' : 'no'}</td>
							<td>
								<span
									class="rounded px-2 py-0.5 text-xs {stateClass(image.state)}"
									data-testid="image-state">{image.state}</span
								>
								{#if image.reasons?.length}
									<ul class="mt-1 text-xs text-muted" data-testid="image-reasons">
										{#each image.reasons as reason (`${reason.kind}:${reason.message}`)}
											<li>{reason.message}</li>
										{/each}
									</ul>
								{/if}
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		{:else}
			<p class="mt-2 text-sm text-muted" data-testid="images-none">
				The project's agents need no image flai makes.
			</p>
		{/if}
		{#if needed && !runs && !failed}
			<div class="mt-3 flex flex-wrap items-center gap-2">
				<button
					type="button"
					class="rounded bg-primary px-3 py-1 text-sm text-on-primary disabled:opacity-50"
					disabled={making || !!refused}
					onclick={make}
					data-testid="images-make">Make the images</button
				>
				{#if makeError}
					<p class="text-xs text-danger" role="alert" data-testid="images-make-error">
						{makeError}
					</p>
				{/if}
			</div>
		{/if}
	</section>

	{@render jobs('Running', data.running ? [data.running] : [], 'images-running')}
	{@render jobs('Waiting', data.queued, 'images-queued')}
	{@render jobs('Ended', data.ended, 'images-ended')}
	{#if nothing}
		<p class="mt-6 text-sm text-muted" data-testid="images-no-jobs">
			No image job has run for the project yet.
		</p>
	{/if}
{/if}
