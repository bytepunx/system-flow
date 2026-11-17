<script lang="ts">
	// A board card's menu (S-0202): what the item's page offers a writer ($lib/cardmenu), then, under
	// the lane's name, what the lane's menu offers there, so a right click on a card still reaches it.
	// Cancel… only where the workflow lets a person cancel the card (S-0753).
	import Menu, { type MenuClose } from './Menu.svelte';
	import { cardMenu, type CardEntry } from '$lib/cardmenu';
	import { cancelMove, laneEntries, type Flow, type LaneAction } from '$lib/lanes';
	import type { StoryActivity } from '$lib/activity';

	let {
		flow,
		card,
		lane,
		laneCount,
		writable,
		activity,
		agentEnabled,
		planEnabled,
		x,
		y,
		oncard,
		onlane,
		onclose
	}: {
		/** The board's lanes and the moves between them, as board.get answers them. */
		flow: Flow;
		card: Parameters<typeof cardMenu>[0];
		lane: string;
		/** How many cards the lane shows, for Archive All on cancelled (S-0343). */
		laneCount: number;
		writable: boolean;
		activity?: StoryActivity;
		agentEnabled: boolean;
		planEnabled: boolean;
		x: number;
		y: number;
		oncard: (entry: CardEntry) => void;
		onlane: (action: LaneAction) => void;
		onclose: (how: MenuClose) => void;
	} = $props();

	const entries = $derived(
		cardMenu(card, { flow, writable, activity, agentEnabled, planEnabled }).filter(
			(e) => e.action !== 'cancel' || cancelMove(flow, card.status, card.type) !== undefined
		)
	);
</script>

<Menu
	{x}
	{y}
	label="{card.id} actions"
	testid="card-menu"
	groups={[{ entries }, { heading: lane, entries: laneEntries(flow, lane, laneCount) }]}
	onpick={(action, group) => {
		if (group === 1) onlane(action as LaneAction);
		else oncard(entries.find((e) => e.action === action)!);
	}}
	{onclose}
/>
