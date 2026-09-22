<script lang="ts">
	// The dashboard managed from its own page (S-0081): what image and version the shared
	// container is running, a Check for updates that changes nothing, and Restart / Upgrade /
	// Stop, gated on the dashboard host action (flai serve enable dashboard, the same shape push
	// and agent already use). Restart and, when it swaps to a new image, upgrade stop the very
	// container answering the request that asked for them: the fetch is expected to end in a
	// network error, not a clean response, so those two treat one as "proceeding" and poll this
	// page's own status until the container is reachable again, rather than trusting the fetch to
	// resolve cleanly.
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
	let busy = $state<'check' | 'restart' | 'upgrade' | 'stop' | null>(null);
	let reconnecting = $state(false);
	let message = $state<string | null>(null);
	let checkResult = $state<{ running?: boolean; upgrade_available?: boolean } | null>(null);
	let failed = $state<string | null>(null);

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

	// A restart or an upgrade that swaps to a new image stops the container answering this very
	// request, so the fetch is raced against a short client-side timeout: whichever comes first,
	// a network error or the timeout, is treated the same as "proceeding", not as a failure.
	async function fireAndExpectMaybeNoAnswer(
		action: 'restart' | 'upgrade'
	): Promise<{ ok: true; body: Record<string, unknown> } | { ok: false; body?: unknown } | 'gone'> {
		const attempt = api('/api/dashboard', {
			method: 'POST',
			body: JSON.stringify({ action })
		})
			.then(async (r) => ({ ok: r.ok, body: await r.json().catch(() => ({})) }))
			.catch(() => 'gone' as const);
		const timeout = sleep(disconnectTimeoutMs).then(() => 'gone' as const);
		return Promise.race([attempt, timeout]);
	}

	async function waitForReconnect(previousImage: string | undefined) {
		reconnecting = true;
		const deadline = Date.now() + reconnectGiveUpMs;
		while (Date.now() < deadline) {
			await sleep(reconnectPollMs);
			if (await load()) {
				message =
					status?.image && status.image !== previousImage
						? `Reconnected — now running ${status.image}.`
						: 'Reconnected.';
				reconnecting = false;
				return;
			}
		}
		reconnecting = false;
		failed = 'Did not reconnect within a minute; on the host, flai dashboard status says more.';
	}

	async function act(action: 'check' | 'restart' | 'upgrade' | 'stop') {
		if (busy) return;
		busy = action;
		message = failed = null;
		checkResult = null;
		const previousImage = status?.image;
		try {
			if (action === 'check') {
				const r = await api('/api/dashboard', {
					method: 'POST',
					body: JSON.stringify({ action })
				});
				const body = await r.json().catch(() => ({}));
				if (!r.ok) {
					failed = body.error ?? r.statusText;
					return;
				}
				checkResult = body;
				return;
			}
			if (action === 'stop') {
				const r = await api('/api/dashboard', {
					method: 'POST',
					body: JSON.stringify({ action })
				});
				const body = await r.json().catch(() => ({}));
				if (!r.ok) {
					failed = body.error ?? r.statusText;
					return;
				}
				message =
					body.state === 'stopped'
						? `${body.container} stopped.`
						: body.state === 'still-running'
							? `Unregistered; still serving ${(body.serves ?? []).join(', ')}.`
							: `${body.container} was not running.`;
				await load();
				return;
			}
			const outcome = await fireAndExpectMaybeNoAnswer(action);
			if (outcome === 'gone') {
				await waitForReconnect(previousImage);
				return;
			}
			if (!outcome.ok) {
				failed = (outcome.body as { error?: string })?.error ?? 'the action failed';
				return;
			}
			// A clean response beat the container's own teardown: still confirm by reconnecting,
			// since "upgraded" can still be mid-swap when the answer arrives.
			await waitForReconnect(previousImage);
		} finally {
			busy = null;
		}
	}

	onMount(() => {
		void load();
	});
</script>

<div class="max-w-2xl rounded border border-line bg-surface p-4 text-sm" data-testid="host-panel">
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

			{#if status.dashboard_enabled}
				<p class="mt-3 flex flex-wrap gap-2">
					<button
						type="button"
						class="rounded border border-line px-2 py-1 disabled:opacity-60"
						onclick={() => act('check')}
						disabled={!!busy}
						data-testid="host-panel-check"
						>{busy === 'check' ? 'Checking…' : 'Check for updates'}</button
					>
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
				</p>
			{:else}
				<p class="mt-3 text-muted">
					Restart, upgrade, and stop from here are off; the operator turns them on in a shell on the
					host with
					<code class="rounded bg-ground px-1 text-ink">flai serve enable dashboard</code>.
				</p>
			{/if}
		{/if}
	{/if}
</div>
