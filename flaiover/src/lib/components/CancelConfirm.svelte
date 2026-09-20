<script lang="ts">
	// Cancelling an item cancels everything open under it (S-0070, ADR-0028):
	// show what goes with it, and take the reason, before it happens.
	// Closing the dialog changes nothing.
	import { api } from '$lib/api';

	type Cascaded = { id: string; type: string; title: string; from: string };
	type Left = { id: string; narrative?: string; branch?: string; worktree?: string };
	type Preview = { id: string; cancelled?: Cascaded[]; left_behind?: Left[] };

	let {
		id,
		onconfirm,
		oncancel
	}: {
		id: string;
		onconfirm: (reason: string) => void | Promise<void>;
		oncancel: () => void;
	} = $props();

	let preview = $state<Preview | null>(null);
	let error = $state<string | null>(null);
	let busy = $state(false);
	let reason = $state('');

	const items = $derived(preview?.cancelled ?? []);
	const inReview = $derived(items.filter((c) => c.from === 'review'));
	const kept = (l: Left) =>
		[l.narrative && 'narrative', l.branch && `branch ${l.branch}`, l.worktree && 'worktree']
			.filter(Boolean)
			.join(', ');

	$effect(() => {
		const target = id;
		preview = null;
		error = null;
		api(`/api/items/${target}/move`, {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ to: 'cancelled', dry_run: true })
		})
			.then(async (r) => {
				const data = await r.json();
				if (!r.ok) error = data.error ?? r.statusText;
				else preview = data;
			})
			.catch((e) => (error = e instanceof Error ? e.message : String(e)));
	});

	async function confirm() {
		busy = true;
		try {
			await onconfirm(reason.trim());
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
>
	<div
		class="w-full max-w-lg rounded-lg border border-line bg-surface p-5 text-sm shadow-lg"
		role="dialog"
		aria-modal="true"
		aria-labelledby="cancel-title"
	>
		<h2 id="cancel-title" class="text-base font-semibold">Cancel {id}?</h2>
		<p class="mt-1 text-muted">Cancelled is final: a cancelled item cannot be moved again.</p>
		{#if error}
			<p class="mt-3 rounded border border-danger bg-danger-soft p-2 text-danger" role="alert">
				{error}
			</p>
		{:else if !preview}
			<p class="mt-3 text-muted">Working out what goes with it…</p>
		{:else}
			{#if items.length}
				<p class="mt-3 font-medium">
					This also cancels {items.length}
					{items.length === 1 ? 'item' : 'items'} still open under it:
				</p>
				<ul class="mt-1 max-h-56 overflow-y-auto rounded border border-line p-2" data-cascade>
					{#each items as c (c.id)}
						<li class={c.type === 'task' ? 'ml-4' : ''}>
							<span class="font-mono">{c.id}</span>
							<span class="text-muted">{c.type}, {c.from}</span>
							{c.title}
						</li>
					{/each}
				</ul>
				{#if inReview.length}
					<p class="mt-2 rounded border border-warn bg-warn-soft p-2" role="note">
						{inReview.map((c) => c.id).join(', ')}
						{inReview.length === 1 ? 'is' : 'are'} in review. The work stays on its branch, unmerged;
						accept it first if you want it.
					</p>
				{/if}
			{:else}
				<p class="mt-3 text-muted">Nothing open under it; only {id} changes.</p>
			{/if}
			{#if preview.left_behind?.length}
				<p class="mt-2 text-muted">
					Left as it is, for you to keep or remove:
					{#each preview.left_behind as l, i (l.id)}{i ? '; ' : ' '}{l.id} ({kept(l)}){/each}.
				</p>
			{/if}
			<label class="mt-3 block">
				<span class="font-medium">Reason</span>
				<textarea
					class="mt-1 w-full rounded border border-line-strong bg-surface p-2"
					rows="2"
					bind:value={reason}
					placeholder="Recorded in the notes of every item cancelled"></textarea>
			</label>
		{/if}
		<div class="mt-4 flex justify-end gap-2">
			<button
				type="button"
				class="rounded border border-line-strong px-3 py-1"
				disabled={busy}
				onclick={oncancel}>Keep it</button
			>
			<button
				type="button"
				class="rounded bg-danger px-3 py-1 text-on-primary disabled:opacity-50"
				disabled={busy || !!error || !preview || !reason.trim()}
				onclick={confirm}
				>{busy
					? 'Cancelling…'
					: items.length
						? `Cancel ${id} and ${items.length} more`
						: `Cancel ${id}`}</button
			>
		</div>
	</div>
</div>
