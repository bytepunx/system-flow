// The Configuration pages (S-0459, TH-0553), replacing the one settings page S-0458's rail linked
// into by anchors: five pages, each with its slug, its title, and the settings sections it shows,
// and the page an old /settings anchor lands on.

/** Every section of the settings, each once, named as its data-testid after section-. */
export const SETTINGS_SECTIONS = [
	'actions',
	'checks',
	'default-agent',
	'providers',
	'planning',
	'shared',
	'import',
	'projects',
	'sub-projects',
	'tokens',
	'agent'
] as const;

/** A section of the settings. */
export type SettingsSection = (typeof SETTINGS_SECTIONS)[number];

/** A Configuration page's slug, its path below /settings. */
export type SettingsPageSlug = 'features' | 'agents' | 'file-system' | 'projects' | 'harnesses';

/** A Configuration page: its slug, its title, and its sections in the order it shows them. */
export type SettingsPage = {
	slug: SettingsPageSlug;
	title: string;
	sections: readonly SettingsSection[];
};

/** The five Configuration pages, in the rail's order. */
export const SETTINGS_PAGES: readonly SettingsPage[] = [
	{ slug: 'features', title: 'Features', sections: ['actions', 'checks'] },
	{ slug: 'agents', title: 'Agents', sections: ['default-agent', 'providers', 'planning'] },
	{ slug: 'file-system', title: 'File System', sections: ['shared', 'import'] },
	{ slug: 'projects', title: 'Projects', sections: ['projects', 'sub-projects', 'tokens'] },
	{ slug: 'harnesses', title: 'Harnesses', sections: ['agent'] }
];

/** The Configuration page with this slug, or undefined when there is none. */
export function settingsPage(slug: string): SettingsPage | undefined {
	return SETTINGS_PAGES.find((p) => p.slug === slug);
}

/** The page an old /settings anchor (S-0458), with or without #, lands on: else Features. */
export function pageOfAnchor(hash: string): SettingsPageSlug {
	return settingsPage(hash.replace(/^#/, ''))?.slug ?? 'features';
}
