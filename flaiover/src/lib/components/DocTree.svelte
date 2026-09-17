<script lang="ts">
	import { resolve } from '$app/paths';
	import DocTree from './DocTree.svelte';

	type Node = {
		name: string;
		path: string;
		kind: 'dir' | 'file';
		title?: string;
		children?: Node[];
	};
	let { nodes, current, depth = 0 }: { nodes: Node[]; current: string; depth?: number } = $props();

	// folders containing the current document start open
	let open = $state<Record<string, boolean>>({});
	const isOpen = (n: Node) => open[n.path] ?? (depth === 0 || current.startsWith(n.path + '/'));
</script>

<ul class={depth === 0 ? '' : 'ml-3 border-l border-zinc-200 pl-2 dark:border-zinc-800'}>
	{#each nodes as n (n.path)}
		<li class="py-0.5">
			{#if n.kind === 'dir'}
				<button
					type="button"
					class="flex w-full items-center gap-1 text-left text-sm font-medium text-zinc-700 hover:text-zinc-900 dark:text-zinc-300 dark:hover:text-zinc-100"
					onclick={() => (open[n.path] = !isOpen(n))}
				>
					<span class="inline-block w-3 text-xs text-zinc-400">{isOpen(n) ? '▾' : '▸'}</span
					>{n.name}
				</button>
				{#if isOpen(n) && n.children?.length}
					<DocTree nodes={n.children} {current} depth={depth + 1} />
				{/if}
			{:else}
				<a
					href={resolve('/docs/[...path]', { path: n.path })}
					class="block truncate rounded px-1 text-sm {current === n.path
						? 'bg-zinc-200 font-medium dark:bg-zinc-800'
						: 'text-zinc-600 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-100'}"
					title={n.path}
				>
					{n.title ?? n.name}
				</a>
			{/if}
		</li>
	{/each}
</ul>
