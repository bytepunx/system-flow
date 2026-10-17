// The Design page's unsaved token edits kept in this browser, per project (S-0441, ADR-0151), so a
// reload keeps them until they are saved or discarded. A draft is the page's Edits: each value a
// DTCG $value keyed by editKey(theme, path), only the edits that differ from the saved tokens. The
// tried tokens are the saved themes with the draft laid over them, editedThemes(saved, draft, raw)
// from tokens.ts. Other tabs hear of a change through the storage event, which onDraft follows.
import { parseEditKey, type Edits } from './tokens';

export type { Edits };

const PREFIX = 'flaiover-design-draft:';

/** The localStorage key of a project's draft; null is the one project the server serves alone. */
export function draftKey(project: string | null): string {
	return PREFIX + (project ?? '');
}

/** The edits a stored draft holds; null when it does not parse or is not an object of edits. */
export function parseDraft(raw: string | null): Edits | null {
	if (raw === null) return null;
	let v: unknown;
	try {
		v = JSON.parse(raw);
	} catch {
		return null;
	}
	if (typeof v !== 'object' || v === null || Array.isArray(v)) return null;
	if (!Object.keys(v).every((k) => parseEditKey(k) !== null)) return null;
	return v as Edits;
}

/** The project's kept edits, or {} when none are kept; a draft that does not parse is dropped. */
export function loadDraft(project: string | null): Edits {
	let raw: string | null;
	try {
		raw = localStorage.getItem(draftKey(project));
	} catch {
		// a private window, or storage blocked: nothing is kept
		return {};
	}
	const edits = parseDraft(raw);
	if (raw !== null && edits === null) clearDraft(project);
	return edits ?? {};
}

/** Keeps the edits as the project's draft; no edits clears it. */
export function saveDraft(project: string | null, edits: Edits): void {
	if (Object.keys(edits).length === 0) {
		clearDraft(project);
		return;
	}
	try {
		localStorage.setItem(draftKey(project), JSON.stringify(edits));
	} catch {
		// a private window, storage blocked, or full: the edits last for this page only
	}
}

/** Drops the project's draft. */
export function clearDraft(project: string | null): void {
	try {
		localStorage.removeItem(draftKey(project));
	} catch {
		// storage blocked: there is nothing kept to drop
	}
}

/**
 * Calls fn with the project's draft each time another tab changes it, {} when it is cleared or
 * does not parse; the storage event does not fire in the tab that wrote. Returns the unsubscribe.
 */
export function onDraft(project: string | null, fn: (edits: Edits) => void): () => void {
	const key = draftKey(project);
	const heard = (e: StorageEvent) => {
		if (e.storageArea !== localStorage) return;
		// a null key is localStorage.clear(), which takes every draft with it
		if (e.key !== null && e.key !== key) return;
		fn(e.key === null ? {} : (parseDraft(e.newValue) ?? {}));
	};
	window.addEventListener('storage', heard);
	return () => window.removeEventListener('storage', heard);
}
