<script lang="ts">
	import { untrack } from 'svelte';
	// Change a lane's WIP limit from its menu (S-0167). flai writes it to wip/kanban/board.md, the
	// one place it is read from; 0 is no limit. Closing changes nothing.
	let {
		lane,
		limit,
		count,
		onconfirm,
		oncancel
	}: {
		lane: string;
		/** the limit now; 0 or none is no limit */
		limit: number | undefined;
		/** the stories in the lane now */
		count: number;
		onconfirm: (limit: number) => void | Promise<void>;
		oncancel: () => void;
	} = $props();

	// A number input binds a number, or null while it is empty. It starts from the limit when the
	// dialog opens; a board that reloads meanwhile does not undo what is typed.
	let value = $state<number | null>(untrack(() => limit ?? 0));
	let busy = $state(false);

	const valid = $derived(Number.isInteger(value) && value! >= 0 && value! <= 99);

	async function confirm(e: Event) {
		e.preventDefault();
		if (!valid || busy) return;
		busy = true;
		try {
			await onconfirm(value!);
		} finally {
			busy = false;
		}
	}
</script>

<div
	class="fixed inset-0 z-50 flex items-start justify-center bg-ink/40 p-4 pt-24"
	role="presentation"
	onclick={(e) => {
		if (e.target === e.currentTarget && !busy) oncancel();
	}}
	onkeydown={(e) => {
		if (e.key === 'Escape' && !busy) oncancel();
	}}
>
	<div
		class="w-full max-w-sm rounded-lg border border-line bg-surface p-5 text-sm shadow-lg"
		role="dialog"
		aria-modal="true"
		aria-labelledby="wip-limit-title"
		data-testid="wip-limit"
	>
		<form onsubmit={confirm}>
			<h2 id="wip-limit-title" class="text-base font-semibold">WIP limit for {lane}</h2>
			<p class="mt-1 text-muted">
				{count}
				{count === 1 ? 'story is' : 'stories are'} in {lane} now. 0 means no limit. A move past the limit
				warns; flai serve starts agents only while in-progress has room.
			</p>
			<label class="mt-3 block">
				<span class="font-medium">Limit</span>
				<input
					class="mt-1 w-24 rounded border border-line-strong bg-surface px-2 py-1"
					type="number"
					min="0"
					max="99"
					step="1"
					bind:value
					data-testid="wip-limit-value"
				/>
			</label>
			{#if !valid}
				<p class="mt-1 text-danger" role="alert">A whole number from 0 to 99.</p>
			{/if}
			<div class="mt-4 flex justify-end gap-2">
				<button
					type="button"
					class="rounded border border-line-strong px-3 py-1"
					disabled={busy}
					onclick={oncancel}>Close</button
				>
				<button
					class="rounded bg-primary px-3 py-1 text-on-primary disabled:opacity-50"
					disabled={!valid || busy}
					data-testid="wip-limit-confirm">{busy ? 'Saving…' : 'Set limit'}</button
				>
			</div>
		</form>
	</div>
</div>
