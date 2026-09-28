<script lang="ts">
	// A checkbox for each work item type the board can show (S-0141), remembered in this browser.
	import { boardTypes, itemTypes, type BoardTypes, type ItemType } from '$lib/boardtypes.svelte';

	let { types = boardTypes }: { types?: BoardTypes } = $props();
	const label: Record<ItemType, string> = { epic: 'epics', story: 'stories', task: 'tasks' };
</script>

<fieldset class="flex flex-wrap items-center gap-3 text-sm" data-testid="board-types">
	<legend class="sr-only">Show</legend>
	{#each itemTypes as type (type)}
		<label class="flex items-center gap-1"
			><input
				type="checkbox"
				data-type={type}
				checked={types.shown[type]}
				onchange={(e) => types.set(type, e.currentTarget.checked)}
			/>
			{label[type]}</label
		>
	{/each}
</fieldset>
