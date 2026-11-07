<script lang="ts">
	// /ops/klusters: the project's kluster (S-0580), flai's kluster.status through /api/kluster. It
	// shows the cluster the manifest's kluster: block names, its profile, provider, whether it runs,
	// and its age, and brings it up or takes it down through flai's kluster.up and kluster.down, Down
	// only once the operator confirms. flai runs either in the background and answers at once, so
	// while one runs the page says so, shows the end of its log, and asks again every few seconds
	// until it ends; then it shows how the last one ended. flai's refusal or error is shown word for
	// word. With the kluster host action off, or no kluster: block, the page says which in place of
	// the buttons: the command that turns the action on and the Features page, or the block to add
	// and the design that describes it.
	import type { KlusterJob, KlusterStatus } from '$lib/server/repo';
	import { api } from '$lib/api';
	import { onMount } from 'svelte';
	import { resolve } from '$app/paths';

	const ENABLE = 'flai serve enable kluster';
	const DESIGN = 'design/system/kluster.md';
	// How often the page asks again while an up or down runs; either takes minutes.
	const POLL_MS = 3000;
	const button = 'rounded px-3 py-1 text-sm disabled:opacity-50';

	let status = $state<KlusterStatus | null>(null);
	let error = $state<string | null>(null);
	// The host's flai answered 501: its message, which says what it cannot do.
	let tooOld = $state<string | null>(null);
	// The host action was refused (403) on a write: what the operator runs to turn it on.
	let refused = $state<string | null>(null);
	// flai's refusal or error for the last Up or Down, word for word.
	let actionError = $state<string | null>(null);
	let busy = $state(false);
	let confirming = $state(false);

	const off = $derived(refused ?? (status && !status.enabled ? status.enable || ENABLE : null));
	const cluster = $derived(status?.cluster);
	const running = $derived(status?.running);
	const shown = $derived(
		!cluster ? null : cluster.running ? 'Running' : cluster.exists ? 'Stopped' : 'Not created'
	);

	const errorText = (e: unknown) => (e instanceof Error ? e.message : String(e));
	const jobLine = (j: KlusterJob) => `kluster ${j.action} of the cluster ${j.cluster}`;

	/** Reads the kluster's status; a failed read keeps what was last shown. */
	async function load() {
		try {
			const r = await api('/api/kluster');
			const body = await r.json().catch(() => ({}));
			if (r.status === 501) {
				tooOld = body.error ?? r.statusText;
				return;
			}
			if (!r.ok) throw new Error(body.error ?? r.statusText);
			status = body as KlusterStatus;
			tooOld = error = null;
			if (status.enabled) refused = null;
		} catch (e) {
			error = `Could not read the kluster: ${errorText(e)}`;
		}
	}

	/** Brings the cluster up or takes it down through flai, taking the status the host answers. */
	async function control(action: 'up' | 'down') {
		confirming = false;
		busy = true;
		actionError = null;
		try {
			const r = await api(`/api/kluster/${action}`, { method: 'POST' });
			const body = await r.json().catch(() => ({}));
			if (r.ok) {
				status = body as KlusterStatus;
				return;
			}
			if (r.status === 501) tooOld = body.error ?? r.statusText;
			else {
				if (r.status === 403) refused = body.enable ?? ENABLE;
				actionError = body.error ?? r.statusText;
			}
		} catch (e) {
			actionError = errorText(e);
		} finally {
			busy = false;
		}
	}

	onMount(() => {
		void load();
		// asked again only while a job runs: nothing else here changes without the page asking
		const timer = setInterval(() => {
			if (status?.running) void load();
		}, POLL_MS);
		return () => clearInterval(timer);
	});
</script>

