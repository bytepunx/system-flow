<script lang="ts">
	// Moving a story to done is acceptance (S-0046): show what it will do
	// before it happens. Cancelling changes nothing; the card stays in review.
	import UnreleasedList from './UnreleasedList.svelte';
	import { api } from '$lib/api';

	type Version = string | { Major: number; Minor: number; Patch: number };
	type Step = {
		component: { name: string };
		delivered: boolean;
		level: string;
		from: Version;
		to: Version;
		tag?: string;
		version?: string;
		files: string[];
	};
	type Preview = {
		id: string;
		branch?: string;
		resumed?: boolean;
		blockers?: string[];
		uncommitted?: string[];
		/** Uncommitted paths in the story's worktree: they block acceptance (S-0140). */
		worktree_uncommitted?: string[];
		plan?: {
			level: string;
			commits: string[];
			steps: Step[];
			skipped?: string;
			/** Components a no-release story touched: code landing on main without a version (ADR-0025). */
			unreleased?: { component: string; files: string[] }[];
		} | null;
	};

	let {
		id,
		onconfirm,
		oncancel
	}: {
		id: string;
		/** include is true when the designer chose to put the uncommitted files in the acceptance commit. */
		onconfirm: (include: boolean) => void | Promise<void>;
		oncancel: () => void;
	} = $props();

	let preview = $state<Preview | null>(null);
	let error = $state<string | null>(null);
	let busy = $state(false);
	// Uncommitted files outside wip: off by default, and accept waits for the choice (S-0051).
	let include = $state(false);
	const needsChoice = $derived(!!preview?.uncommitted?.length && !include);

	// The worktree's uncommitted paths have their own box, with the way to have them committed
	// (S-0140); flai's blocker saying the same is not repeated above it.
	const worktreeBlocker = (b: string) =>
		!!preview?.worktree_uncommitted?.length &&
		b.startsWith('the worktree ') &&
		b.includes('has uncommitted changes');
	const blockers = $derived((preview?.blockers ?? []).filter((b) => !worktreeBlocker(b)));
	let committing = $state(false);
	let committed = $state<string | null>(null);
	let commitError = $state<string | null>(null);

	async function haveCommitted() {
		committing = true;
		commitError = null;
		try {
			const r = await api(`/api/items/${id}/agent`, {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ action: 'commit' })
			});
			const data = await r.json();
			if (!r.ok) commitError = data.error ?? r.statusText;
			else
				committed = `${data.agent ?? 'an agent'} is committing them (pid ${data.pid}). Accept once it has finished: open this again to see.`;
		} catch (e) {
			commitError = e instanceof Error ? e.message : String(e);
		} finally {
			committing = false;
		}
	}

	const v = (x: Version) => (typeof x === 'string' ? x : `${x.Major}.${x.Minor}.${x.Patch}`);

	$effect(() => {
		const target = id;
		preview = null;
		error = null;
		include = false;
		committed = null;
		commitError = null;
		api(`/api/items/${target}/acceptance`)
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
			await onconfirm(include);
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
		aria-labelledby="accept-title"
	>
		<h2 id="accept-title" class="text-base font-semibold">Accept {id}?</h2>
		<p class="mt-1 text-muted">
			Moving a story to done accepts it. This cannot be undone from the dashboard.
		</p>
		{#if error}
			<p class="mt-3 rounded border border-danger bg-danger-soft p-2 text-danger" role="alert">
				{error}
			</p>
		{:else if !preview}
			<p class="mt-3 text-muted">Working out what acceptance will do…</p>
		{:else}
			{#if blockers.length}
				<div class="mt-3 rounded border border-danger bg-danger-soft p-2 text-danger" role="alert">
					<p class="font-medium">This cannot be accepted from here yet:</p>
					<ul class="mt-1 ml-4 list-disc">
						{#each blockers as b (b)}<li>{b}</li>{/each}
					</ul>
				</div>
			{/if}
			{#if preview.worktree_uncommitted?.length}
				<div
					class="mt-3 rounded border border-danger bg-danger-soft p-2 text-danger"
					role="group"
					data-testid="worktree-uncommitted"
				>
					<p class="font-medium">Uncommitted changes in the story's worktree:</p>
					<ul class="mt-1 ml-4 list-disc font-mono text-xs">
						{#each preview.worktree_uncommitted as p (p)}<li>{p}</li>{/each}
					</ul>
					<p class="mt-2">
						Acceptance merges the branch as it is committed and removes the worktree, so it waits
						until these are committed on {preview.branch ?? 'the story branch'}.
					</p>
					{#if committed}
						<p class="mt-2 font-medium" role="status">{committed}</p>
					{:else}
						<button
							type="button"
							class="mt-2 rounded border border-line-strong bg-surface px-2 py-1 text-ink hover:bg-raised disabled:opacity-50"
							disabled={committing || busy}
							onclick={haveCommitted}
							>{committing ? 'Starting an agent…' : 'Have an agent commit them'}</button
						>
					{/if}
					{#if commitError}<p class="mt-2" role="alert">{commitError}</p>{/if}
				</div>
			{/if}
			{#if preview.uncommitted?.length}
				<div class="mt-3 rounded border border-warn bg-warn-soft p-2 text-warn" role="group">
					<p class="font-medium">Uncommitted changes outside wip:</p>
					<ul class="mt-1 ml-4 list-disc font-mono text-xs">
						{#each preview.uncommitted as p (p)}<li>{p}</li>{/each}
					</ul>
					<p class="mt-2">
						Acceptance refuses these by default, so its commit holds only acceptance. Cancel to
						commit or stash them first, or include them.
					</p>
					<label class="mt-2 flex items-center gap-2">
						<input type="checkbox" bind:checked={include} disabled={busy} />
						Include these files in the acceptance commit
					</label>
				</div>
			{/if}
			<ul class="mt-3 space-y-1">
				{#if preview.resumed}
					<li>{id} is already marked done; this completes its acceptance.</li>
				{/if}
				{#if preview.branch}
					<li>
						Rebase <code>{preview.branch}</code>, fast-forward it into the main branch, and remove
						its worktree.
					</li>
				{/if}
				<li>Archive the story, its tasks, and its narrative, and commit.</li>
				{#if preview.plan?.skipped}
					<li>
						No release: {preview.plan.skipped}.
						<UnreleasedList unreleased={preview.plan.unreleased} />
					</li>
				{:else if preview.plan}
					<li>
						Tag and push a {preview.plan.level} release from {preview.plan.commits.length}
						commit{preview.plan.commits.length === 1 ? '' : 's'}:
						<ul class="mt-1 ml-4 list-disc">
							{#each preview.plan.steps as s (s.component.name)}
								<li>
									<span class="font-medium">{s.component.name}</span>
									{v(s.from)} → {v(s.to)}
									<span class="text-muted"
										>({s.level}, {s.delivered ? 'delivered' : 'incidental'}, {s.files.length} files)</span
									>
								</li>
							{/each}
						</ul>
					</li>
				{:else}
					<li>No release is computed for this project.</li>
				{/if}
			</ul>
		{/if}
		<div class="mt-4 flex justify-end gap-2">
			<button
				type="button"
				class="rounded border border-line-strong px-3 py-1"
				disabled={busy}
				onclick={oncancel}>Cancel</button
			>
			<button
				type="button"
				class="rounded bg-primary px-3 py-1 text-on-primary disabled:opacity-50"
				disabled={busy ||
					!!error ||
					!preview ||
					!!preview.blockers?.length ||
					!!preview.worktree_uncommitted?.length ||
					needsChoice}
				onclick={confirm}>{busy ? 'Accepting…' : 'Accept'}</button
			>
		</div>
	</div>
</div>
