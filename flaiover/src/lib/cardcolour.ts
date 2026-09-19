// How a board card is colour coded (S-0055): the background is a pastel tint for the
// item's nature, and a stripe on the left edge is a colour for its type. The colours are
// theme tokens in routes/layout.css. The card and the board's legend both read these maps,
// so they cannot drift apart. The class names are written out in full because Tailwind
// only generates the classes it can find as whole strings.
export const natureTint: Record<string, string> = {
	feature: 'bg-nature-feature',
	improvement: 'bg-nature-improvement',
	remediation: 'bg-nature-remediation',
	research: 'bg-nature-research',
	experiment: 'bg-nature-experiment'
};

export const typeStripe: Record<string, string> = {
	epic: 'border-l-type-epic',
	story: 'border-l-type-story',
	task: 'border-l-type-task'
};

/** The legend's swatch for a type: the stripe's colour as a fill. */
export const typeSwatch: Record<string, string> = {
	epic: 'bg-type-epic',
	story: 'bg-type-story',
	task: 'bg-type-task'
};

/** A nature the schema does not know keeps the plain card. */
export function tintFor(nature: string): string {
	return natureTint[nature] ?? 'bg-ground';
}

/** A known type gets a 4 px stripe that hovering leaves alone; an unknown one keeps the plain edge. */
export function stripeFor(type: string): string {
	const stripe = typeStripe[type];
	return stripe ? `border-l-4 ${stripe}` : 'hover:border-l-line-strong';
}
