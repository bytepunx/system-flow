// S-0231: the License page under Host shows the license the image carries.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));

import LicensePage from './+page.svelte';

const answer = (body: unknown, ok = true) => ({ ok, json: async () => body, statusText: 'nope' });
const settle = async () => {
	for (let i = 0; i < 10; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};

describe('the license page (S-0231)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});

	it('renders the license markdown it is given', async () => {
		api.mockImplementation(async (url: string) => {
			if (url === '/api/license')
				return answer({
					name: 'Example License 2.0',
					text: '# Example License 2.0\n\n## Acceptance\n\nUse it **well**.\n'
				});
			return answer({});
		});
		c = mount(LicensePage, { target: document.body });
		await settle();
		expect(api).toHaveBeenCalledWith('/api/license');
		const article = document.querySelector('article');
		expect(article?.getAttribute('aria-label')).toBe('Example License 2.0');
		expect(article?.querySelector('h2')?.textContent).toContain('Acceptance');
		expect(article?.querySelector('strong')?.textContent).toBe('well');
	});

	it('says why when the image carries no license', async () => {
		api.mockImplementation(async () =>
			answer({ error: 'this image carries no LICENSE.md' }, false)
		);
		c = mount(LicensePage, { target: document.body });
		await settle();
		expect(document.querySelector('article')).toBeNull();
		expect(document.body.textContent).toContain('this image carries no LICENSE.md');
	});
});
