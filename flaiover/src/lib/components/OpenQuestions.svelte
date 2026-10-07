<script lang="ts">
	// A story's open questions, hand-written in its narrative (S-0173): read from the shared inbox,
	// each answerable here as on the inbox. The one an inbox link names is marked and brought into
	// view; a question answered since says so rather than vanishing without a word.
	import { inboxState } from '$lib/inbox.svelte';
	import { tick } from 'svelte';
	import { localTime } from '$lib/localtime';
	import QuestionAnswer from './QuestionAnswer.svelte';

	let {
		story,
		writable = false,
		select
	}: {
		story: string;
		writable?: boolean;
		/** The inbox key of the question a link names, such as question:S-0001:abc. */
		select?: string;
	} = $props();

	const questions = $derived(
		(inboxState.data?.entries ?? []).filter((e) => e.kind === 'question' && e.item === story)
	);
	const gone = $derived(!!select && !!inboxState.data && !questions.some((q) => q.key === select));

	let list = $state<HTMLElement>();
	let shown: string | undefined;
	$effect(() => {
		if (!select || select === shown || !questions.some((q) => q.key === select)) return;
		shown = select;
		void tick().then(() =>
			list
				?.querySelector(`[data-question="${CSS.escape(select)}"]`)
				?.scrollIntoView?.({ block: 'center' })
		);
	});
</script>

{#if questions.length || gone}
	<section class="mt-6 text-sm" data-testid="open-questions">
		<h2 class="font-medium">
			Open questions <span class="text-xs font-normal text-muted">{questions.length}</span>
		</h2>
		<p class="text-xs text-muted">Asked in the narrative; an answer moves to its Decisions.</p>
		{#if gone}
			<p class="mt-2 text-xs text-muted" data-testid="question-gone">
				The question the link named is answered or no longer open.
			</p>
		{/if}
		<ul class="mt-2 space-y-1" bind:this={list}>
			{#each questions as q (q.key)}
				<li
					class="rounded border bg-surface px-3 py-2 {q.key === select
						? 'border-primary ring-1 ring-primary'
						: 'border-line'}"
					data-question={q.key}
					aria-current={q.key === select ? 'true' : undefined}
				>
					<span>{q.title}</span>
					{#if q.at}<span class="block text-xs text-muted">{localTime(q.at)}</span>{/if}
					{#if writable}<QuestionAnswer entry={q} />{/if}
				</li>
			{/each}
		</ul>
	</section>
{/if}
