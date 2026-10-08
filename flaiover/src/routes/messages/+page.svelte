<script lang="ts">
	// /messages: every open conversation between the agents of two stories, and the closed ones on
	// request (S-0336). Read-only: the two stories' agents write them, and none counts in the inbox.
	import Messages from '$lib/components/Messages.svelte';
</script>

<svelte:head><title>Messages · flaiover</title></svelte:head>

<h1 class="text-2xl font-semibold">Messages</h1>
<p class="mt-1 text-sm text-muted">
	The conversations between the agents of two stories, about the paths both touch. They never reach
	the inbox; one the two stories cannot settle is escalated to a thread.
</p>
<Messages />
