import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';

const goto = vi.fn(async () => {});
vi.mock('$app/navigation', () => ({ goto: (...a: unknown[]) => goto(...(a as [])) }));
vi.mock('$app/state', () => ({
	page: { url: new URL('http://192.168.0.162:4242/login?next=/board') }
}));
vi.mock('$app/paths', () => ({ resolve: (p: string) => p }));

import Login from './+page.svelte';

const settle = async () => {
	for (let i = 0; i < 6; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const error = () => document.querySelector('[data-testid="login-error"]')?.textContent ?? null;

async function submit(token: string) {
	const input = document.querySelector<HTMLInputElement>('input[type="password"]')!;
	input.value = token;
	input.dispatchEvent(new Event('input', { bubbles: true }));
	flushSync();
	document
		.querySelector('form')!
		.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));
	await settle();
}

// S-0083: a right token and a dropped cookie used to look like the page asking again, for no reason.
describe('the login page', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		goto.mockClear();
		vi.unstubAllGlobals();
		document.body.innerHTML = '';
	});
	const answers = (login: number, session: number) =>
		vi.stubGlobal(
			'fetch',
			vi.fn(async (url: string) => ({
				ok: (url === '/api/login' ? login : session) < 400,
				status: url === '/api/login' ? login : session
			}))
		);

	it('goes on when the token is accepted and the session is kept', async () => {
		answers(200, 200);
		c = mount(Login, { target: document.body });
		await submit('the-token');
		expect(error()).toBeNull();
		expect(goto).toHaveBeenCalledWith('/board');
	});

	it('says the token was not accepted', async () => {
		answers(401, 401);
		c = mount(Login, { target: document.body });
		await submit('nope');
		expect(error()).toContain('not accepted');
		expect(goto).not.toHaveBeenCalled();
	});

	it('says so when the token is right and the browser did not keep the cookie, and stays', async () => {
		answers(200, 401);
		c = mount(Login, { target: document.body });
		await submit('the-token');
		expect(error()).toContain('The token is right');
		expect(error()).toContain('did not keep the session cookie');
		expect(goto).not.toHaveBeenCalled();
	});
});
