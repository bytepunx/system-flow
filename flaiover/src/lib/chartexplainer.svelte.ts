// Whether a chart's explainer is open (S-0485): open on the first visit, then as the user last
// left it. Like the window, the choice is a per-browser convenience kept in localStorage, and one
// choice serves every chart, so a user who has read the explainers keeps them closed everywhere.
const KEY = 'flaiover-chart-explainer';
const DEFAULT = true;

/** The stored choice: closed only when the user closed it, open for anything else. */
export function parseOpen(raw: string | null): boolean {
	return raw === 'closed' ? false : raw === 'open' ? true : DEFAULT;
}

function stored(): boolean {
	try {
		return parseOpen(localStorage.getItem(KEY));
	} catch {
		return DEFAULT;
	}
}

export class ChartExplainer {
	open = $state<boolean>(stored());

	/** Open or close the explainer on every chart, and remember it. */
	set(open: boolean): void {
		this.open = open;
		try {
			localStorage.setItem(KEY, open ? 'open' : 'closed');
		} catch {
			/* private windows: the choice lasts for the page */
		}
	}
}

export const chartExplainer = new ChartExplainer();
