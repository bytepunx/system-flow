<script lang="ts">
	// What flai on the host did about starting an agent when a story became ready (S-0079): that it
	// started one, for which story, by which command and when; that it could not; or why a ready
	// story is waiting. It shows nothing while the operator has not enabled it. Starting, stopping,
	// and configuring are done on the host; nothing here can.
	// The board asks once and passes what it heard (S-0104), so that its cards' dots and this notice
	// share one answer; without it the notice asks for itself.
	import { api } from '$lib/api';
	import type { HostAgent } from '$lib/activity';

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

	const at = (s: string) => s.replace('T', ' ').replace(/:\d\dZ$/, ' UTC');
	const st = $derived(status?.enabled ? status.state : undefined);
</script>

{#if st && (st.running || st.last || st.waiting)}
	<div
		class="mb-3 rounded border border-line-strong bg-raised p-2 text-sm"
		role="status"
		data-testid="host-agent"
	>
		{#if st.running}
			<p>
				<span class="font-semibold">Agent started</span> for {st.running.story} by
				<code class="rounded bg-surface px-1">{st.running.command}</code> as {st.running.agent}, {at(
					st.running.started
				)}. It is running on the host.
			</p>
		{:else if st.last?.error}
			<p class="text-danger" data-testid="host-agent-failed">
				<span class="font-semibold">No agent could be started</span> for {st.last.story}:
				<code class="rounded bg-surface px-1">{st.last.command}</code>
				{st.last.error} ({at(st.last.started)}). The operator sets the command on the host with
				<code class="rounded bg-surface px-1">flai serve agent set</code>.
			</p>
		{:else if st.last}
			<p>
				The agent started for {st.last.story} by
				<code class="rounded bg-surface px-1">{st.last.command}</code>
				{at(st.last.started)} ended{st.last.ended ? ` ${at(st.last.ended)}` : ''}{st.last.exit
					? ` with exit code ${st.last.exit}`
					: ''}.
			</p>
		{/if}
		{#if st.waiting}
			<p class="mt-1 text-xs" data-testid="host-agent-waiting">
				A ready story is waiting: {st.waiting}.
			</p>
		{/if}
		<p class="mt-1 text-xs text-muted">
			flai on the host starts each ready story's agent{#if st.command}, or <code>{st.command}</code> for
				a story that names no harness{/if}, while the in-progress limit leaves room. Stopping an
			agent is done on the host; the dashboard cannot.
		</p>
	</div>
{/if}
