<script lang="ts">
	// What needs a human (S-0042), grouped by kind, each entry a link to its page.
	import type { Inbox, InboxEntry } from '$lib/inbox.svelte';

	let { inbox }: { inbox: Inbox } = $props();

	const kinds: { kind: InboxEntry['kind']; label: string }[] = [
		{ kind: 'review', label: 'Stories in review' },
		{ kind: 'thread', label: 'Threads awaiting you' },
		{ kind: 'question', label: 'Open questions in narratives' },
		{ kind: 'blocked', label: 'Blocked items' },
		{ kind: 'overlap', label: 'Overlapping touches' }
	];
	const of = (k: InboxEntry['kind']) => inbox.entries.filter((e) => e.kind === k);
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
					</li>
				{/each}
			</ul>
		</section>
	{/if}
{/each}
{#each inbox.notes as n (n)}
	<p class="text-xs text-muted">{n}</p>
{/each}
