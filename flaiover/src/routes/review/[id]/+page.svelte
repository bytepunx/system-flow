<script lang="ts">
	// /review/<id>: review one story and accept it or send it back (S-0041).
	import { page } from '$app/state';
	import Review from '$lib/components/Review.svelte';

	const id = $derived(page.params.id ?? '');
</script>

<svelte:head><title>Review {id} · flaiover</title></svelte:head>

{#if id}
	<Review {id} />
{/if}
