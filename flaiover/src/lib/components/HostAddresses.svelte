<script lang="ts">
	// The hosts the dashboard accepts besides its own address (S-0489, ADR-0193), so that it can be
	// reached through a tunnel or a domain: the host's dashboard.allowed_hosts from settings.get, and
	// the list the running container was started with from dashboard.status. While the settings host
	// action is on for every project the list is edited here through settings.dashboard_hosts, which
	// runs flai config set on the host; until then it is shown with the command that allows it. The
	// container keeps the list it was started with, so a saved list applies on the next restart.
	import { api } from '$lib/api';
	import { allowed, type SettingsView } from '$lib/settings';
	import { onMount } from 'svelte';

	/** What this area reads of GET /api/dashboard: allowed_hosts only for a running container. */
	type Dashboard = { running: boolean; allowed_hosts?: string; dashboard_enabled: boolean };

	let view = $state<SettingsView | null>(null);
	let dashboard = $state<Dashboard | null>(null);
	let error = $state<string | null>(null);
	let hosts = $state('');
	let saving = $state(false);
	let refused = $state<string | null>(null);
	let saved = $state(false);

	const PURPOSE =
		"The address a request comes to is always accepted. List a tunnel's or a domain's host name to reach the dashboard through it.";
	const RESTART = 'Restart in the Dashboard area above applies it.';
	const RESTART_OFF =
		'Restart in the Dashboard area above applies it once the operator turns on the dashboard host action, on the host, with flai serve enable dashboard.';

	/** A normalized list, a,b, as it reads, a, b, or none. */
	const named = (list: string | null) => list?.split(',').filter(Boolean).join(', ') || 'none';

	const configured = $derived(view?.host?.allowed_hosts ?? '');
	const can = $derived(view ? allowed(view, 'dashboard_hosts') : { ok: false, enable: '' });
	/** The running container's list, or null when it is not running or its flai does not say. */
	const running = $derived(
		dashboard?.running && typeof dashboard.allowed_hosts === 'string'
			? dashboard.allowed_hosts
			: null
	);
	const differs = $derived(view !== null && running !== null && running !== configured);

	function savedNote(): string {
		return `Saved. The running dashboard keeps its old list until it is restarted: ${dashboard?.dashboard_enabled ? RESTART : RESTART_OFF}`;
	}

	function differNote(): string {
		return `The running dashboard was started with ${named(running)}; a restart applies ${named(configured)}.`;
	}

	async function loadSettings() {
		try {
			const r = await api('/api/settings');
			const body = await r.json().catch(() => ({}));
			if (!r.ok) {
				error = body.error ?? r.statusText;
				return;
			}
			error = null;
			view = body;
			hosts = body.host?.allowed_hosts?.split(',').join(', ') ?? '';
		} catch (e) {
			error = e instanceof Error ? e.message : String(e);
		}
	}

	async function loadDashboard() {
		try {
			const r = await api('/api/dashboard');
			if (r.ok) dashboard = await r.json();
		} catch {
			// the Dashboard area above says it could not reach the dashboard
		}
	}

	async function save() {
		if (saving) return;
		saving = true;
		refused = null;
		saved = false;
		try {
			const r = await api('/api/settings', {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ kind: 'dashboard_hosts', hosts })
			});
			const body = await r.json().catch(() => ({}));
			if (!r.ok) {
				refused = body.error ?? r.statusText;
				return;
			}
			saved = true;
			await loadSettings();
		} catch (e) {
			refused = e instanceof Error ? e.message : String(e);
		} finally {
			saving = false;
		}
	}

	onMount(() => {
		void loadSettings();
		void loadDashboard();
	});
</script>

<div
	class="mt-4 max-w-3xl rounded border border-line bg-surface p-4 text-sm"
	data-testid="host-addresses"
>
	<h2 class="font-semibold text-ink">Addresses</h2>
	<p class="mt-1 text-xs text-muted">{PURPOSE}</p>

	{#if error}
		<p class="mt-2 text-warn" role="alert" data-testid="host-addresses-error">{error}</p>
	{:else if !view}
		<p class="mt-2 text-muted">Loading…</p>
	{:else if can.ok}
		<label class="mt-2 block">
			<span class="mb-1 block text-xs text-muted"
				>host names, each with an optional :port, separated by commas</span
			>
			<input
				class="w-full rounded border border-line-strong bg-surface px-2 py-1 font-mono text-xs"
				bind:value={hosts}
				placeholder="flow.example.com, board.example.net:8443"
				disabled={saving}
				data-testid="host-addresses-field"
			/>
		</label>
		{#if refused}
			<p class="mt-1 text-xs text-danger" role="alert" data-testid="host-addresses-refused">
				{refused}
			</p>
		{/if}
		{@render differ()}
		<p class="mt-2">
			<button
				type="button"
				class="rounded border border-line px-2 py-1 disabled:opacity-60"
				onclick={save}
				disabled={saving}
				data-testid="host-addresses-save">{saving ? 'Saving…' : 'Save'}</button
			>
		</p>
		{#if saved}
			<p class="mt-2 text-xs text-good" role="status" data-testid="host-addresses-saved">
				{savedNote()}
			</p>
		{/if}
	{:else}
		<p class="mt-2 font-mono" data-testid="host-addresses-list">{named(configured)}</p>
		{@render differ()}
		<div class="mt-2 text-xs text-muted" data-testid="host-addresses-gate">
			<p>Read-only here. On the host, run this to change it from the dashboard:</p>
			<p class="mt-1"><code class="rounded bg-ground px-1 text-ink">{can.enable}</code></p>
		</div>
	{/if}
</div>

{#snippet differ()}
	{#if differs}
		<p class="mt-1 text-xs text-warn" data-testid="host-addresses-differs">{differNote()}</p>
	{/if}
{/snippet}
