// Which natures the board shows (S-0302): the legend's nature tags toggle each one, every
// nature shown until the user changes it. Like the item types, the choice is a per-browser
// convenience kept in localStorage, so it lasts across reloads and navigation.
export type Nature = 'feature' | 'improvement' | 'remediation' | 'research' | 'experiment';
export type ShownNatures = Record<Nature, boolean>;
export const natures: Nature[] = [
	'feature',
	'improvement',
	'remediation',
	'research',
	'experiment'
];
const KEY = 'flaiover-board-natures';
const defaults: ShownNatures = {
	feature: true,
	improvement: true,
	remediation: true,
	research: true,
	experiment: true
};

/** The stored choice, with the default for anything missing or not a boolean. */
export function parseShownNatures(raw: string | null): ShownNatures {
	let v: unknown;
	try {
		v = raw === null ? null : JSON.parse(raw);
	} catch {
		v = null;
	}
	const o = v !== null && typeof v === 'object' ? (v as Record<string, unknown>) : {};
	const out = { ...defaults };
	for (const n of natures) if (typeof o[n] === 'boolean') out[n] = o[n] as boolean;
	return out;
}

function stored(): ShownNatures {
	try {
		return parseShownNatures(localStorage.getItem(KEY));
	} catch {
		return { ...defaults };
	}
}

function isNature(nature: string): nature is Nature {
	return (natures as string[]).includes(nature);
}

export class BoardNatures {
	shown = $state<ShownNatures>(stored());

	/** Show or hide one nature and remember the whole choice. */
	set(nature: Nature, on: boolean): void {
		this.shown[nature] = on;
		try {
			localStorage.setItem(KEY, JSON.stringify(this.shown));
		} catch {
			/* private windows: the choice lasts for the page */
		}
	}

	/** Flip one nature, shown to hidden or back, and remember the whole choice. */
	toggle(nature: Nature): void {
		this.set(nature, !this.shown[nature]);
	}

	/** Whether a card of this nature shows; one the schema does not know is never hidden. */
	isShown(nature: string): boolean {
		return !isNature(nature) || this.shown[nature];
	}
}

export const boardNatures = new BoardNatures();
