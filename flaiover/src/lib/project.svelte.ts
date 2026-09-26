// Which project this browser tab is looking at (S-0080): one dashboard serves every project the
// host flai serves, and this is the one the designer picked, in api()'s calls and in the URL. The
// choice is remembered per browser (localStorage), the way the theme is; with exactly one project
// connected there is nothing to pick, and api() sends no ?project= at all, so a dashboard serving
// only one keeps working exactly as it always has.
const STORAGE_KEY = 'flaiover-project';
const URL_PARAM = 'project';

export type Project = {
	key: string;
	name: string;
	connected: boolean;
	since?: string;
	/** A repository offered for import, not a project yet (S-0098). */
	candidate?: boolean;
	/** flai serve serves it (S-0122). */
	served?: boolean;
	/** Why flai serve has it not connected, when it says. */
	lastError?: string;
};

function stored(): string | null {
	try {
		return localStorage.getItem(STORAGE_KEY);
	} catch {
		return null;
	}
}
function remember(key: string | null): void {
	try {
		if (key) localStorage.setItem(STORAGE_KEY, key);
		else localStorage.removeItem(STORAGE_KEY);
	} catch {
		// a private window, or storage blocked: the choice just does not survive a reload
	}
}

class ProjectState {
	/** Every project flai has ever named here, from /api/projects. */
	list = $state<Project[]>([]);
	loaded = $state(false);
	/** What api() sends as ?project=, and what the switcher shows as current. null means "let the
	 * server's own single-project fallback decide", which is right until there is more than one. */
	current = $state<string | null>(null);
	/** A repository being imported (S-0098): kept in the list, and so on the screen, after flai serve
	 * stops offering it, until the operator has read what became of it and moved on. */
	hold: Project | null = null;

	constructor() {
		if (typeof location !== 'undefined') {
			const fromUrl = new URLSearchParams(location.search).get(URL_PARAM);
			this.current = fromUrl || stored();
		}
	}

	/** The chosen project's row, once the list is known. */
	get project(): Project | undefined {
		return this.list.find((p) => p.key === this.current);
	}

	/** Whether more than one project is known: only then does the choice matter to anyone. */
	get needsChoice(): boolean {
		return this.list.length > 1;
	}

	/** The projects proper, leaving out repositories offered for import (S-0098). */
	get projects(): Project[] {
		return this.list.filter((p) => !p.candidate);
	}

	async refresh(): Promise<void> {
		try {
			const r = await fetch('/api/projects');
			if (!r.ok) return;
			const body = (await r.json()) as { projects?: Project[] };
			this.list = body.projects ?? [];
			const held = this.hold;
			if (held && !this.list.some((p) => p.key === held.key)) this.list = [...this.list, held];
			this.loaded = true;
			// A remembered or url-given key that no longer exists is not silently kept: with exactly
			// one project connected there is a right answer regardless of what was remembered before.
			if (this.list.length === 1) this.pick(this.list[0].key, false);
			else if (this.current && !this.list.some((p) => p.key === this.current)) this.pick(null);
			// one project and repositories offered for import: the project is still the one to show
			// until the operator picks something else (S-0098)
			if (!this.current && this.projects.length === 1) this.pick(this.projects[0].key, false);
		} catch {
			// offline, or nothing connected yet: keep what was chosen and try again later
		}
	}

	/** Choose a project. persist defaults to true; false is for the automatic single-project pick,
	 * which should not overwrite a choice the designer made deliberately in a multi-project session. */
	pick(key: string | null, persist = true): void {
		this.current = key;
		if (persist) remember(key);
		if (typeof location !== 'undefined' && typeof history !== 'undefined') {
			const url = new URL(location.href);
			if (key) url.searchParams.set(URL_PARAM, key);
			else url.searchParams.delete(URL_PARAM);
			history.replaceState(history.state, '', url);
		}
	}

	/** Add ?project=<current> to a same-origin API path, unless it already names one. */
	tag(path: string): string {
		if (!this.current) return path;
		const url = new URL(path, 'http://x');
		if (url.searchParams.has(URL_PARAM)) return path;
		url.searchParams.set(URL_PARAM, this.current);
		return url.pathname + url.search;
	}

	get ready(): boolean {
		return this.loaded;
	}
}

export const projectState = new ProjectState();

/** Test-only: undoes refresh(), back to how a fresh page load starts. Resetting `list` and
 * `current` alone would leave `loaded` stuck true from a previous test, and the next mount's
 * effect would then skip straight to the single-project path on whatever list is left over. */
export function resetForTests(): void {
	projectState.list = [];
	projectState.loaded = false;
	projectState.hold = null;
	projectState.pick(null, false);
}
