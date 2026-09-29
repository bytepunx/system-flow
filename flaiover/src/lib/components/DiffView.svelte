<script lang="ts">
	// A story branch's changes: files with their counts, and each file's hunks in a panel its
	// summary line opens and closes (S-0041). In the panel the + or - of a line is in a margin,
	// a line has a dim background, and consecutive lines added or removed share one outline
	// (S-0164).
	import { patchRuns, type PatchKind } from '$lib/review';

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
	const uid = $props.id();
	let open = $state<Record<string, boolean>>({});

	// the outline round a run, and the background and text of each line in it
	const outline: Record<PatchKind, string> = {
		add: 'border-good',
		del: 'border-danger',
		hunk: 'border-transparent',
		context: 'border-transparent',
		note: 'border-transparent'
	};
	const line: Record<PatchKind, string> = {
		add: 'bg-good-soft text-good',
		del: 'bg-danger-soft text-danger',
		hunk: 'text-muted',
		context: '',
		note: 'text-muted'
	};
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
	{#each diff.files as f, i (f.path)}
		<li class="rounded border border-line bg-ground">
			<button
				type="button"
				class="flex w-full cursor-pointer flex-wrap items-baseline gap-2 px-2 py-1 text-left text-xs"
				aria-expanded={!!open[f.path]}
				aria-controls="{uid}-{i}"
				title={open[f.path] ? 'Hide the changes' : 'Show the changes'}
				data-testid="diff-toggle"
				onclick={() => (open[f.path] = !open[f.path])}
			>
				<span aria-hidden="true" class="w-3 shrink-0 text-muted">{open[f.path] ? '▾' : '▸'}</span>
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
				<div id="{uid}-{i}" class="border-t border-line" data-testid="diff-panel">
					{#if f.binary}
						<p class="px-2 py-1 text-xs text-muted">Binary file, no hunks.</p>
					{:else if !f.patch}
						<p class="px-2 py-1 text-xs text-muted">
							{f.truncated ? 'Left out for size.' : 'No textual changes.'}
						</p>
					{:else}
						<div class="overflow-x-auto p-1 font-mono text-xs leading-snug">
							<!-- as wide as the longest line, so an outline goes round the whole of its run -->
							<div class="w-max min-w-full space-y-px">
								{#each patchRuns(f.patch) as run, r (r)}
									<div class="rounded-sm border {outline[run.kind]}" data-run={run.kind}>
										{#each run.lines as l, n (n)}
											<div class="flex {line[l.kind]}" data-line={l.kind}>
												<span class="w-5 shrink-0 text-center select-none" data-sign>{l.sign}</span
												><span class="pr-2 whitespace-pre">{l.text || ' '}</span>
											</div>
										{/each}
									</div>
								{/each}
							</div>
						</div>
						{#if f.truncated}
							<p class="border-t border-line px-2 py-1 text-xs text-muted">Cut for size.</p>
						{/if}
					{/if}
				</div>
			{/if}
		</li>
	{/each}
</ul>
