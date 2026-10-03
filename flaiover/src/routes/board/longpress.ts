// A long press on a card opens its menu on touch (S-0202): held still for PRESS_MS it fires, and
// the click the press ends in is swallowed, so the card's link does not navigate.

export const PRESS_MS = 500;
/** How far, in pixels, a finger may wander before the press is a scroll or a drag instead. */
export const PRESS_SLOP = 10;
/** How long after the finger lifts a click still belongs to a press that fired. */
const SWALLOW_MS = 600;

/** Follows one touch press at a time; a mouse or a pen press is left alone. */
export function longPress() {
	let timer: ReturnType<typeof setTimeout> | undefined;
	let from = { x: 0, y: 0 };
	let fired = false;
	const cancel = () => {
		clearTimeout(timer);
		timer = undefined;
	};
	return {
		/** A press begins: `fire` runs if it is held long enough without moving. */
		down(e: PointerEvent, fire: () => void) {
			cancel();
			fired = false;
			if (e.pointerType !== 'touch') return;
			from = { x: e.clientX, y: e.clientY };
			timer = setTimeout(() => {
				timer = undefined;
				fired = true;
				fire();
			}, PRESS_MS);
		},
		move(e: PointerEvent) {
			if (timer && Math.hypot(e.clientX - from.x, e.clientY - from.y) > PRESS_SLOP) cancel();
		},
		/** The finger lifts or the browser takes the press over. */
		up() {
			cancel();
			if (fired) setTimeout(() => (fired = false), SWALLOW_MS);
		},
		cancel,
		/** Whether the current press has fired, so the menu is already open. */
		get fired() {
			return fired;
		},
		/** Swallows the click that ends a press that fired. */
		click(e: MouseEvent) {
			if (!fired) return;
			fired = false;
			e.preventDefault();
			e.stopPropagation();
		}
	};
}
