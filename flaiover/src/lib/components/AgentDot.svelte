<script lang="ts">
	// A story's agent at a glance (S-0104): green while it works, yellow while it waits for the
	// designer, red when it failed. An agent that finished its story shows nothing.
	import { activityLine, dotClass, type StoryActivity } from '$lib/activity';

	let { activity }: { activity?: StoryActivity } = $props();
	const colour = $derived(activity ? dotClass(activity.state) : null);
</script>

{#if activity && colour}
	<span
		class="inline-block h-2.5 w-2.5 shrink-0 rounded-full {colour} {activity.state === 'working'
			? 'animate-pulse'
			: ''}"
		role="img"
		aria-label={activityLine(activity)}
		title={activityLine(activity)}
		data-testid="agent-dot"
		data-state={activity.state}
	></span>
{/if}
