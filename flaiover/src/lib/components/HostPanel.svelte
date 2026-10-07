<script lang="ts">
	// The dashboard managed from its own page (S-0081): what image and version the shared
	// container is running, a Check for updates that changes nothing, and Restart / Upgrade /
	// Stop, gated on the dashboard host action (flai serve enable dashboard, the same shape push
	// and agent already use). Restart and, when it swaps to a new image, upgrade stop the very
	// container answering the request that asked for them: the fetch is expected to end in a
	// network error, not a clean response, so those two treat one as "proceeding" and poll this
	// page's own status until the container is reachable again, rather than trusting the fetch to
	// resolve cleanly. Versions (S-0298) lists the published releases, a read, and deploys a chosen
	// one as Upgrade does, with its tag: once, on this container, not pinned (ADR-0117, ADR-0118).
	import { api } from '$lib/api';
	import { onMount } from 'svelte';

	type Status = {
		container: string;
		running: boolean;
		image?: string;
		url?: string;
		serves?: string[];
		dashboard_enabled: boolean;
	};
	/** One published dashboard release as dashboard.versions lists it, newest first. */
	type Release = {
		version: string;
		tag: string;
		running: boolean;
		configured: boolean;
		latest: boolean;
	};

	// The wait before a stalled restart/upgrade POST is treated as "the container answering it is
	// gone, not failed", and the interval between reconnect polls after that: real values, cut to
	// milliseconds in tests, so the tests do not spend real seconds waiting on them.
	let {
		disconnectTimeoutMs = 4000,
		reconnectPollMs = 1500,
		reconnectGiveUpMs = 60_000
	}: {
		disconnectTimeoutMs?: number;
		reconnectPollMs?: number;
		reconnectGiveUpMs?: number;
	} = $props();

	let status = $state<Status | null>(null);
	let loadFailed = $state(false);
	let busy = $state<'check' | 'versions' | 'restart' | 'upgrade' | 'stop' | null>(null);
	let reconnecting = $state(false);
	let message = $state<string | null>(null);
	let checkResult = $state<{ running?: boolean; upgrade_available?: boolean } | null>(null);
	let releases = $state<Release[] | null>(null);
	let chosen = $state<Release | null>(null);
	let failed = $state<string | null>(null);

	// ADR-0118: how long a release chosen here lasts, said with the result of deploying one.
	const LASTS =
		'It keeps running through restarts until the next Upgrade, which goes to the configured tag, or a start after a stop.';

	async function load(): Promise<boolean> {
		try {
			const r = await api('/api/dashboard');
			if (!r.ok) return false;
			status = await r.json();
			loadFailed = false;
			return true;
		} catch {
			loadFailed = true;
			return false;
		}
	}

	function sleep(ms: number): Promise<void> {
		return new Promise((resolve) => setTimeout(resolve, ms));
	}

	// Restart always stops the container answering this very request; a swapping upgrade does on
	// success; stop does when this was the last project registered. All three fetches are raced
	// against a short client-side timeout: whichever comes first, a network error or the timeout,
	// is treated the same as "proceeding", not as a failure — found live (S-0081, T-0325): the
	// first version of this page called plain fetch for stop and always waited to reconnect after
	// upgrade, so an upgrade that changed nothing showed "Reconnected" for a container that was
	// never touched, and a stop that was genuinely the last project threw an unhandled rejection
	// instead of reporting anything.
	async function fireAndExpectMaybeNoAnswer(
		action: 'restart' | 'upgrade' | 'stop',
		tag?: string
	): Promise<{ ok: true; body: Record<string, unknown> } | { ok: false; body?: unknown } | 'gone'> {
		const attempt = api('/api/dashboard', {
			method: 'POST',
			body: JSON.stringify(tag === undefined ? { action } : { action, tag })
		})
			.then(async (r) => ({ ok: r.ok, body: await r.json().catch(() => ({})) }))
			.catch(() => 'gone' as const);
		const timeout = sleep(disconnectTimeoutMs).then(() => 'gone' as const);
		return Promise.race([attempt, timeout]);
	}

	async function waitForReconnect(previousImage: string | undefined, after = '') {
		reconnecting = true;
		const deadline = Date.now() + reconnectGiveUpMs;
		while (Date.now() < deadline) {
			await sleep(reconnectPollMs);
			if (await load()) {
				message =
					(status?.image && status.image !== previousImage
						? `Reconnected — now running ${status.image}.`
						: 'Reconnected.') + after;
				reconnecting = false;
				return;
			}
		}
		reconnecting = false;
		failed = 'Did not reconnect within a minute; on the host, flai dashboard status says more.';
	}

	// tag, given only with upgrade, deploys that published release in place of the configured one.
	async function act(action: 'check' | 'versions' | 'restart' | 'upgrade' | 'stop', tag?: string) {
		if (busy) return;
		busy = action;
		message = failed = null;
		checkResult = null;
		releases = chosen = null;
		const previousImage = status?.image;
		try {
			if (action === 'check' || action === 'versions') {
				const r = await api('/api/dashboard', {
					method: 'POST',
					body: JSON.stringify({ action })
				});
				const body = await r.json().catch(() => ({}));
				if (!r.ok) {
					failed = body.error ?? r.statusText;
					return;
				}
				if (action === 'check') checkResult = body;
				else releases = Array.isArray(body) ? body : [];
				return;
			}
			if (action === 'stop') {
				const outcome = await fireAndExpectMaybeNoAnswer('stop');
				if (outcome === 'gone') {
					// This was the last project: the container is genuinely gone, and with it
					// this page's own way of reaching it again. Nothing to reconnect to.
					message = 'Stopping: this was the last project, so the container is gone.';
					return;
				}
				if (!outcome.ok) {
					failed = (outcome.body as { error?: string })?.error ?? 'the action failed';
					return;
				}
				const body = outcome.body as { container?: string; state?: string; serves?: string[] };
				message =
					body.state === 'stopped'
						? `${body.container} stopped.`
						: body.state === 'still-running'
							? `Unregistered; still serving ${(body.serves ?? []).join(', ')}.`
							: `${body.container} was not running.`;
				if (body.state !== 'stopped') await load();
				return;
			}
			const outcome = await fireAndExpectMaybeNoAnswer(action, tag);
			const after = tag === undefined ? '' : ` ${LASTS}`;
			if (outcome === 'gone') {
				await waitForReconnect(previousImage, after);
				return;
			}
			if (!outcome.ok) {
				failed = (outcome.body as { error?: string })?.error ?? 'the action failed';
				return;
			}
			// A clean response beat the container's own teardown. Restart always cycles the
			// container, so it always reconnects; upgrade only touched anything if it actually
			// swapped — "up-to-date" or "started" never stopped what was already answering.
			const done = outcome.body as { outcome?: string; to?: string };
			if (action === 'restart' || done.outcome === 'upgraded') {
				await waitForReconnect(previousImage, after);
				return;
			}
			const container = status?.container ?? 'flaiover';
			message =
				done.outcome === 'up-to-date'
					? tag === undefined
						? `${container} is already running the latest.`
						: `${container} is already running ${done.to ?? tag}.`
					: `Started.${after}`;
		} finally {
			busy = null;
		}
	}

	onMount(() => {
		void load();
	});
