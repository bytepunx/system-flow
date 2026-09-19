<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import NewItemForm from '$lib/components/NewItemForm.svelte';

	const initialType = $derived(page.url.searchParams.get('type') === 'epic' ? 'epic' : 'story');
</script>

<svelte:head><title>New · flaiover</title></svelte:head>
<h1 class="mb-1 text-2xl font-semibold">New epic or story</h1>
<p class="mb-4 text-sm text-muted">
	Write what it is for. The ID, the file, its front matter, the link in its epic, and the commit are
	made for you. Tasks are written by the agent that pulls the story.
</p>
<NewItemForm {initialType} oncreated={(id) => goto(resolve('/items/[id]', { id }))} />
