<script lang="ts">
	// The processes flai host keeps (S-0107): serve, and each project's MCP server as one group, with
	// the state and version of each, Start / Stop / Restart per row, and one Check for upgrade and
	// one Upgrade for flai itself. Everything but the check is the host action (flai serve enable
	// host). The dashboard reaches the host through serve, so stopping or restarting serve, and an
	// upgrade (which has the host start its children again from the new binary), end the very
	// connection that asked: like the Dashboard area beside it, those treat a network error or a
	// timeout as the work going ahead and poll /api/host until it answers again. Serve stopped
	// cannot be started from here; the Stop confirmation says what brings it back.
	import { api } from '$lib/api';
	import { onMount } from 'svelte';

	type Child = {
		name: string;
		root?: string;
		state: string;
		pid?: number;
		version?: string;
		restarts: number;
		last_error?: string;
	};
	type View =
		| {
				running: true;
				pid: number;
				version: string;
				children: Child[];
				host_enabled: boolean;
		  }
		| { running: false; reason: string; host_enabled: boolean };
	type Process = 'serve' | 'mcp';
	type Action = 'start' | 'stop' | 'restart';
	type Check = { current?: string; latest?: string; up_to_date?: boolean };

	let {
		disconnectTimeoutMs = 4000,
		reconnectPollMs = 1500,
		reconnectGiveUpMs = 90_000
	}: {
		disconnectTimeoutMs?: number;
		reconnectPollMs?: number;
		reconnectGiveUpMs?: number;
	} = $props();

	let view = $state<View | null>(null);
	let loadFailed = $state(false);
	let busy = $state<string | null>(null);
	let reconnecting = $state(false);
	let confirmStopServe = $state(false);
	let message = $state<string | null>(null);
	let failed = $state<string | null>(null);
	let checkResult = $state<Check | null>(null);

	const children = $derived(view?.running ? view.children : []);
	const rows = $derived([
		{
			process: 'serve' as Process,
			label: 'serve',
			kids: children.filter((c) => c.name === 'serve')
		},
		{ process: 'mcp' as Process, label: 'MCP', kids: children.filter((c) => c.name === 'mcp') }
	]);

	async function load(): Promise<boolean> {
		try {
			const r = await api('/api/host');
			if (!r.ok) return false;
			view = await r.json();
			loadFailed = false;
			return !!view?.running;
		} catch {
			loadFailed = true;
			return false;
		}
	}

	function sleep(ms: number): Promise<void> {
		return new Promise((resolve) => setTimeout(resolve, ms));
	}

	type Outcome =
		{ ok: true; body: Record<string, unknown> } | { ok: false; body: unknown } | 'gone';

	function post(body: Record<string, unknown>, mayEndConnection: boolean): Promise<Outcome> {
		const attempt = api('/api/host', { method: 'POST', body: JSON.stringify(body) })
			.then(async (r) => ({ ok: r.ok, body: await r.json().catch(() => ({})) }) as Outcome)
			.catch(() => 'gone' as const);
		if (!mayEndConnection) return attempt;
		return Promise.race([attempt, sleep(disconnectTimeoutMs).then(() => 'gone' as const)]);
	}

	async function waitForReconnect(done: (v: View & { running: true }) => string) {
		reconnecting = true;
		const deadline = Date.now() + reconnectGiveUpMs;
		while (Date.now() < deadline) {
			await sleep(reconnectPollMs);
			if ((await load()) && view?.running) {
				message = done(view);
				reconnecting = false;
				return;
			}
		}
		reconnecting = false;
		failed = 'The host did not answer again in time; on the host, flai host status says more.';
	}

	function refusal(body: unknown): string {
		return (body as { error?: string })?.error ?? 'the action failed';
	}

	async function act(process: Process, action: Action) {
		if (busy) return;
		if (process === 'serve' && action === 'stop' && !confirmStopServe) {
			confirmStopServe = true;
			return;
		}
		confirmStopServe = false;
		busy = `${process}:${action}`;
		message = failed = null;
		checkResult = null;
		const endsConnection = process === 'serve' && action !== 'start';
		try {
			const outcome = await post({ action, process }, endsConnection);
			if (outcome === 'gone' && action === 'stop') {
				// the host is still there, but this page has no way left to reach it
				view = { running: false, reason: 'serve stopped', host_enabled: !!view?.host_enabled };
				message = 'serve stopped. Start it again in a shell on the host: flai host start serve';
				return;
			}
			if (outcome === 'gone') {
				await waitForReconnect(() => 'serve restarted; reconnected.');
				return;
			}
			if (!outcome.ok) {
				failed = refusal(outcome.body);
				return;
			}
			message = `${process === 'mcp' ? 'MCP' : process}: ${action === 'stop' ? 'stopped' : action === 'start' ? 'started' : 'restarted'}.`;
			await load();
		} finally {
			busy = null;
		}
	}

	async function check() {
		if (busy) return;
		busy = 'check';
		message = failed = null;
		checkResult = null;
		try {
			const outcome = await post({ action: 'check' }, false);
			if (outcome === 'gone') failed = 'The host did not answer.';
			else if (!outcome.ok) failed = refusal(outcome.body);
			else checkResult = outcome.body as Check;
		} finally {
			busy = null;
		}
	}

	async function upgrade() {
		if (busy) return;
		busy = 'upgrade';
		message = failed = null;
		checkResult = null;
		const before = view?.running ? view.version : undefined;
		try {
			const outcome = await post({ action: 'upgrade' }, true);
			if (outcome !== 'gone' && !outcome.ok) {
				failed = refusal(outcome.body);
				return;
			}
			if (outcome !== 'gone' && (outcome.body as Check).up_to_date) {
				message = `flai ${before ?? ''} is already the latest.`;
				return;
			}
			await waitForReconnect((v) =>
				v.version !== before
					? `Upgraded flai ${before ?? ''} to ${v.version}; serve and MCP restarted.`
					: `Reconnected; flai is still ${v.version}.`
			);
		} finally {
			busy = null;
		}
	}

	function stateOf(kids: Child[]): string {
		if (!kids.length) return 'none';
		const states = [...new Set(kids.map((k) => k.state))];
		return states.length === 1 ? states[0] : states.join(', ');
	}

	function versionsOf(kids: Child[]): string {
		return [...new Set(kids.map((k) => k.version).filter(Boolean))].join(', ');
	}

	function projectName(root: string | undefined): string {
		return root?.split(/[\\/]/).filter(Boolean).pop() ?? '';
	}

	onMount(() => {
		void load();
	});