</script>

<div class="max-w-3xl rounded border border-line bg-surface p-4 text-sm" data-testid="host-panel">
	<h2 class="font-semibold text-ink">Dashboard</h2>

	{#if !status && !loadFailed}
		<p class="mt-2 text-muted">Loading…</p>
	{:else if loadFailed && !reconnecting}
		<p class="mt-2 text-warn" data-testid="host-panel-unreachable">
			Could not reach the dashboard.
		</p>
	{:else if status}
		<p class="mt-2" data-testid="host-panel-status">
			{#if status.running}
				<span class="font-mono">{status.container}</span> running
				{#if status.image}<span class="font-mono">{status.image}</span>{/if}
				{#if status.serves?.length}, serving {status.serves.join(', ')}{/if}.
			{:else}
				<span class="font-mono">{status.container}</span> is not running.
			{/if}
		</p>

		{#if reconnecting}
			<p class="mt-2 text-muted" role="status" data-testid="host-panel-reconnecting">
				Reconnecting…
			</p>
		{:else}
			{#if message}
				<p class="mt-2 text-good" role="status" data-testid="host-panel-message">{message}</p>
			{/if}
			{#if checkResult}
				<p class="mt-2" data-testid="host-panel-check-result">
					{checkResult.running === false
						? 'Not running; Restart starts it.'
						: checkResult.upgrade_available
							? 'An update is available.'
							: 'Already running the latest.'}
				</p>
			{/if}
			{#if failed}
				<p class="mt-2 text-warn" role="status" data-testid="host-panel-failed">{failed}</p>
			{/if}
			{#if releases}
				{#if releases.length}
					<ul
						class="mt-2"
						aria-label="Published dashboard releases"
						data-testid="host-panel-versions"
					>
						{#each releases as r (r.version)}
							{@const marks = [
								r.running && 'running',
								r.latest && 'newest',
								r.configured && 'configured'
							].filter(Boolean)}
							<li class="flex items-center gap-2 py-0.5" data-testid="host-panel-version">
								<span class="font-mono">{r.version}</span>
								{#if marks.length}<span class="text-muted">{marks.join(', ')}</span>{/if}
								{#if status.dashboard_enabled && !r.running}
									<button
										type="button"
										class="ml-auto rounded border border-line px-2 py-0.5 disabled:opacity-60"
										onclick={() => (chosen = r)}
										disabled={!!busy}
										aria-label={`Deploy ${r.version}`}>Deploy</button
									>
								{/if}
							</li>
						{/each}
					</ul>
				{:else}
					<p class="mt-2 text-muted" data-testid="host-panel-versions-none">
						No published releases were found.
					</p>
				{/if}
			{/if}
			{#if chosen}
				<p class="mt-2 text-warn" role="alert" data-testid="host-panel-confirm-deploy">
					Deploy {status.container}
					{chosen.version}? It applies to this container, not the configuration. {LASTS} This page has
					no control to pin a release; on the host,
					<code class="rounded bg-ground px-1 text-ink">flai config set dashboard.tag</code> does.
					<button
						type="button"
						class="ml-1 rounded border border-warn px-2 py-0.5"
						onclick={() => chosen && act('upgrade', chosen.version)}
						data-testid="host-panel-confirm-deploy-yes">Deploy {chosen.version}</button
					>
					<button
						type="button"
						class="ml-1 rounded border border-line px-2 py-0.5 text-ink"
						onclick={() => (chosen = null)}>Cancel</button
					>
				</p>
			{/if}

			<p class="mt-3 flex flex-wrap gap-2">
				{#if status.dashboard_enabled}
					<button
						type="button"
						class="rounded border border-line px-2 py-1 disabled:opacity-60"
						onclick={() => act('check')}
						disabled={!!busy}
						data-testid="host-panel-check"
						>{busy === 'check' ? 'Checking…' : 'Check for updates'}</button
					>
				{/if}
				<button
					type="button"
					class="rounded border border-line px-2 py-1 disabled:opacity-60"
					onclick={() => act('versions')}
					disabled={!!busy}
					data-testid="host-panel-versions-list"
					>{busy === 'versions' ? 'Listing…' : 'Versions'}</button
				>
				{#if status.dashboard_enabled}
					<button
						type="button"
						class="rounded border border-line px-2 py-1 disabled:opacity-60"
						onclick={() => act('restart')}
						disabled={!!busy}
						data-testid="host-panel-restart"
						>{busy === 'restart' ? 'Restarting…' : 'Restart'}</button
					>
					<button
						type="button"
						class="rounded border border-line px-2 py-1 disabled:opacity-60"
						onclick={() => act('upgrade')}
						disabled={!!busy}
						data-testid="host-panel-upgrade">{busy === 'upgrade' ? 'Upgrading…' : 'Upgrade'}</button
					>
					<button
						type="button"
						class="rounded border border-warn px-2 py-1 text-warn disabled:opacity-60"
						onclick={() => act('stop')}
						disabled={!!busy}
						data-testid="host-panel-stop">{busy === 'stop' ? 'Stopping…' : 'Stop'}</button
					>
				{/if}
			</p>
			{#if !status.dashboard_enabled}
				<p class="mt-3 text-muted">
					Restart, upgrade, deploying a release, and stop from here are off; the operator turns them
					on in a shell on the host with
					<code class="rounded bg-ground px-1 text-ink">flai serve enable dashboard</code>.
				</p>
			{/if}
		{/if}
	{/if}
</div>
