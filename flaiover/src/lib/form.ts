// Forms filled from data a page reads again (S-0574): a live push, a reconnect, or a save elsewhere
// on the page reads the data again, and a form that refills from every read throws away what the
// user typed. rebase lays the user's edits over the new read instead.

type Plain = Record<string, unknown>;

function plain(v: unknown): v is Plain {
	return typeof v === 'object' && v !== null && !Array.isArray(v);
}

function same(a: unknown, b: unknown): boolean {
	if (plain(a) && plain(b)) {
		const keys = new Set([...Object.keys(a), ...Object.keys(b)]);
		for (const k of keys) if (!same(a[k], b[k])) return false;
		return true;
	}
	if (Array.isArray(a) && Array.isArray(b))
		return a.length === b.length && a.every((x, i) => same(x, b[i]));
	return Object.is(a, b);
}

/**
 * The form's values on a new read of its data: each field of fresh, except where current differs
 * from seed, the read current was filled from, which the user changed and keeps. Plain objects are
 * walked key by key; any other value, an array included, is one field. A key fresh lacks is gone,
 * whatever was typed in it; a key fresh adds comes in as read.
 */
export function rebase<T>(seed: T, current: T, fresh: T): T {
	if (plain(seed) && plain(current) && plain(fresh)) {
		const out: Plain = {};
		for (const k of Object.keys(fresh)) {
			out[k] = k in seed && k in current ? rebase(seed[k], current[k], fresh[k]) : fresh[k];
		}
		return out as T;
	}
	return same(current, seed) ? fresh : current;
}
