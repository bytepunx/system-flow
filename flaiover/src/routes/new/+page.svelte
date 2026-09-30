<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import NewItemForm from '$lib/components/NewItemForm.svelte';
	import { startLane } from '$lib/lanes';

	const initialType = $derived(page.url.searchParams.get('type') === 'epic' ? 'epic' : 'story');
	// The board's lane menu opens this page with the lane it was opened on (S-0167).
	const initialLane = $derived(startLane(page.url.searchParams.get('status')));
	// An item page's New link opens this page with the epic of the story it was on (S-0171).
	const initialParent = $derived(page.url.searchParams.get('parent') ?? '');
</script>

<svelte:head><title>New · flaiover</title></svelte:head>
<h1 class="mb-1 text-2xl font-semibold">New epic or story</h1>
<p class="mb-4 text-sm text-muted">
	Write what it is for. The ID, the file, its front matter, the link in its epic, and the commit are
	made for you. Tasks are written by the agent that pulls the story.
</p>
<NewItemForm
	{initialType}
	{initialLane}
	{initialParent}
	oncreated={(id) => goto(resolve('/items/[id]', { id }))}
/>
