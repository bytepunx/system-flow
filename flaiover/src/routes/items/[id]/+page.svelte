<script lang="ts">
	import { themeState } from '$lib/theme.svelte';
	import { api } from '$lib/api';
	import Threads from '$lib/components/Threads.svelte';
	import AcceptConfirm from '$lib/components/AcceptConfirm.svelte';
	import { page } from '$app/state';
	import { resolve } from '$app/paths';
	import { tick } from 'svelte';
	import { render, enhance } from '$lib/markdown';

	type Transition = { to: string; at: string; by: string };
	type Block = { from: string; until?: string; reason: string };
	type Item = {
		id: string;
		type: string;
		nature: string;
		title: string;
		status: string;
		parent?: string;
		owner?: string;
		created: string;
		updated: string;
		transitions: Transition[];
		blocked?: Block[];
		estimate?: string;
		stream?: string;
		tags?: string[];
		path: string;
		archived: boolean;
		body: string;
	};

	let item = $state<Item | null>(null);
	let children = $state<Item[]>([]);
	let html = $state('');
	let error = $state<string | null>(null);
	let notice = $state<string | null>(null);
	let writable = $state(false);
	let content: HTMLElement | undefined = $state();
	let entry = $state('');

	const id = $derived(page.params.id ?? '');
	const allowed: Record<string, string[]> = {
		backlog: ['ready', 'cancelled'],
		ready: ['in-progress', 'cancelled'],
		'in-progress': ['review', 'cancelled'],
		review: ['done', 'in-progress']
	};
	const moves = $derived(
		item
			? [
					...(allowed[item.status] ?? []),
					...(item.type === 'task' && item.status === 'in-progress' ? ['done'] : [])
				]
			: []
	);
	const blocked = $derived((item?.blocked ?? []).some((b) => !b.until));
	const narrative = $derived(
		item?.type === 'story'
			? `wip/${item.archived ? 'archive/agents' : 'agents'}/${item.id}.md`
			: item?.stream
				? `wip/agents/${item.stream}.md`
				: null
	);

	async function load() {
		const r = await api(`/api/items/${id}`);
		if (!r.ok) {
			error = (await r.json()).error ?? r.statusText;
			return;
		}
		const data = await r.json();
		item = data.item;
		children = data.children;
		html = render(item!.body, item!.path);
		writable = (await (await api('/api/board')).json()).writable;
		await tick();
		if (content) await enhance(content, themeState.dark);
	}
	$effect(() => {
		void id;
		load();
	});

	async function post(path: string, body: Record<string, unknown>) {
		notice = null;
		const r = await api(path, {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify(body)
		});
		const data = await r.json();
		if (!r.ok) notice = `refused: ${data.error}`;
		else if (data.push_error)
			notice = `accepted locally but not pushed (${data.push_error}); push the commit and tags from a shell`;
		else if (data.tags?.length) notice = `accepted: released ${data.tags.join(', ')}`;
		else notice = `done${data.warnings?.length ? ': ' + data.warnings.join(' ') : ''}`;
		await load();
	}
	let accepting = $state(false);
	function move(to: string) {
		// A story going to done is an acceptance: confirm with the plan first (S-0046).
		if (to === 'done' && item?.type === 'story' && item.status === 'review' && !accepting) {
			accepting = true;
			return;
		}
		let reason: string | undefined;
		if (to === 'cancelled' || (item?.status === 'review' && to === 'in-progress')) {
			reason = prompt(`Reason for ${to}:`) ?? undefined;
			if (!reason) return;
		}
		post(`/api/items/${id}/move`, { to, reason });
	}
	function block() {
		const reason = prompt('Why is it blocked?');
		if (reason) post(`/api/items/${id}/block`, { reason });
	}
</script>

<svelte:head><title>{id} · flaiover</title></svelte:head>

