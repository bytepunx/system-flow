<script lang="ts" module>
	/** One required check on the pull request's head commit, as flai serve last read it. */
	export type ForgeCheck = {
		name: string;
		state: 'pending' | 'passed' | 'failed';
		url?: string;
		summary?: string;
		since: string;
	};
	/** One review of the pull request, as flai serve last read it. */
	export type ForgeReview = {
		author: string;
		state: 'approved' | 'changes-requested' | 'commented' | 'dismissed' | 'pending';
		url?: string;
		body?: string;
		submitted_at?: string;
	};
	/** A failure flai sent the story back to in-progress for. */
	export type ForgeSendBack = {
		key: string;
		kind: 'check' | 'review';
		name: string;
		commit: string;
		url?: string;
		thread?: string;
		at: string;
	};
	/** A story's pull request as flai serve last read it from the forge (ADR-0167). */
	export type PullRequestStatus = {
		story: string;
		pull_request: number;
		url: string;
		head_sha: string;
		read_at: string;
		checks: ForgeCheck[];
		reviews: ForgeReview[];
		sent_back?: ForgeSendBack[];
	};
</script>

<script lang="ts">
	// A story's pull request, its required checks and its reviews (S-0394), as flai serve last read
	// them into kept state: the page never asks the forge. Each state is a word on its tone: passed
	// and approved good, failed and changes requested danger, pending warn, the rest plain.
	import { localTime } from '$lib/localtime';

	type Props = {
		/** The story's pull request, or null when flai serve has read none. */
		status: PullRequestStatus | null;
		/** Classes for the section's place on the page, such as its margin. */
		class?: string;
	};
	let { status, class: className = '' }: Props = $props();

	const tones: Record<string, string> = {
		passed: 'border-good bg-good-soft text-good',
		approved: 'border-good bg-good-soft text-good',
		failed: 'border-danger bg-danger-soft text-danger',
		'changes-requested': 'border-danger bg-danger-soft text-danger',
		pending: 'border-warn bg-warn-soft text-warn'
	};
	const tone = (state: string) => tones[state] ?? 'border-line text-muted';
	const word = (state: string) => state.replace(/-/g, ' ');
	const head = $derived(status?.head_sha.slice(0, 12) ?? '');
	const readAt = $derived(status ? localTime(status.read_at) : '');
</script>

{#if status}
	<section
		class="rounded border border-line bg-surface p-3 text-sm {className}"
		data-testid="forge-status"
	>
		<h2 class="mb-1 font-medium">
			<a class="underline" href={status.url} target="_blank" rel="external noopener noreferrer"
				>Pull request #{status.pull_request}</a
			>
		</h2>
		<p class="text-xs text-muted">
			Head <code class="rounded bg-ground px-1 text-ink">{head}</code>, read {readAt}.
		</p>

		<h3 class="mt-3 mb-1 font-medium">Required checks</h3>
		{#if status.checks.length}
			<ul class="space-y-1">
				{#each status.checks as check (check.name)}
					<li data-check={check.name}>
						<span class="flex flex-wrap items-baseline gap-2">
							<span
								class="rounded border px-1.5 text-xs {tone(check.state)}"
								data-state={check.state}>{word(check.state)}</span
							>
							<span class="font-medium">{check.name}</span>
							{#if check.url}<a
									class="text-xs underline"
									href={check.url}
									target="_blank"
									rel="external noopener noreferrer"
									aria-label="{check.name} log">log</a
								>{/if}
						</span>
						{#if check.summary}<p class="text-xs text-muted">{check.summary}</p>{/if}
					</li>
				{/each}
			</ul>
		{:else}
			<p class="text-xs text-muted">None required.</p>
		{/if}

		<h3 class="mt-3 mb-1 font-medium">Reviews</h3>
		{#if status.reviews.length}
			<ul class="space-y-1">
				{#each status.reviews as review, i (i)}
					<li class="flex flex-wrap items-baseline gap-2" data-review={review.author}>
						<span
							class="rounded border px-1.5 text-xs {tone(review.state)}"
							data-state={review.state}>{word(review.state)}</span
						>
						<span class="font-medium">{review.author}</span>
						{#if review.url}<a
								class="text-xs underline"
								href={review.url}
								target="_blank"
								rel="external noopener noreferrer"
								aria-label="{review.author}'s review">review</a
							>{/if}
					</li>
				{/each}
			</ul>
		{:else}
			<p class="text-xs text-muted">None yet.</p>
		{/if}
	</section>
{/if}
