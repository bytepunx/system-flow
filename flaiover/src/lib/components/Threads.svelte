<script lang="ts">
	// Threads anchored to a document or item (ADR-0020): read from
	// wip/threads, written through flai. `on` is a repository path or item ID;
	// without it, every thread of the project, each linking to its anchor, and no
	// new thread, which needs an anchor (the threads page, S-0173).
	import { api } from '$lib/api';
	import { render, slug } from '$lib/markdown';
	import { resolve } from '$app/paths';
	import { onMount, tick } from 'svelte';
	import { follow } from '$lib/events';
	import type { StoryActivity } from '$lib/activity';

	type Source = { path: string; heading?: string };
	type Entry = {
		at: string;
		author: string;
		text: string;
		operator?: boolean;
		/** A recommendation the operator confirms to make it the answer, and what it cites (ADR-0090). */
		recommendation?: boolean;
		source?: Source | null;
	};
	type Thread = {
		id: string;
		title: string;
		anchor: { path: string; heading?: string; item?: string };
		status: 'open' | 'answered' | 'resolved';
		participants: string[];
		updated: string;
		entries: Entry[];
		/** The recommendation awaiting the operator's confirmation, if any (ADR-0090). */
		pending_recommendation?: Entry | null;
	};

	let {
		on,
		headings = [],
		writable = true,
		compose,
		select,
		agent
	}: {
		on?: string;
		headings?: string[];
		writable?: boolean;
		/** A heading to start a new thread on: opens the composer with it selected (the editor's "open a thread on this heading", S-0040). */
		compose?: string;
		/** A thread to open on, such as the one an inbox entry leads to (S-0155). */
		select?: string;
		/** What the story's agent is doing, on a story's page (S-0154). */
		agent?: StoryActivity;
	} = $props();

	let threads = $state<Thread[]>([]);
	let showResolved = $state(false);
	let notice = $state<string | null>(null);
	let composing = $state(false);
	let title = $state('');
	let text = $state('');
	let heading = $state('');
	let replies = $state<Record<string, string>>({});
	// One thread is shown at a time (S-0133). The one being read is held by ID so a reload
	// keeps it; when it drops out of the list, the thread now in its place is shown.
	let current = $state<string | null>(null);
	const index = $derived(
		Math.max(
			0,
			threads.findIndex((t) => t.id === current)
		)
	);
	const shown = $derived(threads[index]);
	// A thread longer than LATEST entries shows only its last ones until the reader asks for the
	// earlier ones, per thread and kept while paging (S-0153, TH-0035).
	const LATEST = 2;
	let earlier = $state<Record<string, boolean>>({});

	/**
	 * Whether the story's agent is at work on the operator's reply that ends this thread (S-0154):
	 * it runs, or it ended asking this thread and flai is starting it again now that it is answered.
	 */
	function working(t: Thread): boolean {
		if (!agent || t.status === 'resolved' || !t.entries.at(-1)?.operator) return false;
		return agent.state === 'working' || (agent.state === 'waiting' && agent.thread === t.id);
	}

	$effect(() => {
		if (compose && writable) {
			heading = compose;
			composing = true;
		}
	});
	/** Left and Right page while an arrow of a pager has focus (S-0153). */
	function keys(e: KeyboardEvent) {
		const step = { ArrowLeft: -1, ArrowRight: 1 }[e.key];
		if (!step) return;
		e.preventDefault();
		go(step);
	}

	async function load() {
		const query = [on && `on=${encodeURIComponent(on)}`, showResolved && 'all=1'].filter(Boolean);
		const r = await api(`/api/threads${query.length ? '?' + query.join('&') : ''}`);
		if (!r.ok) return;
		const was = index;
		threads = await r.json();
		if (!threads.some((t) => t.id === current))
			current = threads[Math.min(was, threads.length - 1)]?.id ?? null;
	}
	function go(step: number) {
		const t = threads[index + step];
		if (t) current = t.id;
	}
	// The thread a link names is opened, and the section brought into view, once it is in the
	// list; a reload after the reader pages away leaves them where they are (S-0155).
	let section = $state<HTMLElement>();
	let selected: string | undefined;
	$effect(() => {
		if (!select || select === selected || !threads.some((t) => t.id === select)) return;
		selected = current = select;
		void tick().then(() => section?.scrollIntoView?.({ block: 'start' }));
	});
	$effect(() => {
		void on;
		void showResolved;
		void load();
	});
	// A reply, a new thread, or a resolution by anyone shows without a reload (S-0154); a thread's
	// story is read from the items, so their changes count too (S-0161).
	onMount(() => follow(['thread', 'item'], () => void load()));

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
		return r.ok ? data : null;
	}

	async function open() {
		if (!on || !title.trim() || !text.trim()) return;
		const opened = await post('/api/threads', { on, heading: heading || undefined, title, text });
		if (opened) {
			if (typeof opened.id === 'string') current = opened.id;
			title = text = heading = '';
			composing = false;
		}
	}
	async function reply(id: string) {
		const t = replies[id]?.trim();
		if (!t) return;
		if (await post(`/api/threads/${id}/reply`, { text: t })) replies[id] = '';
	}
	// flai ends an entry that cites a source with a `Source: <path> § <heading>` paragraph; it is
	// shown as a link to the document and heading instead (ADR-0090).
	const body = (e: Entry) => (e.source ? e.text.replace(/\n\nSource: [^\n]*$/, '') : e.text);
	const sourceHref = (s: Source) =>
		resolve('/docs/[...path]', { path: s.path }) + (s.heading ? `#${slug(s.heading)}` : '');
	const badge: Record<Thread['status'], string> = {
		open: 'bg-warn-soft text-warn  ',
		answered: 'bg-info-soft text-info  ',
		resolved: 'bg-raised text-ink-soft  '
	};
