// Which of the nav rail's groups are expanded (S-0458). The groups the operator expands or collapses
// are kept per browser in localStorage, like the board's types; a group they never chose starts
// collapsed. Arriving on a page expands its group for the page's life without storing it, so a
// group the operator collapsed stays collapsed after a reload anywhere else.
import type { GroupKey } from '$lib/sitemenu';

export type Chosen = Partial<Record<GroupKey, boolean>>;
const KEY = 'flaiover-nav-rail';
const groups: GroupKey[] = ['flow', 'agents', 'design', 'documentation', 'configuration'];

/** The stored choices, leaving out anything that is not a known group set to a boolean. */
export function parseChosen(raw: string | null): Chosen {
	let v: unknown;
	try {
		v = raw === null ? null : JSON.parse(raw);
	} catch {
		v = null;
	}
	const o = v !== null && typeof v === 'object' ? (v as Record<string, unknown>) : {};
	const out: Chosen = {};
	for (const g of groups) if (typeof o[g] === 'boolean') out[g] = o[g] as boolean;
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
	/** What the rail shows: the choices, and the groups expanded on arrival. */
	open = $state<Chosen>({ ...this.chosen });

	isExpanded(key: GroupKey): boolean {
		return this.open[key] ?? false;
	}

	/** Expand or collapse a group as the operator chose, and remember the choice. */
	toggle(key: GroupKey): void {
		const on = !this.isExpanded(key);
		this.open[key] = on;
		this.chosen[key] = on;
		try {
			localStorage.setItem(KEY, JSON.stringify(this.chosen));
		} catch {
			/* private windows: the choice lasts for the page */
		}
	}

	/** Expand the group of the page arrived on, without storing it as a choice. */
	expand(key: GroupKey): void {
		this.open[key] = true;
	}
}

export const navRail = new NavRailState();
