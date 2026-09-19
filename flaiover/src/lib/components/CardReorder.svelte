<script lang="ts">
	// Up and down controls for a story card in a column that has a pull order (S-0057): the
	// alternative to dragging. They sit beside the card's link, never inside it, and show on
	// hover and on keyboard focus through the wrapper's `group`. A direction with nowhere to
	// go has no button.
	import type { Placement } from '$lib/reorder';

	let {
		id,
		up,
		down,
		onplace
	}: {
		id: string;
		up: Placement | null;
		down: Placement | null;
		onplace: (p: Placement, dir: 'up' | 'down') => void;
	} = $props();

	const button =
		'rounded border border-line-strong bg-surface px-1.5 leading-4 text-ink hover:bg-raised';
</script>

{#if up || down}
	<div
		class="absolute top-1.5 left-[60%] z-10 hidden -translate-x-1/2 gap-1 text-xs group-focus-within:flex group-hover:flex"
		data-testid="reorder"
	>
		{#if up}
			<button
				type="button"
				class={button}
				data-step="up"
				data-for={id}
				aria-label="Move {id} up in the pull order"
				title="Move up (Alt+↑)"
				onclick={() => onplace(up, 'up')}>↑</button
			>
		{/if}
		{#if down}
			<button
				type="button"
				class={button}
				data-step="down"
				data-for={id}
				aria-label="Move {id} down in the pull order"
				title="Move down (Alt+↓)"
				onclick={() => onplace(down, 'down')}>↓</button
			>
		{/if}
	</div>
{/if}
