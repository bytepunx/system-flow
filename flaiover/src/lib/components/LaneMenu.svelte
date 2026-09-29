<script lang="ts">
	// A lane's right-click menu on the board (S-0167): only what the lane allows, by flai's rules
	// ($lib/lanes). Escape, a click elsewhere, or a choice closes it.
	import { onMount, tick } from 'svelte';
	import { backOf, forwardOf, hasLimit } from '$lib/lanes';

	export type LaneAction = 'create' | 'forward' | 'back' | 'limit';

	let {
		lane,
		x,
		y,
		onpick,
		onclose
	}: {
		lane: string;
		x: number;
		y: number;
		onpick: (action: LaneAction) => void;
		onclose: () => void;
	} = $props();

	const entries = $derived(
		[
			{ action: 'create' as const, label: 'Create item here' },
			forwardOf(lane) && {
				action: 'forward' as const,
				label: `Move stories forward to ${forwardOf(lane)}…`
			},
			backOf(lane) && { action: 'back' as const, label: `Move stories back to ${backOf(lane)}…` },
			hasLimit(lane) && { action: 'limit' as const, label: 'Change WIP limit…' }
		].filter((e) => !!e)
	);

	let menu: HTMLDivElement;
	// none once the menu is gone: a choice can close it before focus moves in
	const items = () => [...(menu?.querySelectorAll<HTMLButtonElement>('[role="menuitem"]') ?? [])];
	function step(by: number) {
		const all = items();
		const at = all.indexOf(document.activeElement as HTMLButtonElement);
		all[(at + by + all.length) % all.length]?.focus();
	}

	onMount(() => {
		void tick().then(() => items()[0]?.focus());
		const away = (e: MouseEvent) => {
			if (!menu?.contains(e.target as Node)) onclose();
		};
		// a click that opened the menu must not close it: listen from the next turn
		const t = setTimeout(() => document.addEventListener('mousedown', away));
		return () => {
			clearTimeout(t);
			document.removeEventListener('mousedown', away);
		};
	});
</script>

<div
	bind:this={menu}
	class="fixed z-50 min-w-56 rounded border border-line bg-surface py-1 text-sm shadow-lg"
	style="left: {x}px; top: {y}px"
	role="menu"
	tabindex="-1"
	aria-label="{lane} lane"
	data-testid="lane-menu"
	onkeydown={(e) => {
		if (e.key === 'Escape') {
			e.preventDefault();
			onclose();
		} else if (e.key === 'ArrowDown') {
			e.preventDefault();
			step(1);
		} else if (e.key === 'ArrowUp') {
			e.preventDefault();
			step(-1);
		} else if (e.key === 'Tab') onclose();
	}}
>
	{#each entries as e (e.action)}
		<button
			type="button"
			role="menuitem"
			class="block w-full px-3 py-1.5 text-left hover:bg-raised focus:bg-raised focus:outline-none"
			data-action={e.action}
			onclick={() => onpick(e.action)}>{e.label}</button
		>
	{/each}
</div>
