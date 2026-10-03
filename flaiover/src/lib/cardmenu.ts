// A board card's menu (S-0202): what the operator may do with the item from the card, as its page
// offers it, in the order the menu lists it.
import { agentAction, type StoryActivity } from './activity';

/** What an entry of the card's menu does. */
export type CardAction = 'open' | 'finalize' | 'agent' | 'block' | 'unblock' | 'cancel';

/** One entry of the card's menu: its action, its label, and for the agent, which flai runs. */
export type CardEntry = { action: CardAction; label: string; agent?: 'start' | 'restart' };

/** The card's menu: Open, then only what a writer may do with the item, in order. */
export function cardMenu(
	card: {
		id: string;
		type: string;
		status: string;
		blocked: boolean;
		draft?: boolean;
		archived?: boolean;
	},
	ctx: { writable: boolean; activity?: StoryActivity; agentEnabled: boolean }
): CardEntry[] {
	const entries: CardEntry[] = [{ action: 'open', label: 'Open' }];
	const story = card.type === 'story';
	const changeable = ctx.writable && !card.archived;
	const open = card.status !== 'done' && card.status !== 'cancelled';
	if (changeable && story && card.draft === true)
		entries.push({ action: 'finalize', label: 'Finalize' });
	const agent = story
		? agentAction(ctx.activity, ctx.agentEnabled, card.status, ctx.writable)
		: null;
	if (agent) entries.push({ action: 'agent', label: agent.label, agent: agent.action });
	if (changeable && card.blocked) entries.push({ action: 'unblock', label: 'Unblock' });
	else if (changeable && open) entries.push({ action: 'block', label: 'Block…' });
	if (changeable && open) entries.push({ action: 'cancel', label: 'Cancel…' });
	return entries;
}
