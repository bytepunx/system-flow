<script lang="ts">
	// /inbox: what needs a human (S-0042). The data is shared with the navigation badge.
	import { inboxState } from '$lib/inbox.svelte';
	import InboxView from '$lib/components/InboxView.svelte';
	import { api } from '$lib/api';
	import { onMount } from 'svelte';

	let writable = $state(false);
	onMount(() => {
		api('/api/board')
			.then(async (r) => (writable = r.ok ? (await r.json()).writable : false))
			.catch(() => (writable = false));
	});
</script>

<svelte:head><title>Inbox · flaiover</title></svelte:head>

<div class="mb-4 flex flex-wrap items-baseline gap-4">
	<h1 class="text-2xl font-semibold">Inbox</h1>
	{#if inboxState.permission !== 'unsupported'}
		<label class="flex items-center gap-2 text-xs text-muted">
			<input
				type="checkbox"
				checked={inboxState.notify}
				disabled={inboxState.permission === 'denied'}
				onchange={(e) => inboxState.setNotify(e.currentTarget.checked)}
			/>
			Desktop notification when something new arrives
			{#if inboxState.permission === 'denied'}(blocked in this browser's settings){/if}
		</label>
	{/if}
</div>
{#if inboxState.error}
	<p class="rounded border border-danger bg-danger-soft p-3 text-sm text-danger" role="alert">
		{inboxState.error}
	</p>
{:else if !inboxState.data}
	<p class="text-sm text-muted">Loading…</p>
{:else}
	<InboxView inbox={inboxState.data} {writable} />
{/if}
