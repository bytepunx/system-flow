<script lang="ts">
	// Names the board's colours (S-0055). Built from the same maps as the cards. Each nature
	// tag is also a toggle (S-0302): a shown nature wears its brighter tint and a stronger
	// border, a hidden one the cards' pastel tint; aria-pressed and the title carry the state
	// for those who cannot tell the tints apart. The global :focus-visible rule rings them.
	import { natureShownTint, natureTint, typeSwatch } from '$lib/cardcolour';
	import { boardNatures, natures } from '$lib/boardnatures.svelte';
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
		{#each Object.entries(typeSwatch) as [type, swatch] (type)}
			<dd class="flex items-center gap-1" data-type={type}>
				<span class="inline-block h-3 w-1 rounded-sm {swatch}" aria-hidden="true"></span>{type}
			</dd>
		{/each}
	</div>
</dl>
