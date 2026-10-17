<!--
	The Components page's source (S-0442): the main checkout, or a story in progress or in review,
	whose worktree the catalogue and the dev server are then read from. With a story, the components
	its branch changes are listed here, before the catalogue, each marked and picked as in the list.
-->
<script lang="ts" module>
	import type { Component } from '$lib/server/repo';

	/** A story the picker offers, as the board's card gives it. */
	export type SourceStory = { id: string; title: string; status: string };
	/** A component the story's branch changes, with its sub-project's manifest name. */
	export type ChangedComponent = { project: string; component: Component };

	/** The board's stories in progress, then those in review: the sources besides Main. */
	export function openStories(
		columns: Partial<Record<string, readonly (SourceStory & { type: string })[]>>
	): SourceStory[] {
		return [...(columns['in-progress'] ?? []), ...(columns.review ?? [])]
			.filter((c) => c.type === 'story')
			.map(({ id, title, status }) => ({ id, title, status }));
	}
</script>

<script lang="ts">
	let {
		stories,
		story,
		onchoose,
		changed,
		changedError = null,
		picked,
		onpick
	}: {
		/** The stories in progress or in review, in the order the board gives them. */
		stories: readonly SourceStory[];
		/** The story picked; undefined for the main checkout. */
		story: string | undefined;
		/** Called with the story chosen, or undefined for Main. */
		onchoose: (story: string | undefined) => void;
		/** The components the story's branch changes; null until its diff is read. */
		changed: readonly ChangedComponent[] | null;
		/** Why the story's diff could not be read; null when it was. */
		changedError?: string | null;
		/** The component framed: its sub-project's manifest name and its path. */
		picked: { project: string; path: string } | null;
		/** Called with the sub-project's manifest name and the component's path. */
		onpick: (project: string, path: string) => void;
	} = $props();

	// A story the URL names that is no longer in progress or in review is still offered, so that the
	// picker shows what the page reads; flai refuses it once it has no worktree.
	const listed = $derived(
		story && !stories.some((s) => s.id === story)
			? [...stories, { id: story, title: 'not in progress or in review', status: '' }]
			: stories
	);
</script>

<div class="flex flex-col gap-2" data-testid="source-picker">
	<label class="text-sm">
		<span class="mb-1 block text-xs text-muted">Source</span>
		<select
			class="max-w-md rounded border border-line-strong bg-surface px-2 py-1"
			value={story ?? ''}
			onchange={(e) => onchoose(e.currentTarget.value || undefined)}
			data-testid="source-picker-select"
		>
			<option value="">Main</option>
			{#each listed as s (s.id)}
				<option value={s.id}>{s.id} {s.title}{s.status ? ` (${s.status})` : ''}</option>
			{/each}
		</select>
	</label>
	{#if story}
		<p class="text-xs text-muted">
			The catalogue and the dev server are those of the worktree of
			<span class="font-mono">story/{story}</span>.
		</p>
		<section aria-labelledby="source-changed-heading" data-testid="source-changed">
			<h2 id="source-changed-heading" class="text-sm font-semibold">Changed on its branch</h2>
			{#if changedError}
				<p class="mt-1 text-xs text-warn" data-testid="source-changed-error">
					Could not read what its branch changes: {changedError}
				</p>
			{:else if changed === null}
				<p class="mt-1 text-xs text-muted">Reading its branch…</p>
			{:else if !changed.length}
				<p class="mt-1 text-xs text-muted" data-testid="source-changed-none">
					Its branch changes no component.
				</p>
			{:else}
				<ul class="mt-1 flex flex-wrap gap-2">
					{#each changed as { project, component: c } (`${project}:${c.path}`)}
						{@const current = picked?.project === project && picked.path === c.path}
						<li>
							<button
								type="button"
								class={[
									'rounded border px-2 py-0.5 text-sm',
									current
										? 'border-primary bg-primary-soft font-medium'
										: 'border-line hover:bg-raised'
								]}
								aria-current={current ? 'true' : undefined}
								title="{project}: {c.path}"
								data-testid="source-changed-{c.path}"
								onclick={() => onpick(project, c.path)}
							>
								{c.name}
								<span class="text-xs text-primary" data-testid="source-changed-mark">changed</span>
							</button>
						</li>
					{/each}
				</ul>
			{/if}
		</section>
	{/if}
</div>
