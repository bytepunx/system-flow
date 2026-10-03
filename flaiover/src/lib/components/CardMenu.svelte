<script lang="ts">
	// A board card's menu (S-0202): what the item's page offers a writer ($lib/cardmenu), then, under
	// the lane's name, what the lane's menu offers there, so a right click on a card still reaches it.
	import Menu, { type MenuClose } from './Menu.svelte';
	import { cardMenu, type CardEntry } from '$lib/cardmenu';
	import { laneEntries, type LaneAction } from '$lib/lanes';
	import type { StoryActivity } from '$lib/activity';

	let {
		card,
		lane,
		writable,
		activity,
		agentEnabled,
		x,
		y,
		oncard,
		onlane,
		onclose
	}: {
		card: Parameters<typeof cardMenu>[0];
		lane: string;
		writable: boolean;
		activity?: StoryActivity;
		agentEnabled: boolean;
		x: number;
		y: number;
		oncard: (entry: CardEntry) => void;
		onlane: (action: LaneAction) => void;
		onclose: (how: MenuClose) => void;
	} = $props();

	const entries = $derived(cardMenu(card, { writable, activity, agentEnabled }));
</script>

<Menu
	{x}
	{y}
	label="{card.id} actions"
	testid="card-menu"
	groups={[{ entries }, { heading: lane, entries: laneEntries(lane) }]}
	onpick={(action, group) => {
		if (group === 1) onlane(action as LaneAction);
		else oncard(entries.find((e) => e.action === action)!);
	}}
	{onclose}
/>
