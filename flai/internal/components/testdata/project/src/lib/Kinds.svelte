<script lang="ts">
	// One prop of each kind, typed by an inline literal.
	import type { Snippet } from 'svelte';
	import type { Theme } from './theme';

	let {
		label,
		count = 3,
		open = false,
		size = 'md',
		tags,
		ids = [],
		options,
		meta,
		onpick,
		children,
		theme,
		ratio
	}: {
		label: string;
		count?: number;
		open?: boolean;
		size?: 'sm' | 'md' | 'lg' | undefined;
		tags: string[];
		ids?: Array<number>;
		options: { dense: boolean };
		meta: Record<string, unknown>;
		onpick: (id: number) => void;
		children?: Snippet<[string]>;
		theme: Theme;
		ratio?: number | null;
	} = $props();
</script>

<p>{label} {count} {open} {size} {tags} {ids} {options} {meta} {theme} {ratio}</p>
<button onclick={() => onpick(count)}>{@render children?.(label)}</button>
