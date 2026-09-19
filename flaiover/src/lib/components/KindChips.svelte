<script lang="ts">
	// An item's nature on its tint and its type beside its stripe colour, as on the
	// board's cards and legend (S-0055). The words stay: colour only repeats them.
	import { natureTint, typeSwatch } from '$lib/cardcolour';

	let { type, nature }: { type?: string; nature?: string } = $props();
</script>

{#if type}<span class="whitespace-nowrap" data-type={type}
		>{#if typeSwatch[type]}<span
				class="mr-1 inline-block h-3 w-1 rounded-sm align-[-1px] {typeSwatch[type]}"
				aria-hidden="true"
			></span>{/if}{type}</span
	>{/if}{#if type && nature}<span class="mx-1">·</span>{/if}{#if nature}<span
		class={natureTint[nature]
			? `rounded border border-line px-1.5 text-ink ${natureTint[nature]}`
			: ''}
		data-nature={nature}>{nature}</span
	>{/if}
