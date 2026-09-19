<script lang="ts">
	import { themeState } from '$lib/theme.svelte';
	import Threads from '$lib/components/Threads.svelte';
	import { resolve } from '$app/paths';
	import { api } from '$lib/api';
	import { page } from '$app/state';
	import { onMount, tick } from 'svelte';
	import DocTree from '$lib/components/DocTree.svelte';
	import { render, enhance } from '$lib/markdown';
	import { headingsOf } from '$lib/edit';
	import { touching, type Worker } from '$lib/touches';

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
	let workers = $state<Worker[]>([]);
	let writable = $state(false);
	// Items in progress or review whose touches cover the open document (ADR-0019);
	// a story in review still owns its branch until it is accepted.
	const touchingNow = $derived(touching(current, workers));
	async function loadWorkers() {
		try {
			const r = await api('/api/items');
			if (r.ok) {
				workers = (await r.json()) as Worker[];
			}
		} catch {
			workers = [];
		}
	}

	onMount(async () => {
		void loadWorkers();
		api('/api/board')
			.then(async (r) => (writable = r.ok ? (await r.json()).writable : false))
			.catch(() => (writable = false));
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
				if (content) await enhance(content, themeState.dark);
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
	<aside class="max-h-[80vh] overflow-auto rounded border border-line bg-surface p-3">
		{#if tree.length}
			<DocTree nodes={tree} {current} />
		{:else}
			<p class="text-sm text-muted">Loading tree…</p>
		{/if}
	</aside>
	<section class="min-w-0">
		{#if error}
			<p class="rounded border border-danger bg-danger-soft p-3 text-sm text-danger">{error}</p>
		{:else if !current}
			<p class="text-sm text-muted">
				Pick a document from the tree. Design, docs, and wip are all here.
			</p>
		{:else}
			{#if touchingNow.length}
				<p class="mb-3 rounded border border-warn bg-warn-soft px-3 py-2 text-xs text-warn">
					Being worked on by
					{#each touchingNow as w, i (w.id)}{i ? ', ' : ' '}<a
							class="font-medium underline"
							href={resolve('/items/[id]', { id: w.id })}>{w.id}</a
						>
						{w.title}{/each}. Edits here may collide with that story's branch.
				</p>
			{/if}
			{#if writable && doc}
				<p class="mb-3 text-right text-xs">
					<a
						class="rounded border border-line px-2 py-1 hover:border-line-strong"
						href={resolve('/edit/[...path]', { path: current })}>Edit</a
					>
				</p>
			{/if}
			{#if doc?.frontMatter}
				<details class="mb-4 rounded border border-line bg-surface p-3 text-xs">
					<summary class="cursor-pointer font-medium">Front matter · {current}</summary>
					<dl class="mt-2 grid grid-cols-[max-content_1fr] gap-x-4 gap-y-1">
						{#each Object.entries(doc.frontMatter) as [k, v] (k)}
							<dt class="text-muted">{k}</dt>
							<dd class="font-mono break-all">
								{typeof v === 'object' ? JSON.stringify(v) : String(v)}
							</dd>
						{/each}
					</dl>
				</details>
			{/if}
			<article bind:this={content} class="prose max-w-none">
				<!-- eslint-disable-next-line svelte/no-at-html-tags -- markdown from the mounted repository, rendered client side -->
				{@html html}
			</article>
			<Threads on={current} headings={doc ? headingsOf(doc.body) : []} />
		{/if}
	</section>
</div>
