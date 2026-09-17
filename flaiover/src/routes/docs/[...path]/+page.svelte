<script lang="ts">
	import { api } from '$lib/api';
	import { page } from '$app/state';
	import { onMount, tick } from 'svelte';
	import DocTree from '$lib/components/DocTree.svelte';
	import { render, enhance } from '$lib/markdown';

	type Node = {
		name: string;
		path: string;
		kind: 'dir' | 'file';
		title?: string;
		children?: Node[];
	};
	type Doc = { path: string; frontMatter: Record<string, unknown> | null; body: string };

	let tree = $state<Node[]>([]);
	let doc = $state<Doc | null>(null);
	let html = $state('');
	let error = $state<string | null>(null);
	let content: HTMLElement | undefined = $state();

	const current = $derived(page.params.path ?? '');

	onMount(async () => {
		const r = await api('/api/docs/tree');
		tree = await r.json();
	});

	$effect(() => {
		const path = current;
		if (!path) {
			doc = null;
			html = '';
			return;
		}
		api(`/api/docs/file?path=${encodeURIComponent(path)}`)
			.then(async (r) => {
				if (!r.ok) throw new Error((await r.json()).error ?? r.statusText);
				doc = await r.json();
				error = null;
				html = render(doc!.body, path);
				await tick();
				if (content) await enhance(content, matchMedia('(prefers-color-scheme: dark)').matches);
			})
			.catch((e) => {
				error = e instanceof Error ? e.message : String(e);
				doc = null;
				html = '';
			});
	});
</script>

<svelte:head><title>{current ? current.split('/').pop() : 'Docs'} · flaiover</title></svelte:head>

<div class="grid grid-cols-1 gap-6 md:grid-cols-[260px_minmax(0,1fr)]">
	<aside
		class="max-h-[80vh] overflow-auto rounded border border-zinc-200 bg-white p-3 dark:border-zinc-800 dark:bg-zinc-900"
	>
		{#if tree.length}
			<DocTree nodes={tree} {current} />
		{:else}
			<p class="text-sm text-zinc-500">Loading tree…</p>
		{/if}
	</aside>
	<section class="min-w-0">
		{#if error}
			<p class="rounded border border-red-300 bg-red-50 p-3 text-sm text-red-800">{error}</p>
		{:else if !current}
			<p class="text-sm text-zinc-500">
				Pick a document from the tree. Design, docs, and wip are all here.
			</p>
		{:else}
			{#if doc?.frontMatter}
				<details
					class="mb-4 rounded border border-zinc-200 bg-white p-3 text-xs dark:border-zinc-800 dark:bg-zinc-900"
				>
					<summary class="cursor-pointer font-medium">Front matter · {current}</summary>
					<dl class="mt-2 grid grid-cols-[max-content_1fr] gap-x-4 gap-y-1">
						{#each Object.entries(doc.frontMatter) as [k, v] (k)}
							<dt class="text-zinc-500">{k}</dt>
							<dd class="font-mono break-all">
								{typeof v === 'object' ? JSON.stringify(v) : String(v)}
							</dd>
						{/each}
					</dl>
				</details>
			{/if}
			<article bind:this={content} class="prose max-w-none prose-zinc dark:prose-invert">
				<!-- eslint-disable-next-line svelte/no-at-html-tags -- markdown from the mounted repository, rendered client side -->
				{@html html}
			</article>
		{/if}
	</section>
</div>