</script>

<!-- One thread at a time (S-0133). The pager reads `← n of m →` and sits beside the heading,
     with a second one under the thread for a reader who has scrolled to its end (S-0153). -->
{#snippet pager(where: 'above' | 'below')}
	{#if threads.length > 1}
		<nav
			class="flex items-center gap-1 text-xs {where === 'below' ? 'mt-3' : ''}"
			aria-label="threads {where}"
			data-pager={where}
		>
			<button
				type="button"
				class="rounded border border-line-strong px-2 py-1 hover:bg-raised disabled:opacity-40"
				aria-label="previous thread"
				onkeydown={keys}
				disabled={index === 0}
				onclick={() => go(-1)}>←</button
			>
			<span
				class="min-w-[4.5em] text-center tabular-nums"
				aria-live={where === 'above' ? 'polite' : 'off'}>{index + 1} of {threads.length}</span
			>
			<button
				type="button"
				class="rounded border border-line-strong px-2 py-1 hover:bg-raised disabled:opacity-40"
				aria-label="next thread"
				onkeydown={keys}
				disabled={index === threads.length - 1}
				onclick={() => go(1)}>→</button
			>
		</nav>
	{/if}
{/snippet}

<section bind:this={section} class="mt-6 text-sm" data-threads={on ?? 'all'}>
	<div class="flex flex-wrap items-center gap-3">
		<h2 class="font-medium">Threads</h2>
		{@render pager('above')}
		<label class="text-xs text-muted"
			><input type="checkbox" bind:checked={showResolved} /> show resolved</label
		>
		{#if writable && on}
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
	{#if shown}
		{@const t = shown}
		{@const hidden = earlier[t.id] ? 0 : Math.max(0, t.entries.length - LATEST)}
		<article class="mt-3 rounded border border-line bg-surface p-3" data-thread={t.id}>
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
				{:else if !t.anchor.item && t.anchor.path !== on}
					<a
						class="text-xs underline"
						href={resolve('/docs/[...path]', { path: t.anchor.path })}
						data-anchor={t.anchor.path}>{t.anchor.path}</a
					>
				{/if}
			</header>
			<!-- A long thread shows its last entries, the earlier ones behind a toggle (S-0153). -->
			{#if t.entries.length > LATEST}
				<button
					type="button"
					class="mt-2 text-xs text-muted underline"
					data-earlier={t.id}
					aria-expanded={!hidden}
					onclick={() => (earlier[t.id] = !earlier[t.id])}
					>{hidden
						? `show ${hidden} earlier ${hidden === 1 ? 'entry' : 'entries'}`
						: 'hide earlier entries'}</button
				>
			{/if}
			<!-- The operator's entries sit on the right in the primary tint, agents' on the left in
			     the neutral one, so who said what reads at a glance (S-0127). -->
			<ol class="mt-2 space-y-2">
				{#each t.entries.slice(hidden) as e (e.at + e.author)}
					<li
						class="flex flex-col text-sm {e.operator ? 'items-end' : 'items-start'}"
						data-from={e.operator ? 'operator' : 'agent'}
					>
						<div class="text-xs text-muted">
							<span class="font-mono">{e.at}</span>
							{e.author}
							{#if e.recommendation}<span
									class="ml-1 rounded border border-info px-1.5 py-0.5 text-[10px] text-info uppercase"
									data-recommendation>recommendation</span
								>{/if}
						</div>
						<!-- Entries are markdown, as they are in the thread's file (S-0126). Authors write
						     repository paths, so relative links resolve from the root. -->
						<div
							class="prose prose-sm mt-0.5 max-w-[85%] rounded-lg border px-3 py-2 [&>:first-child]:mt-0 [&>:last-child]:mb-0 {e.operator
								? 'border-primary bg-primary-soft'
								: 'border-line bg-raised'}"
						>
							<!-- eslint-disable-next-line svelte/no-at-html-tags -- repository markdown, rendered client side as every document is -->
							{@html render(body(e), '')}
						</div>
						{#if e.source}
							<!-- eslint-disable svelte/no-navigation-without-resolve -- the path is resolve()d; the rule does not follow the heading added to it -->
							<p class="mt-0.5 text-xs text-muted" data-source>
								Source:
								<a class="underline" href={sourceHref(e.source)}
									>{e.source.path}{e.source.heading ? ` § ${e.source.heading}` : ''}</a
								>
							</p>
							<!-- eslint-enable svelte/no-navigation-without-resolve -->
						{/if}
					</li>
				{/each}
			</ol>
			{#if working(t)}
				<!-- On the agent's side, where its answer will appear; still for a reader who asks for less motion. -->
				<p
					class="mt-2 flex items-center gap-2 text-xs text-muted"
					role="status"
					data-testid="agent-working"
				>
					<span class="flex gap-1" aria-hidden="true">
						{#each [0, 150, 300] as delay (delay)}
							<span
								class="h-1.5 w-1.5 rounded-full bg-dot-working motion-safe:animate-bounce"
								style="animation-delay: {delay}ms"
							></span>
						{/each}
					</span>
					{agent?.run.agent || 'the agent'} is working on your reply
				</p>
			{/if}
			{#if t.pending_recommendation}
				<!-- A recommendation is no answer until the operator confirms it (ADR-0090). -->
				<div class="mt-2 flex items-center gap-2 text-xs" data-pending={t.id}>
					<span class="text-muted"
						>{t.pending_recommendation.author}'s recommendation awaits your confirmation</span
					>
					{#if writable}
						<button
							type="button"
							class="rounded bg-primary px-2 py-1 text-on-primary"
							data-confirm={t.id}
							onclick={() => post(`/api/threads/${t.id}/confirm`, {})}>Confirm</button
						>
					{/if}
				</div>
			{/if}
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
	{/if}
	{@render pager('below')}
</section>
