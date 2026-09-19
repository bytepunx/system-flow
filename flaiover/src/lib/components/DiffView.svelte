<script lang="ts">
	// A story branch's changes: files with their counts, hunks on demand.
	// Additions and deletions are coloured and keep their + and - (S-0041).
	import { patchLines } from '$lib/review';

	type File = {
		path: string;
		old_path?: string;
		status: string;
		additions: number;
		deletions: number;
		binary: boolean;
		truncated: boolean;
		patch: string;
	};
	type Diff = {
		branch: string;
		base: string;
		commits: number;
		files: File[];
		additions: number;
		deletions: number;
		truncated: boolean;
	};

	let { diff }: { diff: Diff } = $props();
	let open = $state<Record<string, boolean>>({});
</script>

<p class="text-xs text-muted">
	<code>{diff.branch}</code> against <code>{diff.base}</code>: {diff.commits} commit{diff.commits ===
	1
		? ''
		: 's'}, {diff.files.length} file{diff.files.length === 1 ? '' : 's'},
	<span class="text-good">+{diff.additions}</span>
	<span class="text-danger">−{diff.deletions}</span>
	{#if diff.truncated}· some patches were cut for size; see the rest with <code>git diff</code>{/if}
</p>
{#if !diff.files.length}
	<p class="mt-2 text-sm text-muted">The branch changes no files.</p>
{/if}
<ul class="mt-2 space-y-1">
	{#each diff.files as f (f.path)}
		<li class="rounded border border-line bg-ground">
			<button
				type="button"
				class="flex w-full flex-wrap items-baseline gap-2 px-2 py-1 text-left text-xs"
				aria-expanded={!!open[f.path]}
				onclick={() => (open[f.path] = !open[f.path])}
			>
				<span class="w-16 shrink-0 text-muted">{f.status}</span>
				<span class="min-w-0 flex-1 font-mono break-all"
					>{#if f.old_path}{f.old_path} →
					{/if}{f.path}</span
				>
				{#if f.binary}<span class="text-muted">binary</span>{:else}
					<span class="text-good">+{f.additions}</span>
					<span class="text-danger">−{f.deletions}</span>
				{/if}
			</button>
			{#if open[f.path]}
				{#if f.binary}
					<p class="border-t border-line px-2 py-1 text-xs text-muted">Binary file, no hunks.</p>
				{:else if !f.patch}
					<p class="border-t border-line px-2 py-1 text-xs text-muted">
						{f.truncated ? 'Left out for size.' : 'No textual changes.'}
					</p>
				{:else}
					<pre
						class="overflow-auto border-t border-line px-2 py-1 font-mono text-xs leading-snug">{#each patchLines(f.patch) as l, i (i)}<span
								class="block {l.kind === 'add'
									? 'bg-good-soft text-good'
									: l.kind === 'del'
										? 'bg-danger-soft text-danger'
										: l.kind === 'hunk'
											? 'text-muted'
											: ''}">{l.text || ' '}</span
							>{/each}</pre>
					{#if f.truncated}
						<p class="border-t border-line px-2 py-1 text-xs text-muted">Cut for size.</p>
					{/if}
				{/if}
			{/if}
		</li>
	{/each}
</ul>