{#if accepting && item}
	<AcceptConfirm
		id={item.id}
		oncancel={() => (accepting = false)}
		onconfirm={async (include) => {
			await post(`/api/items/${item!.id}/move`, { to: 'done', include_uncommitted: include });
			accepting = false;
		}}
	/>
{/if}

{#if error}
	<p class="rounded border border-danger bg-danger-soft p-3 text-sm text-danger">{error}</p>
{:else if item}
	<div class="grid grid-cols-1 gap-6 lg:grid-cols-[minmax(0,1fr)_320px]">
		<div class="min-w-0">
			<h1 class="text-2xl font-semibold"><span class="font-mono">{item.id}</span> {item.title}</h1>
			<p class="mt-1 text-sm text-muted">
				{item.type} · {item.nature} · <span class="font-medium">{item.status}</span>
				{#if blocked}<span class="ml-1 font-semibold text-danger">BLOCKED</span>{/if}
				{#if item.archived}· archived{/if}
				{#if item.parent}· parent <a
						class="underline"
						href={resolve('/items/[id]', { id: item.parent })}>{item.parent}</a
					>{/if}
			</p>
			{#if notice}<p class="mt-2 rounded border border-line-strong bg-raised p-2 text-sm">
					{notice}
				</p>{/if}
			{#if item.type === 'story' && item.status === 'review'}
				<p class="mt-3 text-sm">
					<a
						class="rounded border border-line-strong px-2 py-1 hover:bg-raised"
						href={resolve('/review/[id]', { id: item.id })}>Review this story</a
					>
					<span class="ml-2 text-xs text-muted"
						>criteria, narrative, threads, the branch's diff, and what accepting does</span
					>
				</p>
			{/if}
			{#if writable && !item.archived}
				<div class="mt-3 flex flex-wrap gap-2">
					{#each moves as to (to)}
						<button
							type="button"
							class="rounded border border-line-strong px-2 py-1 text-xs hover:bg-raised"
							onclick={() => move(to)}>→ {to}</button
						>
					{/each}
					{#if blocked}
						<button
							type="button"
							class="rounded border border-line-strong px-2 py-1 text-xs"
							onclick={() => post(`/api/items/${id}/unblock`, {})}>unblock</button
						>
					{:else if item.status !== 'done' && item.status !== 'cancelled'}
						<button
							type="button"
							class="rounded border border-line-strong px-2 py-1 text-xs"
							onclick={block}>block…</button
						>
					{/if}
				</div>
			{/if}
			<article bind:this={content} class="prose mt-4 max-w-none">
				<!-- eslint-disable-next-line svelte/no-at-html-tags -- markdown from the mounted repository, rendered client side -->
				{@html html}
			</article>
			<Threads on={item.id} {writable} />
		</div>
		<aside class="space-y-4 text-sm">
			<section class="rounded border border-line bg-surface p-3">
				<h2 class="mb-2 font-medium">History</h2>
				<ol class="space-y-1 text-xs">
					<li><span class="font-mono text-muted">{item.created}</span> created</li>
					{#each item.transitions as t (t.at + t.to)}
						<li>
							<span class="font-mono text-muted">{t.at}</span>
							{t.to} <span class="text-muted">by {t.by}</span>
						</li>
					{/each}
				</ol>
				{#if item.blocked?.length}
					<h3 class="mt-3 mb-1 font-medium">Blocked</h3>
					<ul class="space-y-1 text-xs">
						{#each item.blocked as b (b.from)}
							<li>
								<span class="font-mono text-muted">{b.from}</span> → {b.until ?? 'open'}: {b.reason}
							</li>
						{/each}
					</ul>
				{/if}
			</section>
			{#if children.length}
				<section class="rounded border border-line bg-surface p-3">
					<h2 class="mb-2 font-medium">Children</h2>
					<ul class="space-y-1 text-xs">
						{#each children as c (c.id)}
							<li>
								<a class="font-mono underline" href={resolve('/items/[id]', { id: c.id })}>{c.id}</a
								> <span class="text-muted">{c.status}</span>
								{c.title}
							</li>
						{/each}
					</ul>
				</section>
			{/if}
			<section class="rounded border border-line bg-surface p-3">
				<h2 class="mb-2 font-medium">Files</h2>
				<p class="text-xs">
					<a class="underline" href={resolve('/docs/[...path]', { path: item.path })}>{item.path}</a
					>
					{#if writable && !item.archived}
						· <a class="underline" href={resolve('/edit/[...path]', { path: item.path })}
							>edit the body</a
						>
					{/if}
				</p>
				{#if narrative}
					<p class="mt-1 text-xs">
						narrative: <a class="underline" href={resolve('/docs/[...path]', { path: narrative })}
							>{narrative}</a
						>
					</p>
					{#if writable && item.type === 'story' && !item.archived}
						<form
							class="mt-2 flex gap-2"
							onsubmit={(e) => {
								e.preventDefault();
								if (entry.trim()) {
									post(`/api/streams/${item!.id}/log`, { entry });
									entry = '';
								}
							}}
						>
							<input
								class="min-w-0 flex-1 rounded border border-line-strong px-2 py-1 text-xs"
								placeholder="log entry"
								bind:value={entry}
							/>
							<button type="submit" class="rounded border border-line-strong px-2 py-1 text-xs"
								>log</button
							>
						</form>
					{/if}
				{/if}
			</section>
			{#if item.tags?.length || item.owner || item.estimate}
				<section class="rounded border border-line bg-surface p-3 text-xs">
					{#if item.owner}<div>owner: {item.owner}</div>{/if}
					{#if item.estimate}<div>estimate: {item.estimate}</div>{/if}
					{#if item.tags?.length}<div>tags: {item.tags.join(', ')}</div>{/if}
				</section>
			{/if}
		</aside>
	</div>
{:else}
	<p class="text-sm text-muted">Loading…</p>
{/if}
