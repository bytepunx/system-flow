<script lang="ts">
	// Names the board's colours (S-0055). Built from the same maps as the cards.
	import { natureTint, typeSwatch } from '$lib/cardcolour';
</script>

<dl class="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-muted" data-testid="legend">
	<div class="flex flex-wrap items-center gap-x-2 gap-y-1">
		<dt>nature</dt>
		{#each Object.entries(natureTint) as [nature, tint] (nature)}
			<dd class="rounded border border-line px-1.5 py-0.5 text-ink {tint}" data-nature={nature}>
				{nature}
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
