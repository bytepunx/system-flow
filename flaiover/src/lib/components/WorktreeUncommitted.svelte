<script lang="ts" module>
	/**
	 * The blockers other than flai's own for the worktree's uncommitted paths, which this box shows
	 * with the way to have them committed (S-0140).
	 */
	export function otherBlockers(preview: {
		blockers?: string[];
		worktree_uncommitted?: string[];
	}): string[] {
		const all = preview.blockers ?? [];
		if (!preview.worktree_uncommitted?.length) return all;
		return all.filter(
			(b) => !(b.startsWith('the worktree ') && b.includes('has uncommitted changes'))
		);
	}
</script>

<script lang="ts">
	// Work a story's agent left uncommitted in its worktree blocks acceptance: the branch is merged
	// as it is committed and the worktree removed. The operator can have an agent commit it
	// (flai's agent.commit, S-0140).
	import { api } from '$lib/api';

	let {
		id,
		paths,
		branch,
		disabled = false
	}: { id: string; paths: string[]; branch?: string; disabled?: boolean } = $props();

	let starting = $state(false);
	let started = $state<string | null>(null);
	let error = $state<string | null>(null);

	async function haveCommitted() {
		starting = true;
		error = null;
		try {
			const r = await api(`/api/items/${id}/agent`, {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ action: 'commit' })
			});
			const data = await r.json();
			if (!r.ok) error = data.error ?? r.statusText;
			else
				started = `${data.agent ?? 'An agent'} is committing them (pid ${data.pid}). Accept once it has finished: look again then.`;
		} catch (e) {
			error = e instanceof Error ? e.message : String(e);
		} finally {
			starting = false;
		}
	}
</script>

<div
	class="mt-3 rounded border border-danger bg-danger-soft p-2 text-danger"
	role="group"
	data-testid="worktree-uncommitted"
>
	<p class="font-medium">Uncommitted changes in the story's worktree:</p>
	<ul class="mt-1 ml-4 list-disc font-mono text-xs">
		{#each paths as p (p)}<li>{p}</li>{/each}
	</ul>
	<p class="mt-2">
		Acceptance merges the branch as it is committed and removes the worktree, so it waits until
		these are committed on {branch ?? 'the story branch'}.
	</p>
	{#if started}
		<p class="mt-2 font-medium" role="status">{started}</p>
	{:else}
		<button
			type="button"
			class="mt-2 rounded border border-line-strong bg-surface px-2 py-1 text-ink hover:bg-raised disabled:opacity-50"
			disabled={starting || disabled}
			onclick={haveCommitted}
			>{starting ? 'Starting an agent…' : 'Have an agent commit them'}</button
		>
	{/if}
	{#if error}<p class="mt-2" role="alert">{error}</p>{/if}
</div>
