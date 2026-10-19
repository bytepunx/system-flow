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
	/** What the batch's verification found (S-0399): passed, or the tier it stopped at and why. */
	export type BatchResult = {
		passed: boolean;
		stopped_at?: string;
		summary?: string;
		finished_at: string;
	};
	/** One step of a red batch's bisection (S-0400, mergequeue.BisectStep): main with the batch's first stories, verified at commit; no result while it runs. */
	export type BisectStep = {
		stories: string[];
		commit: string;
		started_at: string;
		result?: BatchResult;
	};
	/** A red batch's bisection (mergequeue.Bisect): its steps in the order they ran, and, once found, the first story with which the batch fails. */
	export type QueueBisect = {
		steps: BisectStep[] | null;
		good: number;
		bad: number;
		failing?: string;
		commit?: string;
		result?: BatchResult;
	};
	/** A verified batch's landing (mergequeue.Landing): the stories it lands in order, those landed so far, and main's head the next landing must find. */
	export type QueueLanding = { stories: string[] | null; landed: string[] | null; head: string };
	/** The batch under verification (mergequeue.Batch): its stories, its step, and, once it has them, its result and why it was dropped; since S-0400, each story's heads, a red batch's bisection, and the landing. */
	export type QueueBatch = {
		id?: string;
		branch?: string;
		worktree?: string;
		stories: string[];
		started_at?: string;
		step: string;
		result?: BatchResult;
		dropped?: { reason: string; by: string; at: string };
		heads?: { story: string; verified: string; rebased: string }[] | null;
		bisect?: QueueBisect;
		landing?: QueueLanding;
	};
	/** A story's last outcome in the queue (mergequeue.Outcome): landed, failing, left_out, or accepted_by_hand, among others. */
	export type QueueOutcome = {
		story: string;
		outcome: string;
		at: string;
		detail?: string;
		base?: string;
		head?: string;
	};
	/** The merge queue as board.get carries it (S-0398, mergequeue.Queue), held or not (S-0399). */
	export type MergeQueueState = {
		mode: string;
		off?: boolean;
		off_reason?: string;
		eligible: QueueEntry[] | null;
		ineligible: QueueEntry[] | null;
		enqueued?: { story: string; by: string; at: string }[] | null;
		batch?: QueueBatch;
		outcomes?: QueueOutcome[] | null;
		held?: boolean;
		held_by?: string;
		held_at?: string;
	};
	/** A write to the queue refused (mergequeue.Refusal): flai's sentence, and every reason a story is not eligible. */
	export type QueueRefusal = { text: string; reasons: QueueReason[] };
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

	/**
	 * Whether a batch is under way: not dropped, so its step may change. Since S-0400 a batch with a
	 * result is still bisected or landed until flai serve removes it.
	 */
	export function verifying(q: MergeQueueState | null | undefined): boolean {
		const b = q && !q.off ? q.batch : undefined;
		return !!b && !b.dropped;
	}

	/** What a batched story comes to in its batch (S-0400), as the panel names it. */
	export type BatchMark = 'landed' | 'sent back' | 'to land' | 'still enqueued';

	/**
	 * What story comes to in batch: landed once the landing has landed it; sent back when the
	 * bisection found it failing, or its last outcome since the batch started is failing; to land
	 * while the landing has it still to land; still enqueued when it comes after the failing story
	 * and the batch does not land it. None before the batch has come to any of these.
	 */
	export function batchMark(
		batch: QueueBatch,
		story: string,
		outcomes?: QueueOutcome[] | null
	): BatchMark | undefined {
		const since = batch.started_at ? Date.parse(batch.started_at) : -Infinity;
		const failing =
			batch.bisect?.failing ??
			batch.stories.find((s) =>
				(outcomes ?? []).some(
					(o) => o.story === s && o.outcome === 'failing' && !(Date.parse(o.at) < since)
				)
			);
		if (batch.landing?.landed?.includes(story)) return 'landed';
		if (story === failing) return 'sent back';
		if (batch.landing?.stories?.includes(story)) return 'to land';
		if (failing && batch.stories.indexOf(story) > batch.stories.indexOf(failing))
			return 'still enqueued';
		return undefined;
	}

	/** A failed write's answer as the page shows it: flai's sentence and reasons when it refused, else its error. */
	export function refusalOf(body: unknown, fallback: string): QueueRefusal {
		const b = (body ?? {}) as { text?: unknown; reasons?: unknown; error?: unknown };
		const text =
			typeof b.text === 'string' && b.text
				? b.text
				: typeof b.error === 'string' && b.error
					? b.error
					: fallback;
		return { text, reasons: Array.isArray(b.reasons) ? (b.reasons as QueueReason[]) : [] };
	}
</script>

