<script lang="ts">
	// The count of things that need a human, beside the Inbox link (S-0042).
	import { inboxState } from '$lib/inbox.svelte';
</script>

{#if inboxState.data?.total}
	<span
		class="ml-1 rounded-full bg-primary px-1.5 text-xs text-on-primary"
		data-testid="inbox-count"
		aria-label="{inboxState.data.total} in the inbox">{inboxState.data.total}</span
	>
{/if}
