// Which of the board's lanes are open on a narrow screen, where each lane collapses to its heading
// (S-0462): the lanes where the work moves are open until the user changes it. Like the board's
// type toggles, the choice is a per-browser convenience kept in localStorage. The lanes are the
// board's, as its workflow defines them (S-0753).
import type { Lane } from './lanes';

const KEY = 'flaiover-board-lanes';

/**
 * The lanes open by default: those of kind queue, active, or waiting, but the initial lane, where
 * new work collects; so with the baseline, ready, in-progress, and review.
 */
export function defaultOpen(lanes: Lane[]): string[] {
	const moving = ['queue', 'active', 'waiting'];
	return lanes.filter((l) => !l.initial && moving.includes(l.kind)).map((l) => l.name);
}

/** The stored open lanes in board order; the default unless what is stored names board lanes. */
export function parseOpen(raw: string | null, lanes: Lane[]): string[] {
	let v: unknown;
	try {
		v = raw === null ? null : JSON.parse(raw);
	} catch {
		v = null;
	}
	const names = lanes.map((l) => l.name);
	if (!Array.isArray(v) || !v.every((s) => typeof s === 'string' && names.includes(s)))
		return defaultOpen(lanes);
	return names.filter((s) => v.includes(s));
}

function read(): string | null {
	try {
		return localStorage.getItem(KEY);
	} catch {
		return null;
	}
}

export class LaneCollapse {
	/** The board's lanes, in order. */
	lanes = $state<Lane[]>([]);
	open = $state<string[]>([]);
	/** The choice as last stored or made, kept here too for a browser whose storage refuses it. */
	#raw = read();

	constructor(lanes: Lane[] = []) {
		this.lanes = lanes;
		this.open = parseOpen(this.#raw, lanes);
	}

	/** Take the board's lanes, and read the choice against them when they are not the ones it has. */
	setLanes(lanes: Lane[]): void {
		const same =
			lanes.length === this.lanes.length && lanes.every((l, i) => l.name === this.lanes[i].name);
		if (same) return;
		this.lanes = lanes;
		this.open = parseOpen(this.#raw, lanes);
	}

	/** Whether a lane shows its cards on a narrow board. */
	isOpen(lane: string): boolean {
		return this.open.includes(lane);
	}

	/** Open a closed lane or close an open one, and remember the whole choice. */
	toggle(lane: string): void {
		this.open = this.isOpen(lane)
			? this.open.filter((s) => s !== lane)
			: this.lanes.map((l) => l.name).filter((s) => s === lane || this.open.includes(s));
		this.#raw = JSON.stringify(this.open);
		try {
			localStorage.setItem(KEY, this.#raw);
		} catch {
			/* private windows: the choice lasts for the page */
		}
	}
}

export const laneCollapse = new LaneCollapse();