<script lang="ts">
	// The merge queue above the review lane (S-0398): the eligible stories in the order the queue
	// takes them, and the batch under verification with its step and, once verified, its result
	// (S-0399). Under pull-request acceptance the queue does not run, and flai says why; a flai on the
	// host too old to send the queue sends none. Each card in review carries its own mark, so the
	// reasons of the stories not eligible are there. The hold switch holds and resumes the queue
	// through flai's queue.hold, offered while the merge-queue host action is on. A red batch's
	// bisection shows each step and the story it found failing, and each batched story is marked
	// with what it comes to in the batch: landed, sent back, to land, or still enqueued (S-0400).
	import { api } from '$lib/api';
	import { localTime } from '$lib/localtime';

	type Props = {
		/** The queue from board.get, or undefined when the flai on the host sent none. */
		queue?: MergeQueueState;
		/** Whether the merge-queue host action is on; undefined while unknown, and then nothing is offered. */
		enabled?: boolean;
		/** Whether the dashboard can write through flai. */
		writable?: boolean;
		/** Called once a hold or a resume is written, so that the page reads the queue again. */
		onchanged?: () => void;
	};
	let { queue, enabled, writable = false, onchanged }: Props = $props();

	const eligible = $derived(queue?.eligible ?? []);
	const ineligible = $derived(queue?.ineligible ?? []);
	const batch = $derived(queue?.batch);

	// Each mark's colour; its words carry its meaning.
	const markTone: Record<BatchMark, string> = {
		landed: 'text-good',
		'sent back': 'text-danger',
		'to land': 'text-ink',
		'still enqueued': 'text-warn'
	};

	// Each batched story with the comma before it, and its mark in words and colour after a space.
	const marked = $derived(
		(batch?.stories ?? []).map((story, i) => {
			const mark = batch ? batchMark(batch, story, queue?.outcomes) : undefined;
			return {
				story,
				mark,
				sep: i > 0 ? ', ' : '',
				gap: mark ? ' ' : '',
				tone: mark ? markTone[mark] : ''
			};
		})
	);
	const failingAt = $derived(batch?.bisect?.commit ? ` at ${batch.bisect.commit.slice(0, 7)}` : '');

	// A bisection step's result in a word and its colour, and, when it failed, where and why.
	function stepResult(r: BatchResult | undefined) {
		if (!r) return { word: 'running', tone: 'text-ink', why: '' };
		if (r.passed) return { word: 'passed', tone: 'text-good', why: '' };
		const at = r.stopped_at ? ` at ${r.stopped_at}` : '';
		const why = r.summary ? `${at}: ${r.summary}` : at;
		return { word: 'failed', tone: 'text-danger', why };
	}

	let holding = $state(false);
	let holdError = $state<string | null>(null);

	async function hold(held: boolean) {
		holding = true;
		holdError = null;
		try {
			const r = await api('/api/queue', {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ held })
			});
			const body = await r.json().catch(() => ({}));
			if (!r.ok) {
				holdError = refusalOf(body, r.statusText).text;
				return;
			}
			onchanged?.();
		} catch (e) {
			holdError = e instanceof Error ? e.message : String(e);
		} finally {
			holding = false;
		}
	}
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
		<div class="flex flex-wrap items-center justify-between gap-2">
			<p class="font-semibold text-ink">Merge queue</p>
			{#if enabled && writable}
				<label class="flex items-center gap-1 text-ink">
					<input
						type="checkbox"
						role="switch"
						checked={queue.held === true}
						disabled={holding}
						aria-label="Hold the merge queue"
						data-testid="merge-queue-hold"
						onchange={(e) => hold(e.currentTarget.checked)}
					/>
					Hold
				</label>
			{/if}
		</div>
		{#if queue.held}
			<p class="mt-1 font-semibold text-warn" data-testid="merge-queue-held">
				Held by {queue.held_by}{queue.held_at ? ` at ${localTime(queue.held_at)}` : ''}: no batch is
				built until it is resumed.
			</p>
		{/if}
		{#if enabled === false}
			<p class="mt-1" data-testid="merge-queue-disabled">
				No batch runs: the operator turns the queue on in a shell on the host with
				<code class="rounded bg-ground px-1 text-ink">flai serve enable merge-queue</code>.
			</p>
		{/if}
		{#if holdError}
			<p class="mt-1 text-danger" role="alert" data-testid="merge-queue-hold-error">
				{holdError}
			</p>
		{/if}
		{#if batch}
			<div class="mt-1" data-testid="merge-queue-batch">
				<p>
					Batch
					{#each marked as m (m.story)}{m.sep}<span class="font-mono text-ink" data-story={m.story}
							>{m.story}</span
						>{m.gap}{#if m.mark}<span class="font-semibold {m.tone}" data-testid="merge-queue-mark"
								>({m.mark})</span
							>{/if}{/each}
					at step
					<span class="font-semibold text-ink">{batch.step}</span>
				</p>
				{#if batch.result?.passed}
					<p class="text-good" data-testid="merge-queue-result">
						Passed, {localTime(batch.result.finished_at)}{batch.result.summary
							? `: ${batch.result.summary}`
							: '.'}
					</p>
				{:else if batch.result}
					<p class="text-danger" data-testid="merge-queue-result">
						Failed{batch.result.stopped_at ? ` at ${batch.result.stopped_at}` : ''}, {localTime(
							batch.result.finished_at
						)}{batch.result.summary ? `: ${batch.result.summary}` : '.'}
					</p>
				{/if}
				{#if batch.bisect}
					<ol class="mt-1 space-y-0.5" aria-label="Bisection" data-testid="merge-queue-bisect">
						{#each batch.bisect.steps ?? [] as st, i (i)}
							{@const r = stepResult(st.result)}
							<li data-step={i + 1}>
								Bisect {i + 1}: main with
								<span class="font-mono text-ink">{st.stories.join(', ')}</span> at
								<span class="font-mono">{st.commit.slice(0, 7)}</span>,
								<span class="font-semibold {r.tone}">{r.word}</span>{r.why}
							</li>
						{/each}
					</ol>
					{#if batch.bisect.failing}
						<p class="mt-1 text-danger" data-testid="merge-queue-failing">
							The batch fails with {batch.bisect.failing}{failingAt}: sent back to in-progress.
						</p>
					{/if}
				{/if}
				{#if batch.dropped}
					<p class="text-warn" data-testid="merge-queue-dropped">
						Dropped by {batch.dropped.by}, {localTime(batch.dropped.at)}: {batch.dropped.reason}
					</p>
				{/if}
			</div>
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
