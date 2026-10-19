<script lang="ts" module>
	/** Why a story in review is not eligible for the merge queue, as flai names it. */
	export type QueueReason = {
		kind: string;
		text: string;
		criteria?: number[];
		threads?: string[];
		paths?: string[];
	};
	/** A story in review as the merge queue sees it. */
	export type QueueEntry = {
		story: string;
		title: string;
		entered_review: string;
		head?: string;
		enqueued?: boolean;
		reasons?: QueueReason[];
	};
	/** The merge queue as board.get carries it (S-0398, mergequeue.Queue). */
	export type MergeQueueState = {
		mode: string;
		off?: boolean;
		off_reason?: string;
		eligible: QueueEntry[] | null;
		ineligible: QueueEntry[] | null;
		enqueued?: { story: string; by: string; at: string }[] | null;
		batch?: { stories: string[]; step: string };
		outcomes?: { story: string; outcome: string; at: string; detail?: string }[] | null;
	};
	/** A card's mark in review: its place among the eligible stories, or why it is not eligible. */
	export type QueueMark = { place: number } | { reasons: QueueReason[] };

	/** Each story in review's mark, by ID; none while the queue is off or absent. */
	export function queueMarks(q: MergeQueueState | undefined): Record<string, QueueMark> {
		if (!q || q.off) return {};
		return Object.fromEntries([
			...(q.eligible ?? []).map((e, i) => [e.story, { place: i + 1 }]),
			...(q.ineligible ?? []).map((e) => [e.story, { reasons: e.reasons ?? [] }])
		]);
	}
</script>

<script lang="ts">
	// The merge queue above the review lane (S-0398): the eligible stories in the order the queue
	// takes them, and the batch under verification with its step. Under pull-request acceptance the
	// queue does not run, and flai says why; a flai on the host too old to send the queue sends none.
	// Each card in review carries its own mark, so the reasons of the stories not eligible are there.
	type Props = {
		/** The queue from board.get, or undefined when the flai on the host sent none. */
		queue?: MergeQueueState;
	};
	let { queue }: Props = $props();

	const eligible = $derived(queue?.eligible ?? []);
	const ineligible = $derived(queue?.ineligible ?? []);
</script>

<section
	class="mb-2 rounded border border-line p-2 text-xs text-muted"
	aria-label="Merge queue"
	data-testid="merge-queue"
>
	{#if !queue}
		<p data-testid="merge-queue-absent">
			No merge queue: the flai on the host sends none, being older than the queue or unable to read
			it.
		</p>
	{:else if queue.off}
		<div data-testid="merge-queue-off">
			<p class="font-semibold text-ink">
				The merge queue does not run under {queue.mode} acceptance.
			</p>
			<p class="mt-1">Why: {queue.off_reason}</p>
		</div>
	{:else}
		<p class="font-semibold text-ink">Merge queue</p>
		{#if queue.batch}
			<p class="mt-1" data-testid="merge-queue-batch">
				Batch <span class="font-mono text-ink">{queue.batch.stories.join(', ')}</span> at step
				<span class="font-semibold text-ink">{queue.batch.step}</span>
			</p>
		{:else}
			<p class="mt-1" data-testid="merge-queue-batch">No batch.</p>
		{/if}
		{#if eligible.length}
			<ol class="mt-1 space-y-0.5" data-testid="merge-queue-order">
				{#each eligible as e, i (e.story)}
					<li data-story={e.story}>
						{i + 1}. <span class="font-mono text-ink">{e.story}</span>
						{e.title}{#if e.enqueued}<span class="font-semibold text-good">, enqueued</span>{/if}
					</li>
				{/each}
			</ol>
		{:else}
			<p class="mt-1" data-testid="merge-queue-empty">Empty: no story in review is eligible.</p>
		{/if}
		{#if ineligible.length}
			<p class="mt-1" data-testid="merge-queue-ineligible">
				{ineligible.length} not eligible: the reasons are on {ineligible.length === 1
					? 'its card'
					: 'their cards'}.
			</p>
		{/if}
	{/if}
</section>
