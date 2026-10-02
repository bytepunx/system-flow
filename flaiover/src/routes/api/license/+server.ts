import { json } from '@sveltejs/kit';
import { readLicense } from '$lib/server/license';

/**
 * GET: the license the dashboard is distributed under (S-0231), as the image carries it: its name
 * and its markdown. Nothing of the project is involved, so flai is not asked. An image built
 * without the file answers 503 and says so.
 */
export const GET = () => {
	const license = readLicense();
	if (!license) return json({ error: 'this image carries no LICENSE.md' }, { status: 503 });
	return json({ name: license.name, text: license.text });
};
