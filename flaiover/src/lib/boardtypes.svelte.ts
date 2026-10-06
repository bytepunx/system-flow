// Which work item types the board shows: the legend's type entries toggle each one (S-0303, was
// a checkbox each, S-0141), every type shown until the user changes it. Like the theme, the choice
// is a per-browser convenience kept in localStorage, so it lasts across reloads and navigation.
export type ItemType = 'epic' | 'story' | 'task';
export type Shown = Record<ItemType, boolean>;
export const itemTypes: ItemType[] = ['epic', 'story', 'task'];
const KEY = 'flaiover-board-types';
const defaults: Shown = { epic: true, story: true, task: true };

/** The stored choice, with the default for anything missing or not a boolean. */
export function parseShown(raw: string | null): Shown {
	let v: unknown;
	try {
		v = raw === null ? null : JSON.parse(raw);
	} catch {
		v = null;
	}
	const o = v !== null && typeof v === 'object' ? (v as Record<string, unknown>) : {};
	const out = { ...defaults };
	for (const t of itemTypes) if (typeof o[t] === 'boolean') out[t] = o[t] as boolean;
	return out;
}

function stored(): Shown {
	try {
		return parseShown(localStorage.getItem(KEY));
	} catch {
		return { ...defaults };
	}
}

export class BoardTypes {
	shown = $state<Shown>(stored());

	/** Show or hide one type and remember the whole choice. */
	set(type: ItemType, on: boolean): void {
		this.shown[type] = on;
		try {
			localStorage.setItem(KEY, JSON.stringify(this.shown));
		} catch {
			/* private windows: the choice lasts for the page */
		}
	}

	/** Flip one type, shown to hidden or back, and remember the whole choice. */
	toggle(type: ItemType): void {
		this.set(type, !this.shown[type]);
	}
}

export const boardTypes = new BoardTypes();
