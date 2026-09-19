<script lang="ts">
	import { api } from '$lib/api';
	import { resolve } from '$app/paths';
	import KindChips from '$lib/components/KindChips.svelte';
	import { stripeClass, tintFor } from '$lib/cardcolour';

	type Hit = {
		path: string;
		kind: 'item' | 'doc';
		itemId?: string;
		title: string;
		scope: string;
		status?: string;
		type?: string;
		nature?: string;
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
		class="w-full max-w-lg rounded border border-line-strong bg-surface px-3 py-2 text-sm"
		placeholder="item ID, title, or words in the body"
		bind:value={q}
		oninput={run}
	/>
	<label class="flex items-center gap-2 text-sm"
		><input type="checkbox" bind:checked={docs} onchange={run} /> include docs/</label
	>
	{#if indexed}<span class="text-xs text-muted">{indexed} files indexed</span>{/if}
</form>
<ul class="mt-4 space-y-3">
	{#each hits as h (h.path)}
		<!-- An item hit wears the board card's coding (S-0055); a document stays plain. -->
		<li
			class="rounded border border-line p-3 {h.kind === 'item' && h.nature
				? `${tintFor(h.nature)} ${stripeClass(h.type ?? '')}`
				: 'bg-surface'}"
			data-kind={h.kind}
		>
			<a class="font-medium underline" href={href(h)}>{h.itemId ? `${h.itemId} ` : ''}{h.title}</a>
			<span class="ml-2 text-xs text-muted"
				>{h.scope}{#if h.type || h.nature}<span class="mx-1">·</span><KindChips
						type={h.type}
						nature={h.kind === 'item' ? h.nature : undefined}
					/>{/if}{h.status ? ` · ${h.status}` : ''}</span
			>
			<p class="mt-1 text-sm text-ink-soft">{h.snippet}</p>
			<p class="mt-1 font-mono text-xs text-muted">{h.path}</p>
		</li>
	{:else}
		{#if q.trim()}<li class="text-sm text-muted">No results.</li>{/if}
	{/each}
</ul>
