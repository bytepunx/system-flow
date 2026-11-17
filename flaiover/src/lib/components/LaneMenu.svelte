<script lang="ts" module>
	import type { LaneAction as Action } from '$lib/lanes';
	export type LaneAction = Action;
</script>

<script lang="ts">
	// A lane's right-click menu on the board (S-0167): only what the lane allows, by the workflow's
	// lanes and transitions ($lib/lanes). Escape, a click elsewhere, or a choice closes it.
	import Menu from './Menu.svelte';
	import { laneEntries, type Flow } from '$lib/lanes';

	let {
		flow,
		lane,
		count,
		x,
		y,
		onpick,
		onclose
	}: {
		/** The board's lanes and the moves between them, as board.get answers them. */
		flow: Flow;
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
	groups={[{ entries: laneEntries(flow, lane, count) }]}
	onpick={(action) => onpick(action as LaneAction)}
	{onclose}
/>
