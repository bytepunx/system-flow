// Which of the board's lanes are open on a narrow screen, where each lane collapses to its heading
// (S-0462): Ready, In progress, and Review, where the work moves, are open until the user changes
// it. Like the board's type toggles, the choice is a per-browser convenience kept in localStorage.
export const laneStates = ['backlog', 'ready', 'in-progress', 'review', 'done', 'cancelled'];
const KEY = 'flaiover-board-lanes';
const defaults = ['ready', 'in-progress', 'review'];

/** The stored open lanes, or the default when what is stored is not a list of lane states. */
export function parseOpen(raw: string | null): string[] {
	let v: unknown;
	try {
		v = raw === null ? null : JSON.parse(raw);
	} catch {
		v = null;
	}
	if (!Array.isArray(v) || !v.every((s) => typeof s === 'string' && laneStates.includes(s)))
		return [...defaults];
	return laneStates.filter((s) => v.includes(s));
}

function stored(): string[] {
	try {
		return parseOpen(localStorage.getItem(KEY));
	} catch {
		return [...defaults];
	}
}

export class LaneCollapse {
	open = $state<string[]>(stored());

	/** Whether a lane shows its cards on a narrow board. */
	isOpen(lane: string): boolean {
		return this.open.includes(lane);
	}

	/** Open a closed lane or close an open one, and remember the whole choice. */
	toggle(lane: string): void {
		this.open = this.isOpen(lane)
			? this.open.filter((s) => s !== lane)
			: laneStates.filter((s) => s === lane || this.open.includes(s));
		try {
			localStorage.setItem(KEY, JSON.stringify(this.open));
		} catch {
			/* private windows: the choice lasts for the page */
		}
	}
}

export const laneCollapse = new LaneCollapse();
