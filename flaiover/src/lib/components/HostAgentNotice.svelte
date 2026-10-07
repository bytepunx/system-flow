<script lang="ts">
	// What flai on the host did about starting an agent when a story became ready (S-0079): that it
	// started one, for which story, by which command and when; that it could not; or why a ready
	// story is waiting. It shows nothing while the operator has not enabled it. Starting, stopping,
	// and configuring are done on the host; nothing here can.
	// The board asks once and passes what it heard (S-0104), so that its cards' dots and this notice
	// share one answer; without it the notice asks for itself.
	// The caret beside its bold title collapses it to that title (S-0150). The collapse is kept per
	// browser against what the notice says, so a reload keeps it and any change to it opens it again.
	import { api } from '$lib/api';
	import type { HostAgent } from '$lib/activity';
	import { localTime } from '$lib/localtime';

	let { refresh = 0, status: given }: { refresh?: number; status?: HostAgent | null } = $props();
	let asked = $state<HostAgent | null>(null);
	const status = $derived(given !== undefined ? given : asked);

	async function ask() {
		try {
			const r = await api('/api/host-agent');
			if (r.ok) asked = await r.json();
		} catch {
			// keep what we had
		}
	}
	$effect(() => {
		void refresh;
		if (given === undefined) void ask();
	});

	const st = $derived(status?.enabled ? status.state : undefined);

	const KEY = 'flaiover-host-agent-collapsed';
	function stored(): string | null {
		try {
			return localStorage.getItem(KEY);
		} catch {
			return null;
		}
	}
	function keep(v: string | null) {
		try {
			if (v === null) localStorage.removeItem(KEY);
			else localStorage.setItem(KEY, v);
		} catch {
			// a browser that stores nothing collapses for this page only
		}
	}
	// What the notice says; the 15 s re-ask returns an equal one, which leaves a collapse alone.
	const said = $derived(
		st ? JSON.stringify([st.command, st.running ?? null, st.last ?? null, st.waiting ?? null]) : ''
	);
	let collapsed = $state<string | null>(stored());
	const open = $derived(collapsed !== said);
	$effect(() => {
		if (said && collapsed !== null && collapsed !== said) {
			collapsed = null;
			keep(null);
		}
	});
	function toggle() {
		collapsed = open ? said : null;
		keep(collapsed);
	}
</script>

{#snippet title(text: string)}
	<button
		type="button"
		class="cursor-pointer font-semibold"
		aria-expanded={open}
		data-testid="host-agent-toggle"
		onclick={toggle}
		><span aria-hidden="true" class="inline-block w-3 text-xs text-muted">{open ? '▾' : '▸'}</span
		>{text}</button
	>
{/snippet}

{#if st && (st.running || st.last || st.waiting)}
	<div
		class="mb-3 rounded border border-line-strong bg-raised p-2 text-sm"
		role="status"
		data-testid="host-agent"
	>
		{#if st.running}
			<p>
				{@render title('Agent started')}
				{#if open}for {st.running.story} by
					<code class="rounded bg-surface px-1">{st.running.command}</code> as {st.running.agent},
					{localTime(st.running.started)}. It is running on the host.{/if}
			</p>
		{:else if st.last?.error}
			<p class="text-danger" data-testid="host-agent-failed">
				{@render title('No agent could be started')}
				{#if open}for {st.last.story}:
					<code class="rounded bg-surface px-1">{st.last.command}</code>
					{st.last.error} ({localTime(st.last.started)}). The operator sets the command on the host
					with <code class="rounded bg-surface px-1">flai serve agent set</code>.{/if}
			</p>
		{:else if st.last}
			<p>
				{@render title('Agent ended')}{#if open}: the one started for {st.last.story} by
					<code class="rounded bg-surface px-1">{st.last.command}</code>
					{localTime(st.last.started)} ended{st.last.ended
						? ` ${localTime(st.last.ended)}`
						: ''}{st.last.exit ? ` with exit code ${st.last.exit}` : ''}.{/if}
			</p>
		{:else}
			<p>
				{@render title('A ready story is waiting')}{#if open}:
					<span data-testid="host-agent-waiting">{st.waiting}</span>.{/if}
			</p>
		{/if}
		{#if open}
			{#if st.waiting && (st.running || st.last)}
				<p class="mt-1 text-xs" data-testid="host-agent-waiting">
					A ready story is waiting: {st.waiting}.
				</p>
			{/if}
			<p class="mt-1 text-xs text-muted">
				flai on the host starts each ready story's agent{#if st.command}, or <code
						>{st.command}</code
					> for a story that names no harness{/if}, while the in-progress limit leaves room.
				Stopping an agent is done on the host; the dashboard cannot.
			</p>
		{/if}
	</div>
{/if}
