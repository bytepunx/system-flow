<script lang="ts">
	// Move stories of one lane a column forward or back, from the lane's menu (S-0167). The
	// designer ticks which; the card the menu was opened on starts ticked. Closing changes nothing.
	import { untrack } from 'svelte';
	import { needsReason } from '$lib/lanes';

	type Story = { id: string; title: string };

	let {
		from,
		to,
		stories,
		picked = [],
		onconfirm,
		oncancel
	}: {
		from: string;
		to: string;
		stories: Story[];
		picked?: string[];
		onconfirm: (ids: string[], reason: string) => void | Promise<void>;
		oncancel: () => void;
	} = $props();

	// Ticked when the dialog opens; a board that reloads meanwhile does not undo the designer's ticks.
	let ticked = $state<string[]>(
		untrack(() => picked.filter((id) => stories.some((s) => s.id === id)))
	);
	let reason = $state('');
	let busy = $state(false);

	const asking = $derived(needsReason(from, to));
	const all = $derived(stories.length > 0 && ticked.length === stories.length);
	const can = $derived(!busy && ticked.length > 0 && (!asking || reason.trim() !== ''));

	async function confirm() {
		busy = true;
		try {
			// in the lane's order, whatever order they were ticked in
			await onconfirm(
				stories.filter((s) => ticked.includes(s.id)).map((s) => s.id),
				reason.trim()
			);
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
		class="w-full max-w-lg rounded-lg border border-line bg-surface p-5 text-sm shadow-lg"
		role="dialog"
		aria-modal="true"
		aria-labelledby="lane-move-title"
		data-testid="lane-move"
	>
		<h2 id="lane-move-title" class="text-base font-semibold">Move stories from {from} to {to}</h2>
		{#if stories.length}
			<label class="mt-3 flex items-center gap-2 font-medium">
				<input
					type="checkbox"
					checked={all}
					data-testid="lane-move-all"
					onchange={(e) => (ticked = e.currentTarget.checked ? stories.map((s) => s.id) : [])}
				/>
				All {stories.length}
			</label>
			<ul class="mt-1 max-h-72 overflow-y-auto rounded border border-line p-2">
				{#each stories as s (s.id)}
					<li>
						<label class="flex items-baseline gap-2">
							<input type="checkbox" value={s.id} bind:group={ticked} />
							<span class="font-mono">{s.id}</span>
							<span>{s.title}</span>
						</label>
					</li>
				{/each}
			</ul>
			{#if asking}
				<label class="mt-3 block">
					<span class="font-medium">Reason</span>
					<textarea
						class="mt-1 w-full rounded border border-line-strong bg-surface p-2"
						rows="2"
						bind:value={reason}
						data-testid="lane-move-reason"
						placeholder="Why it goes back; recorded in each story's notes"></textarea>
				</label>
			{/if}
		{:else}
			<p class="mt-3 text-muted">No stories in {from}.</p>
		{/if}
		<div class="mt-4 flex justify-end gap-2">
			<button
				type="button"
				class="rounded border border-line-strong px-3 py-1"
				disabled={busy}
				onclick={oncancel}>Close</button
			>
			<button
				type="button"
				class="rounded bg-primary px-3 py-1 text-on-primary disabled:opacity-50"
				disabled={!can}
				data-testid="lane-move-confirm"
				onclick={confirm}
				>{busy
					? 'Moving…'
					: `Move ${ticked.length} ${ticked.length === 1 ? 'story' : 'stories'} to ${to}`}</button
			>
		</div>
	</div>
</div>
