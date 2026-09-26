<script lang="ts">
	// Threads anchored to a document or item (ADR-0020): read from
	// wip/threads, written through flai. `on` is a repository path or item ID.
	import { api } from '$lib/api';
	import { render } from '$lib/markdown';
	import { resolve } from '$app/paths';

	type Entry = { at: string; author: string; text: string };
	type Thread = {
		id: string;
		title: string;
		anchor: { path: string; heading?: string; item?: string };
		status: 'open' | 'answered' | 'resolved';
		participants: string[];
		updated: string;
		entries: Entry[];
	};

	let {
		on,
		headings = [],
		writable = true,
		compose
	}: {
		on: string;
		headings?: string[];
		writable?: boolean;
		/** A heading to start a new thread on: opens the composer with it selected (the editor's "open a thread on this heading", S-0040). */
		compose?: string;
	} = $props();

	let threads = $state<Thread[]>([]);
	let showResolved = $state(false);
	let notice = $state<string | null>(null);
	let composing = $state(false);
	let title = $state('');
	let text = $state('');
	let heading = $state('');
	let replies = $state<Record<string, string>>({});

	$effect(() => {
		if (compose && writable) {
			heading = compose;
			composing = true;
		}
	});

	async function load() {
		const r = await api(`/api/threads?on=${encodeURIComponent(on)}${showResolved ? '&all=1' : ''}`);
		if (r.ok) threads = await r.json();
	}
	$effect(() => {
		void on;
		void showResolved;
		void load();
	});

	async function post(path: string, body: Record<string, unknown>) {
		notice = null;
		const r = await api(path, {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify(body)
		});
		const data = await r.json().catch(() => ({}));
		if (!r.ok) notice = `refused: ${data.error ?? r.statusText}`;
		await load();
		return r.ok;
	}

	async function open() {
		if (!title.trim() || !text.trim()) return;
		if (await post('/api/threads', { on, heading: heading || undefined, title, text })) {
			title = text = heading = '';
			composing = false;
		}
	}
	async function reply(id: string) {
		const t = replies[id]?.trim();
		if (!t) return;
		if (await post(`/api/threads/${id}/reply`, { text: t })) replies[id] = '';
	}
	const badge: Record<Thread['status'], string> = {
		open: 'bg-warn-soft text-warn  ',
		answered: 'bg-info-soft text-info  ',
		resolved: 'bg-raised text-ink-soft  '
	};
</script>

<section class="mt-6 text-sm" data-threads={on}>
	<div class="flex items-center gap-3">
		<h2 class="font-medium">Threads</h2>
		<label class="text-xs text-muted"
			><input type="checkbox" bind:checked={showResolved} /> show resolved</label
		>
		{#if writable}
			<button
				type="button"
				class="ml-auto rounded border border-line-strong px-2 py-1 text-xs"
				onclick={() => (composing = !composing)}>{composing ? 'cancel' : 'new thread'}</button
			>
		{/if}
	</div>
	{#if notice}<p class="mt-2 text-xs text-danger">{notice}</p>{/if}
	{#if composing}
		<form
			class="mt-3 space-y-2 rounded border border-line bg-surface p-3"
			onsubmit={(e) => {
				e.preventDefault();
				void open();
			}}
		>
			<input
				class="w-full rounded border border-line-strong px-2 py-1 text-sm"
				placeholder="title"
				bind:value={title}
			/>
			{#if headings.length}
				<select
					class="w-full rounded border border-line-strong px-2 py-1 text-xs"
					bind:value={heading}
				>
					<option value="">whole document</option>
					{#each headings as h (h)}<option value={h}>{h}</option>{/each}
				</select>
			{/if}
			<textarea
				class="w-full rounded border border-line-strong px-2 py-1 text-sm"
				rows="3"
				placeholder="what do you want to ask or say?"
				bind:value={text}></textarea>
			<button type="submit" class="rounded bg-primary px-3 py-1 text-xs text-on-primary"
				>post</button
			>
		</form>
	{/if}
	{#if threads.length === 0}
		<p class="mt-2 text-xs text-muted">No threads here.</p>
	{/if}
	{#each threads as t (t.id)}
		<article class="mt-3 rounded border border-line bg-surface p-3">
			<header class="flex flex-wrap items-center gap-2">
				<span class="font-mono text-xs text-muted">{t.id}</span>
				<span class="font-medium">{t.title}</span>
				<span class="rounded px-1.5 py-0.5 text-[10px] uppercase {badge[t.status]}">{t.status}</span
				>
				{#if t.anchor.heading}<span class="text-xs text-muted">§ {t.anchor.heading}</span>{/if}
				{#if t.anchor.item && t.anchor.item !== on}
					<a class="text-xs underline" href={resolve('/items/[id]', { id: t.anchor.item })}
						>{t.anchor.item}</a
					>
				{/if}
			</header>
			<ol class="mt-2 space-y-2">
				{#each t.entries as e (e.at + e.author)}
					<li class="text-sm">
						<div class="text-xs text-muted">
							<span class="font-mono">{e.at}</span>
							{e.author}
						</div>
						<!-- Entries are markdown, as they are in the thread's file (S-0126). Authors write
						     repository paths, so relative links resolve from the root. -->
						<div class="prose prose-sm max-w-none">
							<!-- eslint-disable-next-line svelte/no-at-html-tags -- repository markdown, rendered client side as every document is -->
							{@html render(e.text, '')}
						</div>
					</li>
				{/each}
			</ol>
			{#if writable}
				<form
					class="mt-2 flex gap-2"
					onsubmit={(e) => {
						e.preventDefault();
						void reply(t.id);
					}}
				>
					<input
						class="min-w-0 flex-1 rounded border border-line-strong px-2 py-1 text-xs"
						placeholder="reply"
						bind:value={replies[t.id]}
					/>
					<button type="submit" class="rounded border border-line-strong px-2 py-1 text-xs"
						>reply</button
					>
					{#if t.status !== 'resolved'}
						<button
							type="button"
							class="rounded border border-line-strong px-2 py-1 text-xs"
							onclick={() => post(`/api/threads/${t.id}/resolve`, {})}>resolve</button
						>
					{/if}
				</form>
			{/if}
		</article>
	{/each}
</section>
