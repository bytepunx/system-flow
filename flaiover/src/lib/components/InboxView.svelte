<script lang="ts">
	// What needs a human (S-0042), grouped by kind, each entry a link to its page, never to a
	// document (S-0173). A hand-written open question (kind "question") has no thread of its own to
	// reply to, so it is answered right here (S-0090) as well as on its story's page. A thread's
	// recommendation awaiting the designer is shown with its source and confirmed right here too
	// (ADR-0090).
	import type { Inbox, InboxEntry } from '$lib/inbox.svelte';
	import { inboxState } from '$lib/inbox.svelte';
	import { api } from '$lib/api';
	import { resolve } from '$app/paths';
	import { slug } from '$lib/markdown';
	import { localTime } from '$lib/localtime';
	import QuestionAnswer from './QuestionAnswer.svelte';

	let { inbox, writable = false }: { inbox: Inbox; writable?: boolean } = $props();

	const kinds: { kind: InboxEntry['kind']; label: string }[] = [
		{ kind: 'review', label: 'Stories in review' },
		{ kind: 'thread', label: 'Threads awaiting you' },
		{ kind: 'question', label: 'Open questions in narratives' },
		{ kind: 'blocked', label: 'Blocked items' },
		{ kind: 'overlap', label: 'Overlapping touches' }
	];
	const of = (k: InboxEntry['kind']) => inbox.entries.filter((e) => e.kind === k);

	type Source = NonNullable<NonNullable<InboxEntry['recommendation']>['source']>;
	const sourceHref = (s: Source) =>
		resolve('/docs/[...path]', { path: s.path }) + (s.heading ? `#${slug(s.heading)}` : '');

	let confirming = $state<string | null>(null);
	let refused = $state<Record<string, string>>({});
	/** Confirm the thread's pending recommendation as its answer; the entry leaves on the refresh. */
	async function confirm(e: InboxEntry) {
		confirming = e.key;
		delete refused[e.key];
		try {
			const r = await api(`/api/threads/${e.key.replace(/^thread:/, '')}/confirm`, {
				method: 'POST'
			});
			if (!r.ok) {
				refused[e.key] = (await r.json().catch(() => ({}))).error ?? r.statusText;
				return;
			}
			await inboxState.refresh();
		} finally {
			confirming = null;
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
						{#if e.at}<span class="block text-xs text-muted">{localTime(e.at)}</span>{/if}
						{#if e.recommendation}
							{@const rec = e.recommendation}
							<div class="mt-2 border-l-2 border-info pl-2 text-xs" data-recommendation={e.key}>
								<span class="text-[10px] text-info uppercase">recommendation</span>
								{rec.author}
								<p class="mt-0.5 line-clamp-4 whitespace-pre-line text-ink">{rec.text}</p>
								{#if rec.source}
									<!-- eslint-disable svelte/no-navigation-without-resolve -- the path is resolve()d; the rule does not follow the heading added to it -->
									<p class="mt-0.5 text-muted" data-source>
										Source:
										<a class="underline" href={sourceHref(rec.source)}
											>{rec.source.path}{rec.source.heading ? ` § ${rec.source.heading}` : ''}</a
										>
									</p>
									<!-- eslint-enable svelte/no-navigation-without-resolve -->
								{/if}
								{#if writable}
									<button
										type="button"
										class="mt-1 rounded bg-primary px-2 py-1 text-on-primary disabled:opacity-50"
										data-confirm={e.key}
										disabled={confirming === e.key}
										onclick={() => confirm(e)}
										>{confirming === e.key ? 'confirming…' : 'Confirm'}</button
									>
								{/if}
								{#if refused[e.key]}
									<p class="mt-1 text-danger" role="alert">{refused[e.key]}</p>
								{/if}
							</div>
						{/if}
						{#if kind === 'question' && writable && e.item}<QuestionAnswer entry={e} />{/if}
					</li>
				{/each}
			</ul>
		</section>
	{/if}
{/each}
{#each inbox.notes as n (n)}
	<p class="text-xs text-muted">{n}</p>
{/each}
