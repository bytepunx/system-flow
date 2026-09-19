<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import NewAdrForm from '$lib/components/NewAdrForm.svelte';
</script>

<svelte:head><title>New ADR · flaiover</title></svelte:head>
<h1 class="mb-1 text-2xl font-semibold">New architecture decision</h1>
<p class="mb-4 text-sm text-muted">
	One decision per record. Write the context, the decision, its consequences, and what lost. The
	number, the file, its front matter, the row in the index, and the commit are made for you.
</p>
<NewAdrForm oncreated={(path) => goto(resolve('/docs/[...path]', { path }))} />
