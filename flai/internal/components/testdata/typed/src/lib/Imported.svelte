<script lang="ts">
	// The whole props type is imported, so each prop's kind is unknown.
	import type { Props } from '$lib/types';

	let { story, count }: Props = $props();
</script>

<p>{story.id} {count}</p>
