<script lang="ts">
	// The processes flai host keeps (S-0107): serve, and each project's MCP server as one group, with
	// the state and version of each, Start / Stop / Restart per row, and one Check for upgrade and
	// one Upgrade for flai itself. Everything but the check is the host action (flai serve enable
	// host). The dashboard reaches the host through serve, so stopping or restarting serve, and an
	// upgrade (which has the host start its children again from the new binary), end the very
	// connection that asked: like the Dashboard area beside it, those treat a network error or a
	// timeout as the work going ahead and poll /api/host until it answers again. Serve stopped
	// cannot be started from here; the Stop confirmation says what brings it back. Versions (S-0298)
	// lists the published releases, a read, and installs a chosen one as Upgrade does, with its
	// version; one below a served project's flai.minimum is warned about, not refused (ADR-0117).
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
	/** One published flai release as host.versions lists it, newest first. */
	type Release = {
		version: string;
		tag: string;
		published?: string;
		installed: boolean;
		latest: boolean;
		below_minimum?: { project: string; minimum: string }[];
	};

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
	let releases = $state<Release[] | null>(null);
	let chosen = $state<Release | null>(null);

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
		// The dashboard outlives serve: when serve goes down mid-request, the route answers 502 "the
		// host flai went away before it answered", which for these writes is the work going ahead.
		const attempt = api('/api/host', { method: 'POST', body: JSON.stringify(body) })
			.then(async (r) =>
				mayEndConnection && r.status === 502
					? ('gone' as const)
					: ({ ok: r.ok, body: await r.json().catch(() => ({})) } as Outcome)
			)
			.catch(() => 'gone' as const);
		if (!mayEndConnection) return attempt;
		return Promise.race([attempt, sleep(disconnectTimeoutMs).then(() => 'gone' as const)]);
	}

	type Running = View & { running: true };

	// Polls until the host answers in a way `accept` takes, since what is still answering right
	// after the request may be the old process, before it went down.
	async function waitForReconnect(
		done: (v: Running) => string,
		accept: (v: Running) => boolean = () => true,
		gaveUp = 'The host did not answer again in time; on the host, flai host status says more.'
	) {
		reconnecting = true;
		const deadline = Date.now() + reconnectGiveUpMs;
		while (Date.now() < deadline) {
			await sleep(reconnectPollMs);
			if ((await load()) && view?.running && accept(view)) {
				message = done(view);
				reconnecting = false;
				return;
			}
		}
		reconnecting = false;
		failed = gaveUp;
	}

	function refusal(body: unknown): string {
		return (body as { error?: string })?.error ?? 'the action failed';
	}

	// Starts an action: what the last one left on the page is cleared.
	function begin(what: string) {
		busy = what;
		message = failed = null;
		checkResult = null;
		releases = chosen = null;
	}

	async function act(process: Process, action: Action) {
		if (busy) return;
		if (process === 'serve' && action === 'stop' && !confirmStopServe) {
			confirmStopServe = true;
			return;
		}
		confirmStopServe = false;
		begin(`${process}:${action}`);
		const endsConnection = process === 'serve' && action !== 'start';
		const servePid = children.find((k) => k.name === 'serve')?.pid;
		try {
			const outcome = await post({ action, process }, endsConnection);
			if (outcome === 'gone' && action === 'stop') {
				// the host is still there, but this page has no way left to reach it
				view = { running: false, reason: 'serve stopped', host_enabled: !!view?.host_enabled };
				message = 'serve stopped. Start it again in a shell on the host: flai host start serve';
				return;
			}
			if (outcome === 'gone') {
				await waitForReconnect(
					() => 'serve restarted; reconnected.',
					(v) => !servePid || v.children.find((k) => k.name === 'serve')?.pid !== servePid
				);
				return;
			}
			if (!outcome.ok) {
				failed = refusal(outcome.body);
				return;
			}
			await load();
			message = `${process === 'mcp' ? 'MCP' : process}: ${action === 'stop' ? 'stopped' : action === 'start' ? 'started' : 'restarted'}.`;
		} finally {
			busy = null;
		}
	}

	// The two reads: check for a newer flai, and versions, the published releases.
	async function read(action: 'check' | 'versions') {
		if (busy) return;
		begin(action);
		try {
			const outcome = await post({ action }, false);
			if (outcome === 'gone') failed = 'The host did not answer.';
			else if (!outcome.ok) failed = refusal(outcome.body);
			else if (action === 'check') checkResult = outcome.body as Check;
			else releases = Array.isArray(outcome.body) ? (outcome.body as Release[]) : [];
		} finally {
			busy = null;
		}
	}

	// version, a published release chosen from Versions, is installed in place of the newest.
	async function upgrade(version?: string) {
		if (busy) return;
		begin('upgrade');
		const before = view?.running ? view.version : undefined;
		try {
			const outcome = await post(
				version === undefined ? { action: 'upgrade' } : { action: 'upgrade', version },
				true
			);
			if (outcome !== 'gone' && !outcome.ok) {
				failed = refusal(outcome.body);
				return;
			}
			// flai host upgrade answers {upgrade, restarting}: nothing installed, nothing restarts
			if (outcome !== 'gone' && outcome.body.restarting === false) {
				const up = (outcome.body.upgrade ?? {}) as Check;
				message = `flai ${up.current ?? before ?? ''} is already the latest.`;
				return;
			}
			await waitForReconnect(
				(v) =>
					version === undefined
						? `Upgraded flai ${before ?? ''} to ${v.version}; serve and MCP restarted.`
						: `Deployed flai ${v.version} over ${before ?? ''}; serve and MCP restarted.`,
				(v) => v.version !== before,
				`The host did not come back on ${version === undefined ? 'a newer flai' : `flai ${version}`} in time; on the host, flai host status says more.`
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
	class="mt-4 max-w-3xl rounded border border-line bg-surface p-4 text-sm"
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
		<!-- Fixed layout (S-0113): Process and the actions keep set widths, and State, Version, and
		     Restarts split the rest equally, padded apart; the MCP project list wraps under Process. -->
		<table class="mt-2 w-full table-fixed text-left" data-testid="host-processes-table">
			<thead class="text-muted">
				<tr>
					<th class="w-28 py-1 pr-4 font-normal">Process</th>
					<th class="px-4 py-1 font-normal">State</th>
					<th class="px-4 py-1 font-normal">Version</th>
					<th class="px-4 py-1 font-normal">Restarts</th>
					{#if view.host_enabled}<th class="w-52 py-1 pl-4"></th>{/if}
				</tr>
			</thead>
			<tbody>
				{#each rows as row (row.process)}
					{@const state = stateOf(row.kids)}
					<tr class="border-t border-line align-top" data-testid={`host-processes-${row.process}`}>
						<td class="py-1 pr-4">
							<span class="font-mono">{row.label}</span>
							{#if row.process === 'mcp' && row.kids.length}
								<span
									class="block text-xs break-words text-muted"
									data-testid="host-processes-mcp-projects"
									>{row.kids.map((k) => projectName(k.root)).join(', ')}</span
								>
							{/if}
						</td>
						<td
							class="px-4 py-1"
							class:text-good={state === 'running'}
							class:text-warn={state !== 'running'}
							data-testid={`host-processes-${row.process}-state`}>{state}</td
						>
						<td class="px-4 py-1 font-mono" data-testid={`host-processes-${row.process}-version`}
							>{versionsOf(row.kids)}</td
						>
						<td class="px-4 py-1" data-testid={`host-processes-${row.process}-restarts`}
							>{row.kids.reduce((n, k) => n + k.restarts, 0)}</td
						>
						{#if view.host_enabled}
							<td class="py-1 pl-4 text-right whitespace-nowrap">
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
		{#if releases}
			{#if releases.length}
				<ul class="mt-2" aria-label="Published flai releases" data-testid="host-processes-versions">
					{#each releases as r (r.version)}
						{@const marks = [r.installed && 'installed', r.latest && 'newest'].filter(Boolean)}
						<li class="flex items-center gap-2 py-0.5" data-testid="host-processes-version">
							<span class="font-mono">{r.version}</span>
							{#if r.published}<span class="text-muted">{r.published.slice(0, 10)}</span>{/if}
							{#if marks.length}<span class="text-muted">{marks.join(', ')}</span>{/if}
							{#each r.below_minimum ?? [] as m (m.project)}
								<span class="text-warn">below {m.project}'s minimum {m.minimum}</span>
							{/each}
							{#if view.host_enabled && !r.installed}
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
				<p class="mt-2 text-muted" data-testid="host-processes-versions-none">
					No published releases were found.
				</p>
			{/if}
		{/if}
		{#if chosen}
			<p class="mt-2 text-warn" role="alert" data-testid="host-processes-confirm-deploy">
				Install flai {chosen.version} on the host? The host restarts serve and the MCP servers on it,
				as Upgrade does.
				{#each chosen.below_minimum ?? [] as m (m.project)}
					It is below {m.project}'s
					<code class="rounded bg-ground px-1 text-ink">flai.minimum</code>,
					{m.minimum}: flai serve leaves {m.project} unserved, and flai's commands refuse it, until flai
					is at least {m.minimum} again.
				{/each}
				<button
					type="button"
					class="ml-1 rounded border border-warn px-2 py-0.5"
					onclick={() => chosen && upgrade(chosen.version)}
					data-testid="host-processes-confirm-deploy-yes">Deploy {chosen.version}</button
				>
				<button
					type="button"
					class="ml-1 rounded border border-line px-2 py-0.5 text-ink"
					onclick={() => (chosen = null)}>Cancel</button
				>
			</p>
		{/if}

		<p class="mt-3 flex flex-wrap gap-2">
			<button
				type="button"
				class="rounded border border-line px-2 py-1 disabled:opacity-60"
				onclick={() => read('check')}
				disabled={!!busy}
				data-testid="host-processes-check"
				>{busy === 'check' ? 'Checking…' : 'Check for upgrade'}</button
			>
			<button
				type="button"
				class="rounded border border-line px-2 py-1 disabled:opacity-60"
				onclick={() => read('versions')}
				disabled={!!busy}
				data-testid="host-processes-versions-list"
				>{busy === 'versions' ? 'Listing…' : 'Versions'}</button
			>
			{#if view.host_enabled}
				<button
					type="button"
					class="rounded border border-line px-2 py-1 disabled:opacity-60"
					onclick={() => upgrade()}
					disabled={!!busy}
					data-testid="host-processes-upgrade"
					>{busy === 'upgrade' ? 'Upgrading…' : 'Upgrade'}</button
				>
			{/if}
		</p>
		{#if !view.host_enabled}
			<p class="mt-2 text-muted" data-testid="host-processes-off">
				Start, stop, restart, upgrade, and deploying a release from here are off; the operator turns
				them on in a shell on the host with
				<code class="rounded bg-ground px-1 text-ink">flai serve enable host</code>.
			</p>
		{/if}
	{/if}
</div>