</script>

<div
	class="mt-4 max-w-2xl rounded border border-line bg-surface p-4 text-sm"
	data-testid="host-processes"
>
	<h2 class="font-semibold text-ink">
		flai host
		{#if view?.running}
			<span class="font-normal text-muted" data-testid="host-processes-host"
				>{view.version}, pid {view.pid}</span
			>
		{/if}
	</h2>

	{#if !view && !loadFailed}
		<p class="mt-2 text-muted">Loading…</p>
	{:else if reconnecting}
		<p class="mt-2 text-muted" role="status" data-testid="host-processes-reconnecting">
			Reconnecting…
		</p>
	{:else if !view || !view.running}
		<p class="mt-2 text-warn" data-testid="host-processes-none">
			{#if message}{message}{:else}
				No flai host is answering{view && !view.running ? ` (${view.reason})` : ''}. Start one in a
				shell on the host with
				<code class="rounded bg-ground px-1 text-ink">flai host start</code>.
			{/if}
		</p>
	{:else}
		<table class="mt-2 w-full text-left">
			<thead class="text-muted">
				<tr>
					<th class="py-1 font-normal">Process</th>
					<th class="py-1 font-normal">State</th>
					<th class="py-1 font-normal">Version</th>
					<th class="py-1 font-normal">Restarts</th>
					{#if view.host_enabled}<th class="py-1"></th>{/if}
				</tr>
			</thead>
			<tbody>
				{#each rows as row (row.process)}
					{@const state = stateOf(row.kids)}
					<tr class="border-t border-line align-top" data-testid={`host-processes-${row.process}`}>
						<td class="py-1">
							<span class="font-mono">{row.label}</span>
							{#if row.process === 'mcp' && row.kids.length}
								<span class="block text-xs text-muted" data-testid="host-processes-mcp-projects"
									>{row.kids.map((k) => projectName(k.root)).join(', ')}</span
								>
							{/if}
						</td>
						<td
							class="py-1"
							class:text-good={state === 'running'}
							class:text-warn={state !== 'running'}
							data-testid={`host-processes-${row.process}-state`}>{state}</td
						>
						<td class="py-1 font-mono" data-testid={`host-processes-${row.process}-version`}
							>{versionsOf(row.kids)}</td
						>
						<td class="py-1">{row.kids.reduce((n, k) => n + k.restarts, 0)}</td>
						{#if view.host_enabled}
							<td class="py-1 text-right whitespace-nowrap">
								{#each ['start', 'stop', 'restart'] as const as action (action)}
									<button
										type="button"
										class="ml-1 rounded border border-line px-2 py-0.5 capitalize disabled:opacity-60"
										class:border-warn={action === 'stop'}
										class:text-warn={action === 'stop'}
										onclick={() => act(row.process, action)}
										disabled={!!busy ||
											!row.kids.length ||
											(action === 'start' && row.kids.every((k) => k.state === 'running')) ||
											(action === 'stop' && row.kids.every((k) => k.state === 'stopped'))}
										data-testid={`host-processes-${row.process}-${action}`}
										>{busy === `${row.process}:${action}` ? `${action}…` : action}</button
									>
								{/each}
							</td>
						{/if}
					</tr>
					{#each row.kids.filter((k) => k.last_error) as k (k.root ?? k.name)}
						<tr>
							<td colspan="5" class="pb-1 text-xs text-warn"
								>{k.root ? projectName(k.root) + ': ' : ''}{k.last_error}</td
							>
						</tr>
					{/each}
				{/each}
			</tbody>
		</table>

		{#if confirmStopServe}
			<p class="mt-2 text-warn" role="alert" data-testid="host-processes-confirm-stop">
				This dashboard reaches the host through serve, so it cannot start serve again: only
				<code class="rounded bg-ground px-1 text-ink">flai host start serve</code> in a shell on the
				host does.
				<button
					type="button"
					class="ml-1 rounded border border-warn px-2 py-0.5"
					onclick={() => act('serve', 'stop')}
					data-testid="host-processes-confirm-stop-yes">Stop serve</button
				>
				<button
					type="button"
					class="ml-1 rounded border border-line px-2 py-0.5 text-ink"
					onclick={() => (confirmStopServe = false)}>Cancel</button
				>
			</p>
		{/if}
		{#if message}
			<p class="mt-2 text-good" role="status" data-testid="host-processes-message">{message}</p>
		{/if}
		{#if checkResult}
			<p class="mt-2" data-testid="host-processes-check-result">
				{checkResult.up_to_date
					? `flai ${checkResult.current ?? ''} is the latest.`
					: `flai ${checkResult.latest ?? ''} is available (running ${checkResult.current ?? ''}).`}
			</p>
		{/if}
		{#if failed}
			<p class="mt-2 text-warn" role="status" data-testid="host-processes-failed">{failed}</p>
		{/if}

		<p class="mt-3 flex flex-wrap gap-2">
			<button
				type="button"
				class="rounded border border-line px-2 py-1 disabled:opacity-60"
				onclick={check}
				disabled={!!busy}
				data-testid="host-processes-check"
				>{busy === 'check' ? 'Checking…' : 'Check for upgrade'}</button
			>
			{#if view.host_enabled}
				<button
					type="button"
					class="rounded border border-line px-2 py-1 disabled:opacity-60"
					onclick={upgrade}
					disabled={!!busy}
					data-testid="host-processes-upgrade"
					>{busy === 'upgrade' ? 'Upgrading…' : 'Upgrade'}</button
				>
			{/if}
		</p>
		{#if !view.host_enabled}
			<p class="mt-2 text-muted" data-testid="host-processes-off">
				Start, stop, restart, and upgrade from here are off; the operator turns them on in a shell
				on the host with
				<code class="rounded bg-ground px-1 text-ink">flai serve enable host</code>.
			</p>
		{/if}
	{/if}
</div>
