<script lang="ts">
	// /edit/<path>?line=&story=: the editor for one document (ADR-0023, S-0040), at a line and on a
	// story's branch when the link names them (S-0378).
	import { page } from '$app/state';
	import { beforeNavigate } from '$app/navigation';
	import DocEditor from '$lib/components/DocEditor.svelte';

	const path = $derived(page.params.path ?? '');
	const line = $derived(Number(page.url.searchParams.get('line')) || undefined);
	const story = $derived(page.url.searchParams.get('story') || undefined);
	let dirty = $state(false);

	// Unsaved text is only in this tab: ask before leaving it, inside the app and out of it.
	beforeNavigate((nav) => {
		if (dirty && !confirm('Leave without saving your changes?')) nav.cancel();
	});
</script>

<svelte:window
	onbeforeunload={(e) => {
		if (dirty) e.preventDefault();
	}}
/>
<svelte:head><title>Edit {path.split('/').pop()} · flaiover</title></svelte:head>

{#if path}
	<DocEditor {path} {line} {story} ondirty={(d) => (dirty = d)} />
{:else}
	<p class="text-sm text-muted">No document given.</p>
{/if}
