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
	import PageSidebar from '$lib/components/PageSidebar.svelte';
	import { hashId, type SidebarEntry } from '$lib/pagesidebar.svelte';

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
			const r = await api('/api/items?archived=false');
			if (r.ok) {
				workers = (await r.json()) as Worker[];
			}
		} catch {
			workers = [];
		}
	}

	// S-0461: the open document's headings, as the sidebar's scroll entries. Their ids are the ones
	// the renderer gave the headings (markdown-it-anchor, unique within the document), read from the
	// rendered HTML, so an entry always names its heading; their labels are the headings' text.
	type Heading = { id: string; label: string; level: number };
	function headingsIn(rendered: string): Heading[] {
		if (!rendered || typeof DOMParser === 'undefined') return [];
		const body = new DOMParser().parseFromString(rendered, 'text/html').body;
		return [...body.querySelectorAll<HTMLElement>('h1, h2, h3, h4, h5, h6')]
			.filter((h) => h.id)
			.map((h) => ({ id: h.id, label: h.textContent!.trim(), level: Number(h.tagName[1]) }));
	}
	// Each heading sits under the nearest heading above it of a higher level.
	function nest(hs: Heading[]): SidebarEntry[] {
		const top: SidebarEntry[] = [];
		const open: { level: number; entry: SidebarEntry }[] = [];
		for (const h of hs) {
			const entry: SidebarEntry = { id: h.id, label: h.label, mode: 'scroll' };
			while (open.length && open[open.length - 1].level >= h.level) open.pop();
			const parent = open[open.length - 1]?.entry;
			if (parent) (parent.children ??= []).push(entry);
			else top.push(entry);
			open.push({ level: h.level, entry });
		}
		return top;
	}
	const headings = $derived(headingsIn(html));
	// Under a label, whose id has a space and so is never a heading's slug.
	const sections = $derived<SidebarEntry[]>(
		headings.length ? [{ id: 'On this page', label: 'On this page', children: nest(headings) }] : []
	);

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
				// A link to one of its headings, from another document or shared, lands on it once drawn.
				const at = hashId(location.hash);
				if (at && headings.some((h) => h.id === at))
					document.getElementById(at)?.scrollIntoView?.({ block: 'start' });
			})
			.catch((e) => {
				error = e instanceof Error ? e.message : String(e);
				doc = null;
				html = '';
			});
	});
</script>

<svelte:head><title>{current ? current.split('/').pop() : 'Docs'} · flaiover</title></svelte:head>

<!-- S-0461: the tree of documents heads the sidebar, and the open document's headings follow it as
     scroll entries; the document, its front matter, and its threads are the content. -->
<PageSidebar label="Documents" entries={sections}>
	{#snippet header()}
		<div class="max-h-96 overflow-y-auto rounded border border-line bg-surface p-3" data-doc-tree>
			{#if tree.length}
				<DocTree nodes={tree} {current} />
			{:else}
				<p class="text-sm text-muted">Loading tree…</p>
			{/if}
		</div>
	{/snippet}
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
			<article bind:this={content} class="prose max-w-none overflow-x-auto">
				<!-- eslint-disable-next-line svelte/no-at-html-tags -- markdown from the mounted repository, rendered client side -->
				{@html html}
			</article>
			<Threads on={current} headings={doc ? headingsOf(doc.body) : []} />
		{/if}
	</section>
</PageSidebar>
