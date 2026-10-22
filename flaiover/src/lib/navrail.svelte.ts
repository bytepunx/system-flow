// Which of the nav rail's groups are expanded (S-0458). The groups the operator expands or collapses
// are kept per browser in localStorage, like the board's types; a group they never chose starts
// collapsed. Arriving on a page expands its group for the page's life without storing it, so a
// group the operator collapsed stays collapsed after a reload anywhere else. Whether Charts is
// expanded is kept with the groups' choices, under the same key, and in the same way (S-0484).
import type { GroupKey } from '$lib/sitemenu';

/** What the rail expands and collapses: its groups, and the Charts page with its charts. */
export type RailKey = GroupKey | 'charts';
export type Chosen = Partial<Record<RailKey, boolean>>;
const KEY = 'flaiover-nav-rail';
const keys: RailKey[] = ['flow', 'agents', 'design', 'documentation', 'configuration', 'charts'];

/** The stored choices, leaving out anything that is not a known key set to a boolean. */
export function parseChosen(raw: string | null): Chosen {
	let v: unknown;
	try {
		v = raw === null ? null : JSON.parse(raw);
	} catch {
		v = null;
	}
	const o = v !== null && typeof v === 'object' ? (v as Record<string, unknown>) : {};
	const out: Chosen = {};
	for (const k of keys) if (typeof o[k] === 'boolean') out[k] = o[k] as boolean;
	return out;
}

function stored(): Chosen {
	try {
		return parseChosen(localStorage.getItem(KEY));
	} catch {
		return {};
	}
}

export class NavRailState {
	/** What the operator chose, as stored. */
	private chosen: Chosen = stored();
	/** What the rail shows: the choices, and what was expanded on arrival. */
	open = $state<Chosen>({ ...this.chosen });

	isExpanded(key: RailKey): boolean {
		return this.open[key] ?? false;
	}

	/** Expand or collapse a group, or Charts, as the operator chose, and remember the choice. */
	toggle(key: RailKey): void {
		const on = !this.isExpanded(key);
		this.open[key] = on;
		this.chosen[key] = on;
		try {
			localStorage.setItem(KEY, JSON.stringify(this.chosen));
		} catch {
			/* private windows: the choice lasts for the page */
		}
	}

	/** Expand the group, or Charts, of the page arrived on, without storing it as a choice. */
	expand(key: RailKey): void {
		this.open[key] = true;
	}
}

export const navRail = new NavRailState();
