<script lang="ts">
	// /edit/<path>: the editor for one document (ADR-0023, S-0040).
	import { page } from '$app/state';
	import { beforeNavigate } from '$app/navigation';
	import DocEditor from '$lib/components/DocEditor.svelte';

	const path = $derived(page.params.path ?? '');
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
	<DocEditor {path} ondirty={(d) => (dirty = d)} />
{:else}
	<p class="text-sm text-muted">No document given.</p>
{/if}