{#snippet job(j: KlusterJob, now: boolean)}
	<div
		class="mt-4 rounded border border-line p-3 text-sm"
		data-testid={now ? 'kluster-running' : 'kluster-last'}
	>
		{#if now}
			<p role="status">{jobLine(j)} runs, since {j.started_at}.</p>
		{:else}
			<p>
				Last: {jobLine(j)}
				<span class={j.outcome === 'failed' ? 'text-danger' : 'text-good-strong'}>{j.outcome}</span>
				at {j.ended_at}.
			</p>
			{#if j.error}
				<p class="mt-1 text-danger" data-testid="kluster-last-error">{j.error}</p>
			{/if}
		{/if}
		<p class="mt-1 text-xs text-muted">Its output: <code class="font-mono">{j.log}</code></p>
		{#if j.log_tail}
			<pre
				class="mt-2 max-h-64 overflow-auto rounded bg-surface p-2 font-mono text-xs"
				data-testid="kluster-log">{j.log_tail}</pre>
		{/if}
	</div>
{/snippet}

<svelte:head><title>Klusters · flaiover</title></svelte:head>

<h1 class="text-2xl font-semibold">Klusters</h1>
<p class="mt-1 text-sm text-muted">
	The project's local Kubernetes cluster, made by kluster from the manifest's kluster: block. flai
	on the host brings it up and takes it down.
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
		data-testid="kluster-too-old"
	>
		<p>This flai cannot show the kluster yet; upgrade flai on the host.</p>
		<p class="mt-1 text-xs">{tooOld}</p>
	</div>
{:else if status === null}
	{#if !error}
		<p class="mt-4 text-sm text-muted">Loading…</p>
	{/if}
{:else if !status.block}
	<div
		class="mt-4 rounded border border-warn bg-warn-soft p-3 text-sm text-warn"
		role="status"
		data-testid="kluster-no-block"
	>
		<p>This project has no kluster: block, so flai runs no cluster for it.</p>
		{#if status.no_block}
			<p class="mt-1">{status.no_block}</p>
		{/if}
		<p class="mt-1">
			Add the block to <code class="font-mono">system-flow.yaml</code>, as
			<a class="underline" href={resolve('/docs/[...path]', { path: DESIGN })}>{DESIGN}</a> describes.
		</p>
	</div>
{:else}
	<section class="mt-6" aria-labelledby="kluster-heading">
		<h2 id="kluster-heading" class="text-lg font-semibold">{cluster?.name ?? 'The cluster'}</h2>
		{#if cluster}
			<dl
				class="mt-2 grid grid-cols-[max-content_1fr] gap-x-4 gap-y-1 text-sm"
				data-testid="kluster-cluster"
			>
				<dt class="text-muted">Name</dt>
				<dd>{cluster.name}</dd>
				<dt class="text-muted">Profile</dt>
				<dd>{cluster.profile || '—'}</dd>
				<dt class="text-muted">Provider</dt>
				<dd>{cluster.provider}</dd>
				<dt class="text-muted">State</dt>
				<dd data-testid="kluster-state">{shown}</dd>
				<dt class="text-muted">Age</dt>
				<dd>{cluster.age || '—'}</dd>
			</dl>
		{/if}
		{#if status.status_error}
			<p
				class="mt-4 rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
				role="alert"
				data-testid="kluster-status-error"
			>
				{status.status_error}
			</p>
		{/if}

		{#if off}
			<div
				class="mt-4 rounded border border-warn bg-warn-soft p-3 text-sm text-warn"
				role="status"
				data-testid="kluster-off"
			>
				<p>
					The kluster host action is off, so the dashboard cannot bring the cluster up or take it
					down.
				</p>
				<p class="mt-1">
					Turn it on from the
					<a class="underline" href={resolve('/settings/[page]', { page: 'features' })}
						>Features page</a
					>, or run <code class="font-mono">{off}</code> on the host, in the project.
				</p>
			</div>
		{:else}
			<div class="mt-4 flex flex-wrap items-center gap-2">
				<button
					type="button"
					class="{button} bg-primary text-on-primary"
					disabled={busy || !!running}
					onclick={() => control('up')}>Up</button
				>
				<button
					type="button"
					class="{button} bg-danger text-on-primary"
					disabled={busy || !!running || confirming}
					onclick={() => (confirming = true)}>Down</button
				>
			</div>
			{#if confirming}
				<div
					class="mt-2 rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
					role="alertdialog"
					aria-label="Take the cluster down"
					data-testid="kluster-confirm-down"
				>
					<p>
						Take the cluster {cluster?.name ?? ''} down? kluster deletes it, with everything deployed
						in it.
					</p>
					<div class="mt-2 flex gap-2">
						<button
							type="button"
							class="{button} bg-danger text-on-primary"
							onclick={() => control('down')}>Take it down</button
						>
						<button
							type="button"
							class="{button} border border-line"
							onclick={() => (confirming = false)}>Keep it</button
						>
					</div>
				</div>
			{/if}
		{/if}
		{#if actionError}
			<p
				class="mt-4 rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
				role="alert"
				data-testid="kluster-action-error"
			>
				{actionError}
			</p>
		{/if}
		{#if running}
			{@render job(running, true)}
		{:else if status.last}
			{@render job(status.last, false)}
		{/if}
	</section>
{/if}
