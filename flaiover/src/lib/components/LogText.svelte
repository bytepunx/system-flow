<script lang="ts">
	// Free text from flai, a status, a why, or an error, with each log path it names as a LogLink to
	// the log page and the rest as text (S-0785), so no message leaves the operator a path to open in
	// a shell. Which paths are log paths is splitLogPaths'.
	import LogLink from '$lib/components/LogLink.svelte';
	import { splitLogPaths } from '$lib/logs';

	let {
		text
	}: {
		/** The text as flai gives it. */
		text: string;
	} = $props();

	const parts = $derived(splitLogPaths(text));
</script>

<span data-testid="log-text">
	{#each parts as p, i (i)}
		{#if 'path' in p}
			<LogLink path={p.path} />
		{:else}
			{p.text}
		{/if}
	{/each}
</span>
