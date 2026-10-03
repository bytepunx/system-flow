<script lang="ts" module>
	/** One entry of a menu: what it does and what it says. */
	export type MenuEntry = { action: string; label: string };
	/** Entries shown together, under an optional heading, divided from the group before. */
	export type MenuGroup = { heading?: string; entries: MenuEntry[] };
	/** How a menu closed without a choice: Escape, or a click or focus elsewhere. */
	export type MenuClose = 'escape' | 'away';
</script>

<script lang="ts">
	// A menu at a point on the page, the lane's (S-0167) and the card's (S-0202): its first entry
	// takes focus, the arrows move through it, and Escape, Tab, a click elsewhere, focus leaving it,
	// or a choice closes it. It stays inside the viewport.
	import { onMount, tick } from 'svelte';

	let {
		x,
		y,
		label,
		testid,
		groups,
		onpick,
		onclose
	}: {
		x: number;
		y: number;
		/** The menu's accessible name. */
		label: string;
		testid?: string;
		groups: MenuGroup[];
		/** The entry chosen and the index of its group. */
		onpick: (action: string, group: number) => void;
		onclose: (how: MenuClose) => void;
	} = $props();

	let menu: HTMLDivElement;
	// how far the menu moves back from x,y so it does not overflow the right or bottom edge
	let back = $state({ x: 0, y: 0 });
	// none once the menu is gone: a choice can close it before focus moves in
	const items = () => [...(menu?.querySelectorAll<HTMLButtonElement>('[role="menuitem"]') ?? [])];
	function step(by: number) {
		const all = items();
		const at = all.indexOf(document.activeElement as HTMLButtonElement);
		all[(at + by + all.length) % all.length]?.focus();
	}

	onMount(() => {
		const box = menu.getBoundingClientRect();
		back = {
			x: box.right > window.innerWidth ? box.width : 0,
			y: box.bottom > window.innerHeight ? box.height : 0
		};
		void tick().then(() => items()[0]?.focus());
		const away = (e: MouseEvent) => {
			if (!menu?.contains(e.target as Node)) onclose('away');
		};
		// a click that opened the menu must not close it: listen from the next turn
		const t = setTimeout(() => document.addEventListener('mousedown', away));
		return () => {
			clearTimeout(t);
			document.removeEventListener('mousedown', away);
		};
	});
</script>

<!-- The menu itself takes focus (tabindex -1), so a click on an entry that does not focus it, as in
     Safari, keeps focus inside and is not taken for focus leaving. -->
<div
	bind:this={menu}
	class="fixed z-50 min-w-56 rounded border border-line bg-surface py-1 text-sm shadow-lg"
	style="left: {Math.max(0, x - back.x)}px; top: {Math.max(0, y - back.y)}px"
	role="menu"
	tabindex="-1"
	aria-label={label}
	data-testid={testid}
	onkeydown={(e) => {
		if (e.key === 'Escape') {
			e.preventDefault();
			onclose('escape');
		} else if (e.key === 'ArrowDown') {
			e.preventDefault();
			step(1);
		} else if (e.key === 'ArrowUp') {
			e.preventDefault();
			step(-1);
		} else if (e.key === 'Tab') onclose('away');
	}}
	onfocusout={(e) => {
		if (!menu?.contains(e.relatedTarget as Node | null)) onclose('away');
	}}
>
	{#each groups as g, i (i)}
		{#if i > 0}<div class="my-1 border-t border-line" role="separator"></div>{/if}
		<div role="group" aria-label={g.heading}>
			{#if g.heading}
				<div class="px-3 pt-1 pb-0.5 text-xs text-muted" aria-hidden="true">{g.heading}</div>
			{/if}
			{#each g.entries as e (e.action)}
				<button
					type="button"
					role="menuitem"
					class="block w-full px-3 py-1.5 text-left hover:bg-raised focus:bg-raised focus:outline-none"
					data-action={e.action}
					onclick={() => onpick(e.action, i)}>{e.label}</button
				>
			{/each}
		</div>
	{/each}
</div>
