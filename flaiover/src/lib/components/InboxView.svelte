<script lang="ts">
	// What needs a human (S-0042), grouped by kind, each entry a link to its page.
	// A hand-written open question (kind "question") can also be answered right
	// here: it has no thread of its own to reply to, so this is its only place
	// in the dashboard (S-0090). Answering moves it out of Open questions into
	// the narrative's Decisions, and it drops out of this list once it does.
	import type { Inbox, InboxEntry } from '$lib/inbox.svelte';
	import { inboxState } from '$lib/inbox.svelte';
	import { api } from '$lib/api';

	let { inbox, writable = false }: { inbox: Inbox; writable?: boolean } = $props();

	const kinds: { kind: InboxEntry['kind']; label: string }[] = [
		{ kind: 'review', label: 'Stories in review' },
		{ kind: 'thread', label: 'Threads awaiting you' },
		{ kind: 'question', label: 'Open questions in narratives' },
		{ kind: 'blocked', label: 'Blocked items' },
		{ kind: 'overlap', label: 'Overlapping touches' }
	];
	const of = (k: InboxEntry['kind']) => inbox.entries.filter((e) => e.kind === k);

	let drafts = $state<Record<string, string>>({});
	let answering = $state<string | null>(null);
	let failed = $state<Record<string, string>>({});

	async function answer(e: InboxEntry) {
		const text = drafts[e.key]?.trim();
		if (!text || !e.item) return;
		answering = e.key;
		failed = { ...failed, [e.key]: '' };
		try {
			const r = await api(`/api/streams/${e.item}/answer`, {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ question: e.title, answer: text })
			});
			const data = await r.json().catch(() => ({}));
			if (!r.ok) {
				failed = { ...failed, [e.key]: data.error ?? r.statusText };
				return;
			}
			drafts = { ...drafts, [e.key]: '' };
			await inboxState.refresh();
		} finally {
			answering = null;
		}
	}
</script>

{#if !inbox.total}
	<p class="text-sm text-muted">Nothing needs you right now.</p>
{/if}
{#each kinds as { kind, label } (kind)}
	{#if of(kind).length}
		<section class="mb-4">
			<h2 class="mb-1 text-sm font-medium">
				{label} <span class="text-xs font-normal text-muted">{of(kind).length}</span>
			</h2>
			<ul class="space-y-1">
				{#each of(kind) as e (e.key)}
					<li class="rounded border border-line bg-surface px-3 py-2 text-sm">
						<!-- eslint-disable-next-line svelte/no-navigation-without-resolve -- the server builds these dashboard paths -->
						<a class="underline" href={e.href}>{e.title}</a>
						{#if e.detail}<span class="block text-xs text-muted">{e.detail}</span>{/if}
						{#if e.at}<span class="block text-xs text-muted">{e.at}</span>{/if}
						{#if kind === 'question' && writable && e.item}
							<form
								class="mt-2 flex gap-2"
								onsubmit={(ev) => {
									ev.preventDefault();
									void answer(e);
								}}
							>
								<input
									class="min-w-0 flex-1 rounded border border-line-strong px-2 py-1 text-xs"
									placeholder="answer"
									data-testid="question-answer-input"
									bind:value={drafts[e.key]}
								/>
								<button
									type="submit"
									class="rounded border border-line-strong px-2 py-1 text-xs disabled:opacity-50"
									disabled={!drafts[e.key]?.trim() || answering === e.key}
									data-testid="question-answer-submit"
									>{answering === e.key ? 'answering…' : 'answer'}</button
								>
							</form>
							{#if failed[e.key]}
								<p class="mt-1 text-xs text-danger" role="alert">{failed[e.key]}</p>
							{/if}
						{/if}
					</li>
				{/each}
			</ul>
		</section>
	{/if}
{/each}
{#each inbox.notes as n (n)}
	<p class="text-xs text-muted">{n}</p>
{/each}
