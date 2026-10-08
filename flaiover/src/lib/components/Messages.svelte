<script lang="ts">
	// Conversations between the agents of two stories (S-0336, ADR-0120): read from wip/messages
	// through /api/messages, open before closed, and read-only, since only the two stories' agents
	// write them. `story` narrows them to that story's and names the other side of each; without it,
	// every conversation of the project, each with both its stories (the messages page).
	import { api } from '$lib/api';
	import { render } from '$lib/markdown';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import { follow } from '$lib/events';
	import { localTime } from '$lib/localtime';
	import { age } from '$lib/age';

	type Entry = { at: string; author: string; story: string; text: string };
	type Conversation = {
		id: string;
		title: string;
		from: string;
		to: string;
		about: string[];
		status: 'open' | 'closed';
		closed: boolean;
		closed_reason: string;
		/** The story whose reply it waits for while open, '' once closed. */
		awaiting: string;
		participants: string[];
		created: string;
		updated: string;
		path: string;
		entries: Entry[];
	};

	let {
		story
	}: {
		/** A story's ID: only its conversations, each labelled by the other story (a story's page). */
		story?: string;
	} = $props();

	let conversations = $state<Conversation[]>([]);
	let showClosed = $state(false);
	let notice = $state<string | null>(null);
	let loaded = $state(false);
	let now = $state(Date.now());

	async function load() {
		const query = [story && `story=${encodeURIComponent(story)}`, showClosed && 'all=1'].filter(
			Boolean
		);
		const r = await api(`/api/messages${query.length ? '?' + query.join('&') : ''}`);
		if (!r.ok) {
			const data = await r.json().catch(() => ({}));
			notice = `could not read the conversations: ${data.error ?? r.statusText}`;
			return;
		}
		notice = null;
		conversations = await r.json();
		now = Date.now();
		loaded = true;
	}
	$effect(() => {
		void story;
		void showClosed;
		void load();
	});
	// A conversation's file is under the wip folder, so a document to $lib/changes; its stories'
	// states, which close it, are items.
	onMount(() => follow(['document', 'item'], () => void load()));

	/**
	 * The thread an escalation opened, as flai's escalatedText words the entry it adds, and the
	 * escalating story, whose thread it is; the last such entry wins.
	 */
	function escalation(c: Conversation): { thread: string; story: string } | null {
		for (const e of [...c.entries].reverse()) {
			const m = /^Asked the operator on (TH-\d+)/.exec(e.text);
			if (m) return { thread: m[1], story: e.story };
		}
		return null;
	}
	const since = (at: string) => age(Math.max(0, (now - Date.parse(at)) / 1000));
	const item = (id: string) => resolve('/items/[id]', { id });
	// A conversation reads as closed when either story is done, not only when its file says so.
	const badge = (c: Conversation) =>
		c.closed ? 'bg-raised text-ink-soft' : 'bg-warn-soft text-warn';
</script>

{#snippet storyLink(id: string)}
	<a class="font-mono text-xs underline" href={item(id)} data-story={id}>{id}</a>
{/snippet}

<section class="mt-6 text-sm" data-messages={story ?? 'all'}>
	<div class="flex flex-wrap items-center gap-3">
		<h2 class="font-medium">Messages</h2>
		<label class="text-xs text-muted"
			><input type="checkbox" bind:checked={showClosed} /> show closed</label
		>
	</div>
	{#if notice}<p class="mt-2 text-xs text-danger">{notice}</p>{/if}
	{#if loaded && conversations.length === 0}
		<p class="mt-2 text-xs text-muted" data-empty>
			{showClosed ? 'No conversations here.' : 'No open conversations here.'}
		</p>
	{/if}
	{#if conversations.length}
		<ul class="mt-3 space-y-3">
			{#each conversations as c (c.id)}
				{@const raised = escalation(c)}
				<li class="rounded border border-line bg-surface p-3" data-conversation={c.id}>
					<header class="flex flex-wrap items-center gap-2">
						<span class="font-mono text-xs text-muted">{c.id}</span>
						<span class="font-medium">{c.title}</span>
						<span class="rounded px-1.5 py-0.5 text-[10px] uppercase {badge(c)}"
							>{c.closed ? 'closed' : 'open'}</span
						>
						{#if raised}
							<!-- eslint-disable svelte/no-navigation-without-resolve -- the item is resolve()d; the rule does not follow the thread added to it -->
							<a
								class="rounded bg-danger-soft px-1.5 py-0.5 text-[10px] text-danger uppercase underline"
								href={item(raised.story) + `?thread=${encodeURIComponent(raised.thread)}`}
								data-escalated={raised.thread}>escalated on {raised.thread}</a
							>
							<!-- eslint-enable svelte/no-navigation-without-resolve -->
						{/if}
					</header>
					<p class="mt-1 flex flex-wrap items-center gap-x-2 text-xs text-muted">
						{#if story === c.from}
							<span data-with>to {@render storyLink(c.to)}</span>
						{:else if story === c.to}
							<span data-with>from {@render storyLink(c.from)}</span>
						{:else}
							<span data-between>{@render storyLink(c.from)} → {@render storyLink(c.to)}</span>
						{/if}
						{#if c.closed}
							<span data-closed>closed{c.closed_reason ? `: ${c.closed_reason}` : ''}</span>
						{:else if c.awaiting === story}
							<span data-awaiting={c.awaiting}>awaits this story</span>
						{:else if c.awaiting}
							<span data-awaiting={c.awaiting}>awaits {@render storyLink(c.awaiting)}</span>
						{/if}
						<span title={localTime(c.updated)} data-age>updated {since(c.updated)} ago</span>
					</p>
					{#if c.about.length}
						<ul class="mt-1 flex flex-wrap gap-1 text-xs" aria-label="about">
							{#each c.about as p (p)}<li><code data-about>{p}</code></li>{/each}
						</ul>
					{/if}
					<!-- An open conversation's entries are shown, a closed one's on request. -->
					<details class="mt-2" open={!c.closed}>
						<summary class="cursor-pointer text-xs text-muted"
							>{c.entries.length}
							{c.entries.length === 1 ? 'entry' : 'entries'}</summary
						>
						<ol class="mt-2 space-y-2">
							{#each c.entries as e, i (i)}
								<li class="flex flex-col items-start text-sm" data-entry={e.story}>
									<div class="text-xs text-muted">
										<span class="font-mono">{localTime(e.at)}</span>
										{e.author}
										{#if e.story}<span class="font-mono">{e.story}</span>{/if}
									</div>
									<!-- Entries are markdown, as they are in the conversation's file; authors write
									     repository paths, so relative links resolve from the root. -->
									<div
										class="prose prose-sm mt-0.5 max-w-[85%] rounded-lg border border-line bg-raised px-3 py-2 [&>:first-child]:mt-0 [&>:last-child]:mb-0"
									>
										<!-- eslint-disable-next-line svelte/no-at-html-tags -- repository markdown, rendered client side as every document is -->
										{@html render(e.text, '')}
									</div>
								</li>
							{/each}
						</ol>
					</details>
				</li>
			{/each}
		</ul>
	{/if}
</section>
