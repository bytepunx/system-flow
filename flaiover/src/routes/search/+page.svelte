<script lang="ts">
	import { api } from '$lib/api';
	import { resolve } from '$app/paths';

	type Hit = {
		path: string;
		kind: 'item' | 'doc';
		itemId?: string;
		title: string;
		scope: string;
		status?: string;
		type?: string;
		snippet: string;
		route: string;
	};
	let q = $state('');
	let docs = $state(false);
	let hits = $state<Hit[]>([]);
	let indexed = $state(0);
	let timer: ReturnType<typeof setTimeout> | undefined;

	function run() {
		clearTimeout(timer);
		timer = setTimeout(async () => {
			if (!q.trim()) {
				hits = [];
				return;
			}
			const r = await api(`/api/search?q=${encodeURIComponent(q)}&docs=${docs}`);
			const body = await r.json();
			hits = body.hits;
			indexed = body.indexed;
		}, 150);
	}
	const href = (h: Hit) =>
		h.kind === 'item'
			? resolve('/docs/[...path]', { path: h.path })
			: resolve('/docs/[...path]', { path: h.path });
</script>

<svelte:head><title>Search · flaiover</title></svelte:head>
<h1 class="text-2xl font-semibold">Search</h1>
<form
	class="mt-3 flex flex-wrap items-center gap-3"
	onsubmit={(e) => {
		e.preventDefault();
		run();
	}}
>
	<input
		class="w-full max-w-lg rounded border border-zinc-300 bg-white px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
		placeholder="item ID, title, or words in the body"
		bind:value={q}
		oninput={run}
	/>
	<label class="flex items-center gap-2 text-sm"
		><input type="checkbox" bind:checked={docs} onchange={run} /> include docs/</label
	>
	{#if indexed}<span class="text-xs text-zinc-500">{indexed} files indexed</span>{/if}
</form>
<ul class="mt-4 space-y-3">
	{#each hits as h (h.path)}
		<li class="rounded border border-zinc-200 bg-white p-3 dark:border-zinc-800 dark:bg-zinc-900">
			<a class="font-medium underline" href={href(h)}>{h.itemId ? `${h.itemId} ` : ''}{h.title}</a>
			<span class="ml-2 text-xs text-zinc-500"
				>{h.scope}{h.type ? ` · ${h.type}` : ''}{h.status ? ` · ${h.status}` : ''}</span
			>
			<p class="mt-1 text-sm text-zinc-600 dark:text-zinc-400">{h.snippet}</p>
			<p class="mt-1 font-mono text-xs text-zinc-400">{h.path}</p>
		</li>
	{:else}
		{#if q.trim()}<li class="text-sm text-zinc-500">No results.</li>{/if}
	{/each}
</ul>
