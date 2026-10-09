<script lang="ts" module>
	import type { LaneAction as Action } from '$lib/lanes';
	export type LaneAction = Action;
</script>

<script lang="ts">
	// A lane's right-click menu on the board (S-0167): only what the lane allows, by flai's rules
	// ($lib/lanes). Escape, a click elsewhere, or a choice closes it.
	import Menu from './Menu.svelte';
	import { laneEntries } from '$lib/lanes';

	let {
		lane,
		count,
		x,
		y,
		onpick,
		onclose
	}: {
		lane: string;
		/** How many cards the lane shows, for Archive All on cancelled (S-0343). */
		count: number;
		x: number;
		y: number;
		onpick: (action: LaneAction) => void;
		onclose: () => void;
	} = $props();
</script>

<Menu
	{x}
	{y}
	label="{lane} lane"
	testid="lane-menu"
	groups={[{ entries: laneEntries(lane, count) }]}
	onpick={(action) => onpick(action as LaneAction)}
	{onclose}
/>
