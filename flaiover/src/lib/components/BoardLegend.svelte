<script lang="ts">
	// Names the board's colours (S-0055). Built from the same maps as the cards. Each nature
	// tag is also a toggle (S-0302): a shown nature wears its brighter tint and a stronger
	// border, a hidden one the cards' pastel tint. The type entries are toggles too (S-0303):
	// a shown type is highlighted in a tint of its own colour, a hidden one keeps the plain
	// look. aria-pressed and the title carry the state for those who cannot tell the colourings
	// apart. The global :focus-visible rule rings them.
	import { natureShownTint, natureTint, typeHighlight, typeSwatch } from '$lib/cardcolour';
	import { boardNatures, natures } from '$lib/boardnatures.svelte';
	import { boardTypes, itemTypes, type BoardTypes, type ItemType } from '$lib/boardtypes.svelte';

	let { types = boardTypes }: { types?: BoardTypes } = $props();
	const plural: Record<ItemType, string> = { epic: 'epics', story: 'stories', task: 'tasks' };
</script>

<dl class="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-muted" data-testid="legend">
	<div class="flex flex-wrap items-center gap-x-2 gap-y-1">
		<dt>nature</dt>
		{#each natures as nature (nature)}
			{@const shown = boardNatures.shown[nature]}
			<dd>
				<button
					type="button"
					class="cursor-pointer rounded border px-1.5 py-0.5 text-ink {shown
						? `border-line-strong ${natureShownTint[nature]}`
						: `border-line ${natureTint[nature]}`}"
					data-nature={nature}
					aria-pressed={shown}
					title="click to {shown ? 'hide' : 'show'} {nature}"
					onclick={() => boardNatures.toggle(nature)}
				>
					{nature}
				</button>
			</dd>
		{/each}
	</div>
	<div class="flex flex-wrap items-center gap-x-2 gap-y-1">
		<dt>type</dt>
		{#each itemTypes as type (type)}
			{@const shown = types.shown[type]}
			<dd>
				<button
					type="button"
					class="flex cursor-pointer items-center gap-1 rounded border px-1.5 py-0.5 {shown
						? `text-ink ${typeHighlight[type]}`
						: 'border-line'}"
					data-type={type}
					aria-pressed={shown}
					title="click to {shown ? 'hide' : 'show'} {plural[type]}"
					onclick={() => types.toggle(type)}
					><span class="inline-block h-3 w-1 rounded-sm {typeSwatch[type]}" aria-hidden="true"
					></span>{type}</button
				>
			</dd>
		{/each}
	</div>
</dl>
