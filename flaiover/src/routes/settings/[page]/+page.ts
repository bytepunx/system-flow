import { error } from '@sveltejs/kit';
import { SETTINGS_PAGES, settingsPage } from '$lib/settingspages';
import type { PageLoad } from './$types';

/** The Configuration page the slug names (S-0459), or a 404 naming the pages there are. */
export const load: PageLoad = ({ params }) => {
	const found = settingsPage(params.page);
	if (!found)
		error(
			404,
			`No settings page "${params.page}": open one of ${SETTINGS_PAGES.map((p) => `/settings/${p.slug}`).join(', ')}`
		);
	return found;
};
